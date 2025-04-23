# How to Deploy Weshnet Infrastructure

This guide explains how to set up the essential services for a Weshnet network infrastructure, including rendez-vous points, emitter.io, and relay services.

## Prerequisites

Before starting the deployment, ensure you have the following tools installed on your system:

1. **Docker**: Required to run all services in containers
   - For Ubuntu/Debian: `sudo apt update && sudo apt install docker.io docker-compose`
   - For macOS: Download and install Docker Desktop from [docker.com](https://www.docker.com/products/docker-desktop)
   - For Windows: Download and install Docker Desktop from [docker.com](https://www.docker.com/products/docker-desktop)

2. **Make**: Used to simplify deployment commands
   - For Ubuntu/Debian: `sudo apt install make`
   - For macOS: Install Xcode Command Line Tools with `xcode-select --install`
   - For Windows: Install via [Chocolatey](https://chocolatey.org/) with `choco install make`

Verify installations with `docker --version` and `make --version`.

## Service Overview

Weshnet relies on three main components to facilitate peer-to-peer communication:

1. **Rendez-vous Point (RDVP)**: Acts as a meeting point for peers to discover each other on the network. It helps peers establish connections without needing to know each other's exact network location beforehand.

2. **Emitter.io**: Provides a pub/sub messaging system that allows peers to broadcast their presence and receive notifications about other peers. This service facilitates real-time communication and discovery within the network.

3. **Relay Service**: Helps peers connect when they're behind NATs or firewalls. It relays traffic between peers that cannot establish direct connections, ensuring connectivity even in challenging network environments.

## Rendez-vous Point and Emitter.io Services

### Setting Up Rendez-vous Point

1. Generate a new private key for the rendez-vous point service:
   ```sh
   cd rdvp
   docker run --rm --entrypoint rdvp bertytech/berty:kubo-v0.29.0 genkey
   ```

2. When the command completes, you'll receive a key in this format:
   `CAESQHW91QjcGJN1RrIXtzCf8aC5EHCIB2Q+CSJ6KI68E7WLn49INScVKtToDjCMk4TxnncKWFcys59TjCgu8yBDOD8=`

3. Copy this key to the `RDVP_PK` variable in your `.env` file.
   
4. Add your public IP address to the `ANNOUNCE_SERVER` variable in the same file.

### Setting Up Emitter.io

1. Generate a license and secret key for emitter.io:
   ```sh
   cd rdvp
   docker run --rm emitter/server:v3.1
   ```

2. From the output, copy:
   - The license to the `EMITTER_LICENSE` variable
   - The secret key to the `EMITTER_SECRET_KEY` variable in your `.env` file

### Starting the Services

Once your configuration is complete, start both services with:

```sh
make up
```

## Relay Service

The relay service facilitates peer connections through NATs and firewalls.

### Configuration and Deployment

1. Edit the relay configuration file:
   - Open `relay/config.json`
   - Update the `Network/AnnounceAddrs` section with your public IP address

2. Deploy your relay:
   ```sh
   cd relay
   make build  # Build the relay Docker image
   make up     # Start the relay service
   ```

## Verifying Your Deployment

After deployment, you can verify your services are running correctly by checking:
- Logs for each service
- Network connectivity through the announced addresses
- Peer connections through your infrastructure

## Troubleshooting

If you encounter issues:
- Check that ports are properly forwarded on your network
- Verify your public IP is correctly configured in all services
- Ensure Docker has sufficient resources allocated
- Review service logs for specific error messages
