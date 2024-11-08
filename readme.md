# GophKeeper

![GophKeeper](https://pictures.s3.yandex.net/resources/gophkeeper_2x_1650456239.png)

## Description

GophKeeper is a client-server system for securely storing logins, passwords, and other private data. The main functionality is concentrated in the client part, which allows users to manage their data. The client part can operate without synchronization with the server, providing the ability to manage data autonomously. The client stores data in the user's profile folder, settings in a JSON file, and data in an SQLite database.

## Dependencies

- Go
- Postgres (if Docker is not used)
- Docker and Docker Compose
- make

## Installation

1. Clone this repository and navigate to its folder:

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

### Client Side:

- **Help**: Help is available for each command with the `--help` and `-h` flags.  
  For example: `gophkeeper save --help` or `gophkeeper save card --help`.
- **Encryption**: Encryption is performed using a randomly generated encryption key. This encryption key is generated once, when first needed, encrypted with a passphrase requested from the user, and stored in the profile settings. Subsequent access to the encryption key is through a passphrase request.
- **List of Saved Data**: Retrieving a list of keys of existing saved data does not require a passphrase or encryption key - only open data is displayed.
- **Data Request by Key**: Users can access their data by requesting it with a unique key. To do this, they need to enter a passphrase that unlocks the encryption key to unpack the encrypted data.
- **Data Synchronization**: Ability to synchronize data to the server specified in the settings. Data is transmitted in the same encrypted form as it is stored in the local database. The encryption key, encrypted with a passphrase, is stored in the user's settings on the server. Registration on the server is done with a separate synchronization password, and further authorization is done with a client token obtained during registration.

### Server Side

- User registration and authentication.
- Storage of the encrypted encryption key.
- Storage of the user's encrypted data.
- Data synchronization with authorized clients.

### Client Usage Examples

#### Version Information, Help

  ```bash
  gophkeeper
  gophkeeper -v
  gophkeeper -h
  gophkeeper [command] -h
  ```

#### Shell

```bash
gophkeeper shell
```

#### Saving Data

When saving data, the key is automatically generated based on the data type and current date; if necessary, a custom key can be set using the `-k|--key` parameter.

```bash
gophkeeper save [auth|bin|card|text] [args]
```

```bash
gophkeeper save --help
gophkeeper save text -k "manual key name" -t "some text"
gophkeeper save bin -f filename -d "description"
gophkeeper save card --num 2222-4444-5555-1111 --exp 10/29 --cvv 123 --owner "Max Space"
gophkeeper save auth -l login -p password -k "my-key-name" -d site.com
```

#### Retrieving Data by Key

```bash
gophkeeper view <key name>
```

#### Settings

```bash
gophkeeper config global
gophkeeper config user
```

##### Configuring Email, Synchronization Server

```bash
gophkeeper config user -e <email>
gophkeeper config user -s <synchronization server address>
```

#### Synchronization with Remote Server

##### Registration

  ```bash
  gophkeeper sync register
  ```

##### Synchronization

```bash
gophkeeper sync now
```

##### Changing Server Authorization Password

```bash
gophkeeper sync password
```

##### Complete Account Deletion from Server

```bash
gophkeeper sync delete
```

## License

This project is licensed under the terms of the [MIT License](LICENSE).
