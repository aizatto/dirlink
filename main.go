package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aizatto/dirlink/copy"
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Usage: dirlink <source> <destination>")
		os.Exit(1)
	}

	source := os.Args[1]
	destination := os.Args[2]
	errs := Dirlink(source, destination)
	if len(errs) != 0 {
		for _, err := range errs {
			log.Println(err)
		}
		os.Exit(1)
	}
}

func Dirlink(source, destinationroot string) []error {
	var errs []error
	if _, err := os.Stat(source); os.IsNotExist(err) {
		errs = append(errs, err)
	}

	stat, err := os.Lstat(destinationroot)
	if os.IsNotExist(err) {
		err = os.Mkdir(destinationroot, 0744)
	} else if !stat.IsDir() {
		err = fmt.Errorf("Path is not a directory: %s", destinationroot)
	}

	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return errs
	}

	wd, err := os.Getwd()
	if err != nil {
		return []error{err}
	}

	basename := time.Now().Format("20060102 150405")

	destination := ResolvePath(wd, path.Join(destinationroot, basename))

	err = copy.DirCopy(source, destination, copy.Content, true)
	// err = os.Rename(source, destination)
	if err != nil {
		return []error{err}
	}

	now := time.Now()
	err = os.Chtimes(destination, now, now)
	if err != nil {
		return []error{err}
	}

	symlinkroot := ResolvePath(wd, destinationroot)
	// symlinkdir := ResolvePath(wd, path.Join(destinationroot, "current"))
	symlinkdir := path.Join(symlinkroot, "current")
	_, err = os.Lstat(symlinkdir)
	if os.IsNotExist(err) {
		log.Printf("symlink does not exist: %s\n", symlinkdir)
	} else {
		log.Printf("removing: %s\n", symlinkdir)
		err = os.Remove(symlinkdir)
		if err != nil {
			err2 := os.Remove(destination)
			return []error{err, err2}
		}
	}

	log.Printf("Copying %s to %s\n", source, destination)
	log.Printf("Symlinking %s to %s\n", destination, symlinkdir)

	var files FileInfos
	files, err = ioutil.ReadDir(symlinkroot)
	if err != nil {
		return []error{err}
	}

	if len(files) > 5 {
		sort.Sort(files)
		for _, file := range files[5:] {
			log.Printf("%s: %s", file.Name(), file.ModTime())
			olddir := path.Join(symlinkroot, file.Name())
			err = os.RemoveAll(olddir)
			if err != nil {
				errs = append(errs, err)
			} else {
				log.Printf("Removing old directory: %s\n", olddir)
			}
		}
	}

	if len(errs) != 0 {
		return errs
	}

	err = os.Symlink(basename, symlinkdir)
	if err != nil {
		return []error{err}
	}
	return nil
}

func ResolvePath(wd string, givenpath string) string {
	if strings.HasPrefix(givenpath, filepath.VolumeName(givenpath)) {
		return path.Clean(givenpath)
	}

	return path.Clean(path.Join(wd, givenpath))
}

type FileInfos []os.FileInfo

func (p FileInfos) Len() int {
	return len(p)
}

func (p FileInfos) Less(i, j int) bool {
	return p[i].ModTime().After(p[j].ModTime())
}

func (p FileInfos) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
