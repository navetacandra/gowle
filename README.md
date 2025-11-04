# GOWLE

Gowle is a lightweight and high-performance file system watcher written in Go.
It monitors project directories in real-time and automatically restarts or reloads your application whenever file changes are detected.
Designed to provide a fast and dependency-free alternative to tools like nodemon, Gowle delivers an efficient development workflow with minimal resource usage.

## Features

- **Lightweight and Fast**: Gowle is designed to be lightweight and fast, with minimal resource usage.
- **Cross-Platform**: Gowle is cross-platform and can be used on Linux, macOS, and Windows.
- **Dependency-Free**: Gowle is dependency-free and can be used without any external dependencies.
## Installation

To install Gowle, you need to build it from source.

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/navetacandra/gowle.git
    cd gowle
    ```

2.  **Build the executable:**

    ```bash
    mkdir -p build
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o build/gowle cmd/gowle.go
    ```

    This will create an executable named `gowle` (or `gowle.exe` on Windows) in the `build` directory.

3.  **Move the executable to your PATH (optional):**

    You can move the `gowle` executable to a directory in your system's PATH to run it from anywhere. For example:

    ```bash
    sudo mv build/gowle /usr/local/bin/
    ```

## Usage

To use Gowle, you can run the following command in your project directory:

```bash
gowle
```

Gowle will automatically detect the changes in your project directory and restart your application.

You can also use the following commands:

- `rs`: to manually restart the application.
- `.exit`: to stop Gowle.

### Options

- `-v`: to enable verbose mode and print every change.

## Configuration

Gowle can be configured with a `.gowle` file in your project directory.

Here is an example of a `.gowle` file:

```
WATCH=src
IGNORE=.git,vendor,node_modules
COMMAND=go run main.go --port 8080
```

- `WATCH`: (Optional) The directory to watch for changes. Defaults to the current directory if not specified.
- `IGNORE`: (Optional) A comma-separated list of directories or files to ignore.
- `COMMAND`: (Optional) The command to execute when changes are detected. Defaults to `go run <entrypoint>` if not specified.

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

Gowle is licensed under the MIT License.
