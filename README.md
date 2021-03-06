# Keystore
## Server side
<p align="center">
    <img src="docs/logo.png?raw=true" width="200px"/>
</p>

## What is Keystore?
It is server side of keystore -- simple self-hosted key storage
You can run server on your local network and don't be afraid of password leaks.

## How to install it?
Firstly, I highly recommend you go to <a href="docs/Configuration.md">config</a> section before install, cause config effects on installation.
Just install go and build server in few seconds with make.

```
~ » sudo pacman -S go git make postgresql --noconfirm
~ » sudo systemctl start postgresql
~ » git clone https://github.com/ACLzz/keystore-server.git
~ » make setup_bin
```

If all ok, binary must be at `bin` folder.
!!! Don't forget to copy .env file (./bin/.env) if you will move binary of server

## Docker
You can use docker image for keystore, just run with root:
```
~ » make setup_docker
```

## How to use it?
Keystore has an http API, so all communications goes through that API.<br/>
Keystore has official client <a href="https://github.com/ACLzz/keywarden">Keywarden</a><br/>
For a documentation about API methods go <a href="docs/API/">here</a>.

## Logs
Logs are stores in `~/.kss/logs` folder
