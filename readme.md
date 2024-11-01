# GophKeeper

![GophKeeper](https://pictures.s3.yandex.net/resources/gophkeeper_2x_1650456239.png)

## Description

GophKeeper is a client-server system for the secure storage of logins, passwords, and other private data. The main functionality is concentrated in the client part, which allows users to manage their data. The client part can operate without synchronization with the server, providing the possibility of autonomous data management. The client part stores data in the user's profile folder, settings in a JSON file, and data in an SQLite database.

## Dependencies

- Go
- Postgres (if not using Docker)
- Docker and Docker Compose
- make

## Installation

1. Clone this repository and navigate to the folder:

   ```bash
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Use make to build and run the applications:
    - To build and run the client:

      ```bash
      make run_client
      ```

    - To build and run the server:

      ```bash
      make run_server
      ```

3. The server and database can be started using Docker Compose:

   ```bash
   make docker-up
   ```

## Features

### Client Part:

- **Help**: Each command has help available with the `--help` and `-h` flags.  
  For example: `gophkeeper save --help` or `gophkeeper save card --help`.
- **Encryption**: Encryption is performed using a randomly generated encryption key. This encryption key is generated once, when first needed, encrypted with a passphrase requested from the user, and stored in the profile settings. Subsequent access to the encryption key is through a passphrase request.
- **List of saved data**: Retrieving a list of keys of saved data does not require a passphrase or encryption key - only open data is displayed.
- **Data request by key**: Users can access their data by requesting it with a unique key. To do this, you need to enter a passphrase that unlocks the encryption key to unpack the encrypted data.
- **Data synchronization**: The ability to synchronize data to the server specified in the settings. Data is transmitted in the same encrypted form as stored in the local database. The user's settings on the server store the encryption key encrypted with the passphrase. Registration on the server is done with a separate synchronization password, and further authorization is done with a client token obtained during registration.

### Server Part

- User registration and authentication.
- Storage of the encrypted encryption key.
- Storage of the user's encrypted data.
- Data synchronization with authorized clients.

### Examples of using the client part

- **Version information, help**

  ```bash
  gophkeeper
  gophkeeper -v
  gophkeeper -h
  gophkeeper [command] -h
  ```

- **Shell**

  ```bash
  gophkeeper shell
  ```

- **Saving data**:

  ```bash
  gophkeeper save --help
  gophkeeper save bin -f filename -d "description"
  gophkeeper save card --num 2222-4444-5555-1111 --exp 10/29 --cvv 123 --owner "Max Space"
  gophkeeper save auth -l login -p password -k "my-key-name" -d site.com
  ```

- **Retrieving data by key**:

  ```bash
  gophkeeper view <key name>
  ```

- **Settings**

  ```bash
  gophkeeper config global
  gophkeeper config user
  ```

  - **Setting email, synchronization server**

      ```bash
      gophkeeper config user -e <email>
      gophkeeper config user -s <server address for synchronization>
      ```

- **Synchronization with a remote server**
  - **Registration**:

      ```bash
      gophkeeper sync register
      ```

    - **Synchronization**:

      ```bash
      gophkeeper sync now
      ```

    - **Changing the authorization password on the server**:

      ```bash
      gophkeeper sync password
      ```

    - **Complete account deletion from the server**:

      ```bash
      gophkeeper sync delete
      ```

## License

This project is licensed under the terms of the [MIT License](LICENSE).
