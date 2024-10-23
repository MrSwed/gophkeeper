---

### English Version

# GophKeeper

![GophKeeper](https://pictures.s3.yandex.net/resources/gophkeeper_2x_1650456239.png)

## Description

GophKeeper is a client-server system that allows users to securely store logins, passwords, binary data, and other private information. The system provides users with the ability to authenticate, authorize, and synchronize data between multiple clients.

## Functional Features

### Server Side:
- User registration, authentication, and authorization.
- Storage of private data.
- Synchronization of data between multiple authorized clients of the same owner.
- Providing private data to the owner upon request.

### Client Side:
- User authentication and authorization on a remote server.
- Access to private data upon request.

## Installation

### Requirements
- Go 1.16 or higher
- gRPC protocol

### Installing the Server Side
1. Clone the repository:
   git clone <URL>
   cd gophKeeper/server
2. Install dependencies:
   go mod tidy
3. Run the server:
   go run main.go
### Installing the Client Side
1. Clone the repository:
   git clone <URL>
   cd gophKeeper/client
2. Install dependencies:
   go mod tidy
3. Run the client:
   go run main.go
## Usage

### Registering a New User
1. Start the client and run the registration command:
   gophkeeper register --email <email> --password <password>
### Authentication
1. Start the client and run the authentication command:
   gophkeeper login --email <email> --password <password>
### Data Synchronization
1. After authentication, run the command to synchronize data:
   gophkeeper sync
### Accessing Data
1. Request data:
   gophkeeper get --key <key>
## Architecture

GophKeeper consists of client and server parts implemented in Go. The server uses gRPC to handle requests, while the client provides a command-line interface for interaction with the server. Data is encrypted to ensure security.

## Testing

The code of the system is covered by unit tests at least 80%. Each exported function, type, variable, and package contains comprehensive documentation.

## License

This project is licensed under the terms of the [MIT License](LICENSE).

---

Если у вас есть дополнительные пожелания или изменения, дайте знать!