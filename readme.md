# GophKeeper
![GophKeeper](https://pictures.s3.yandex.net/resources/gophkeeper_2x_1650456239.png)

## Description
GophKeeper is a client—server system for the secure storage of usernames, passwords and other private data. The main functionality is concentrated in the client part, which allows users to manage their data. The client part can work without synchronization with the server, providing the possibility of autonomous data management.

## Functionality

### The client part:
- **Authentication and authorization**: Registration and login to access data.
- **Offline operation**: The ability to manage data locally without syncing with the server.
- **Key Data Request**: Users can access their data by requesting it with a unique key.
- **Data synchronization**: The ability to synchronize data between the client and the server.

### Examples of using the client part:
- **Registration of a new client**:
  gophkeeper view <key name>
- **Getting data by key**:
  gophkeeper get <key name>
- **Adding a new element (synchronization)**:
  gophkeeper sync <key name>

### The backend:
- Registration and authentication of users.
- Data synchronization with authorized clients.
- Storage and management of private data.

### Synchronization examples:
- **Data synchronization with the server**:
  gophkeeper sync
- **Deleting an item and syncing**:
  gophkeeper delete <key name>
  gophkeeper sync

## Installation

### Installation using Makefile

Clone this repository and go to the folder with it.  
Use make for build, run and other, for example, to build and launch the applications:
   - To build and launch the client:  
     make run_client
   - To build and run the server:  
     make run_server

## Architecture
GophKeeper consists of client and server parts implemented in the Go language. The client provides a command-line interface for data management, and the server uses gRPC to process requests and synchronize data. All data is encrypted to ensure security, which guarantees the protection of confidential user information.

## License
This project is licensed under the terms of the [MIT License] (LICENSE).