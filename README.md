# Background
When `create-react-app` starts a new build, it wipes out the contents in the `build/` directory.

This can break production if production reads from `build/`.

This tool copies a `source` folder into a destination folder with timestamped dates.

The idea is to have a directory structure [similar to Capistrano](https://capistranorb.com/documentation/getting-started/structure/).

Sample directory structure:

```
builds/
├── 20201025 121625
│   └── static
│       ├── css
│       └── js
├── 20201025 121829
│   └── static
│       ├── css
│       └── js
├── 20201025 124435
│   └── static
│       ├── css
│       └── js
├── 20201025 125148
│   └── static
│       ├── css
│       └── js
├── 20201025 130305
│   └── static
│       ├── css
│       └── js
└── current -> 20201025 130305
```

After building, `dirlink` copies the contents into `builds/$timestamp`, and then create a symbolic link from `builds/current` to `builds/$timestamp`.

## Testing

After a successful build run `dirlink build builds`.

You will see a `builds` directory create and populated.

## Setup yarn

Update your `package.json` `build` script to run `dirlink build builds`

```json
{
  "scripts": {
    "build": "react-scripts build && dirlink build builds"
  }
}
```

Execute: `yarn build`

## Setting up Nginx
If using `docker-compose`, update your `docker-compose.yml` and point the volume to `builds` folder
```yml
volumes:
- /home/aizat/src/app/builds:/var/www/app
```

## Nginx Config
In the `nginx` configuration, point the `root` path to the `current` directory.
```
root /var/www/app/current;
```