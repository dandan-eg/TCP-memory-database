# TCP Memory Database

This TCP Memory Database is a simple educational project aimed at providing a better understanding of the TCP protocol and implementing a practical application. It incorporates the Observer pattern to notify subscribers and follows the principle of dependency inversion by introducing interfaces for the Saver and Sourcer components.

## Purpose

The purpose of this project is to gain a better understanding of TCP protocol and implement a functional database application. It demonstrates the usage of the Observer pattern for notification and dependency inversion by decoupling the DB structure from the Saver and Sourcer interfaces.

## Usage

To use the TCP Memory Database, follow these steps:

1. Clone the repository to your local machine.
2. Run the project by executing `go run main.go` in the terminal.
3. Optionally, you can specify additional flags:
    - `-out`: Specify the save path for the database data. Default is set to `./`.
    - `-save`: Specify the save format (json or csv). Default is set to `json`.
    - `-src`: Specify the source file path for loading initial data.
4. Interact with the database by establishing a TCP connection to the server. You can use telnet or any other TCP client.
5. Use the following actions to interact with the database:
    - `GET <key>`: Retrieve the value associated with a key from the database.
    - `SET <key> <value>`: Set a key-value pair in the database.
    - `DEL <key>`: Delete the key-value pair from the database.
    - `EXIT`: Exit the database application.
    - `SAVE`: Save the current state of the database to a file. The file format (JSON or CSV) is determined by the `-save` flag specified during program execution.
    - `CLOSE`: Close the database connection.
6. The database server will process your commands and provide appropriate responses based on the actions performed.

Please note that while this application has been tested, it is intended for educational purposes and not meant for production use. It serves as a basic demonstration of a memory database. For more robust and feature-rich database solutions, it is recommended to use established tools like Redis.

## Dependencies

This project relies solely on the native Go libraries and does not have any external dependencies.

## Project Structure

The project structure is as follows:

- `saver/`: Contains the components responsible for saving the database data.
    - `saver.go`: Defines the Saver interface.
    - `csv.go`: Implements the Saver interface for CSV format.
    - `json.go`: Implements the Saver interface for JSON format.

- `src/`: Contains the components responsible for loading initial data into the database.
    - `source.go`: Defines the Source interface.
    - `csv.go`: Implements the Source interface for CSV format.
    - `json.go`: Implements the Source interface for JSON format.

- `memory/`: Contains the main database functionality.
    - `db.go`: Defines the DB struct and its methods.
    - `action.go`: Defines the database actions.
    - `display.go`: Handles sending messages to the database connections.

- `main.go`: Entry point of the application.

## Conclusion

The TCP Memory Database project provides a simple implementation of a memory database with basic functionality. It can be used to gain insights into the TCP protocol and serve as an educational tool. However, for more advanced features, performance, and production-ready applications, it is recommended to use established database solutions like Redis or other similar tools.

If you have any questions or need further assistance, feel free to reach out.
