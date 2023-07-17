# Key Value DB (Redis v0.0.1) in Go

This readme file provides an overview and instructions for using the keyvaluedb CLI tool written in Golang.

## Description

The CLI tool is a key-value database (KVDB) that allows users to interact with a simple in-memory database through a TCP server. It provides a command-line interface for executing various commands and retrieving results.

## Installation

To use the CLI tool, you need to have Golang(1.18) installed on your machine. Follow the steps below to install and set up the tool:

1. Clone the repository or download the source code files.
2. Open a terminal and navigate to the project directory.
3. Run the following command to build the CLI tool:

   ```shell
   go build -o kvdb
   ```

   This will create an executable file named `kvdb` in the project directory.

4. Optionally, you can move the `kvdb` executable to a directory included in your system's `PATH` environment variable for easier access.

## Usage

To start the TCP server and use the CLI tool, follow these steps:

1. Ensure the environment variable `APP_PORT` is set to the desired port number on which the TCP server should listen. For example, you can set it to `9736` by running:

   ```shell
   export APP_PORT=9736
   ```

   Replace `9736` with the desired port number.

2. Optionally, if you want to specify the number of in-memory databases (`DB_COUNT`), set the environment variable as well. For example:

   ```shell
   export DB_COUNT=16
   ```

   Replace `16` with the desired number of databases. If not set, the default value is `16`.

3. Run the following command to start the TCP server:

   ```shell
   ./kvdb
   ```

   Replace `./kvdb` with the actual path to the `kvdb` executable if it's not in the current directory or not in your `PATH`.

4. The TCP server will start and display a message indicating that it is listening on the specified port.

5. Open another terminal or use a tool like `nc` to connect to the TCP server. For example:

   ```shell
   nc localhost 9736
   ```

   Replace `localhost` with the appropriate host if the server is running on a different machine, and `9736` with the correct port number.

6. Once connected, you can start interacting with the CLI tool by entering commands. The command prompt is denoted by a `$` symbol.

7. The available commands are case-insensitive and can be entered in the following format:

   ```
   COMMAND [argument1] [argument2] ...
   ```

   Replace `COMMAND` with one of the supported commands and provide the necessary arguments.

8. The CLI tool supports the following commands:
  
    - `SET key value`: Sets the value of the specified key in the current database.
    - `GET key`: Retrieves the value of the specified key from the current database.
    - `DEL key`: Deletes the specified key from the current database.
    - `INCR key`: Increments the value of the specified key by 1.
    - `INCRBY key increment`: Increments the value of the specified key by the specified increment.
    - `MULTI`: Starts a transaction block.
    - `EXEC`: Executes all commands in a transaction block.
    - `DISCARD`: Discards all commands in a transaction block.
    - `COMPACT`: Compacts the database by removing expired keys.
    - `SELECT` index: Switches to the specified database index (0-based).

    Replace key, value, index, and increment with the appropriate values.

9. After entering a command, the CLI tool will display the result of the command. If the result is a list, each item will be numbered.

10. To exit the CLI tool, close the `nc` connection or terminate the terminal session.

## Dependencies

The CLI tool depends on the following external packages:

- `github.com/joho/godotenv`: Used for loading environment variables from a `.env` file.

Make sure to install these dependencies using a package manager like `go get` or by including them in your Go module dependencies.