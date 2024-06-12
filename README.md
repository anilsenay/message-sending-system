# message-sending-system

---

This project is an automated message sending service. The system is designed to periodically retrieve a specified number of unsent messages from the database and send them to the designated recipients. The retrieval and sending process takes place at regular intervals to ensure that messages are delivered on time. The interval duration and the number of messages retrieved can be configured based on the specific requirements of the application.

---

- [Getting Started](#getting-started)
- [Project Layout](#project-layout)
- [Swagger](#swagger)
- [Diagrams](#diagrams)

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker

### How to run

1. Clone the project

```sh
git clone git@github.com:anilsenay/message-sending-system.git
cd message-sending-system
```

2. Build and Deploy project via Docker Compose

```sh
docker-compose up --build
```

3. Visit [http://localhost:8080/swagger/](http://localhost:8080/swagger/)

---

## Project Layout

    .
    ├── cmd                       # Main applications
    │   └── server                  # Server app's main.go placed here
    │
    ├── internal                  # Private application packages
    │   ├── apps                    # Dependency injections and application setups
    │   ├── client                  # Clients used in application such as redis and message client
    │   ├── config                  # The configs to be used in the application
    │   ├── handlers                # The "endpoints" that the application serves
    │   ├── models                  # The data models used in the application
    │   ├── repositories            # Database operations within the application
    │   ├── services                # It provides communication between "handlers" and "repositories" within the application and logical operations are performed here.
    │   └── worker                  # Automatic message sender worker
    │
    └── pkg                       # The codes, written with the intention that they can be used in other projects, are available here as a package.

---

## Swagger

I used my own package [swagno](https://github.com/go-swagno/swagno) for Swagger documentation.
SwaggerUI can be accessible from [http://localhost:8080/swagger/](http://localhost:8080/swagger/)

![resim](https://github.com/anilsenay/message-sending-system/assets/1047345/1d482dd5-0480-4b50-a835-7cbc040e07da)

---

### Diagrams

![resim](https://github.com/anilsenay/message-sending-system/assets/1047345/63940ed9-a69d-4591-8410-e1fb7cace16d)

![resim](https://github.com/anilsenay/message-sending-system/assets/1047345/a073ef2c-9ea5-47cc-a91e-5b830016306b)
