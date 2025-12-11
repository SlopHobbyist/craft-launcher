# Craft Launcher - Update Server

This folder contains the infrastructure for hosting the modpack update server.
It uses **Nginx** to serve files and a custom script to generate the update manifest automatically.

## Structure

```
craftlauncher-server-side/
├── docker-compose.yml       # Docker deployment config
├── nginx.conf               # Nginx server configuration
├── generate-manifest.sh     # Script to generate manifest.json
└── files/                   # (Created at runtime) Modpack files go here
```

## How It Works

1.  **Nginx** serves the static files (mods, configs, assets) via HTTP.
2.  **Manifest Generator** runs on container startup. It scans the `files/` directory and creates `manifest.json`.
3.  **The Launcher** downloads `manifest.json`, compares it with the local state, and downloads changed files.

## Setup & Deployment

1.  **Prepare Files**:
    Put your entire modpack structure (mods, config, options.txt, etc.) into a folder named `files` in this directory (or mount it via volume).
    
    ```text
    files/
    ├── mods/
    │   └── my-mod.jar
    ├── config/
    ├── options.txt
    └── servers.dat
    ```

2.  **Version Control**:
    Create a file named `.version` inside `files/`. This file should contain a single integer (e.g., `5`).
    Increment this number to trigger an update for all clients.
    
    `echo "5" > files/.version`

3.  **Run with Docker**:
    ```bash
    docker-compose up -d
    ```

## File Overrides

The `generate-manifest.sh` script automatically handles "safe" updates for user configuration files.

-   **Enforced Files** (`override: true`): Most files (mods, scripts). If the user deletes or modifies them, the launcher will repair them.
-   **User Files** (`override: false`): Configs like `options.txt`, `servers.dat`, `usercache.json`.
    -   If the user has these files, the launcher **will not** download the server's version (preserving user settings).
    -   If the user is missing these files (fresh install), the launcher **will** download the server's version (providing defaults).

## Adding New Mods

1.  Stop the server or just access the volume.
2.  Add the new `.jar` to `files/mods/`.
3.  Remove the old `.jar`.
4.  Increment the version in `files/.version`.
5.  Restart the container to regenerate the manifest:
    ```bash
    docker-compose restart
    ```
