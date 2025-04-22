# How to deploy weshnet infrastructure

## Rendez-vous point and emitter.io services

Execute the following command to generate a new private key for the rendez-vous point service:

```sh
cd rdvp
docker run --rm --entrypoint rdvp bertytech/berty:kubo-v0.29.0 genkey
```

Copy this key to the `RDVP_PK` variable in the `.env` file and add your public IP address in `ANNOUNCE_SERVER`.

We also use [emitter.io](https://emitter.io/) as a discovery service.

You have to generate a license and a secret key copy them to the `.env` file. Execute the following command:

```sh
cd rdvp
docker run --rm emitter/server:v3.1
```

Copy the license to the `EMITTER_LICENSE` variable and the secret key to the `EMITTER_SECRET_KEY` variable in the `.env` file.

Execute the following command to start the services:

```sh
make up
```

## Relay service

The relay service helps to connect peers behind NATs.

Firstly, edit the configuration file in `relay/config.json` with your public IP address in the `Network/AnnounceAddrs` section.

To deploy your own relay:

```sh
cd relay
make build # build the relay image
make up
```
