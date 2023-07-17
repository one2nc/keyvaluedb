# Key Value DB (Redis v0.0.1) in Go

This repository is a reference implementation of the problem statement available at https://playbook.one2n.in/go-bootcamp/go-projects/key-value-db-redis-in-go/key-value-db-redis-exercise. The solution is a command-line utility that allows you to interact with a simple in-memory key-value database through a TCP server. It provides a command-line interface for executing various commands and retrieving results.

## Installation

You must have Go(1.18) installed on your machine to use the CLI tool. Follow the steps below to install and set up the tool:

1. Clone the repository or download the source code files.
2. Open a terminal and navigate to the project directory.
3. Run the following command to build the CLI tool:

   ```shell
   go build -o kvdb
   ```

   This will create an executable file named `kvdb` in the project directory.

4. Optionally, you can move the `kvdb` executable to a directory in your system's `PATH` environment variable for easier access.

## Usage

To start the TCP server and use the CLI tool, follow these steps:


1. Create a `.env` file in the root directory and set the desired values to the below environment variables. We have provided a `.env.example` file in the root directory for reference.

   1. Ensure the environment variable `APP_PORT` is set to the desired port number on which the TCP server should listen. For example, you can set it to `9736` by running:

      ```shell
      export APP_PORT=9736
      ```

      Replace `9736` with the desired port number.

   2. Optionally, set the environment variable to specify the number of in-memory databases (`DB_COUNT`). For example:

      ```shell
      export DB_COUNT=16
      ```

      Replace `16` with the desired number of databases. If not set, the default value is `16`.

2. Run the following command to start the TCP server:

   ```shell
   ./kvdb
   ```

   Replace `./kvdb` with the actual path to the `kvdb` executable if it's not in the current directory or your `PATH`.

3. The TCP server will start and display a message indicating it listens on the specified port.

4. Open another terminal and use a tool like `nc` to connect to the TCP server. For example:

   ```shell
   nc localhost 9736
   ```

   Replace `localhost` with the appropriate host if the server runs on a different machine and `9736` with the correct port number.

5. Once connected, you can interact with the CLI tool by entering commands. A `$` symbol denotes the command prompt.

6. The available commands are case-insensitive and can be entered in the following format:

   ```
   COMMAND [argument1] [argument2] ...
   ```

   Replace `COMMAND` with a supported command and provide the necessary arguments.

7. The CLI tool supports the following commands:
  
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

8. After entering a command, the CLI tool will display the command result. If the result is a list, each item will be numbered.

9. To exit the CLI tool, close the `nc` connection or terminate the terminal session.

## Dependencies

The CLI tool depends on the following external packages:

- `github.com/joho/godotenv`: Used for loading environment variables from a `.env` file.

Make sure to install these dependencies using a package manager like `go get` or by including them in your Go module dependencies.

## License
This project is licensed under the [MIT License](./LICENSE)
