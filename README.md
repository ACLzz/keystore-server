# Keystore
## Server side
<p align="center">
    <img src="docs/logo.png?raw=true" width="200px"/>
</p>

## What is Keystore?
It is server side of keystore -- simple self-hosted key storage
You can run server on your local network and don't be afraid of password leaks.

## How to install it?
Firstly, I highly recommend you go to <a href="docs/Configuration.md">config</a> section before install.
Just install go and build server in few seconds with make.

```
~ » sudo pacman -S go git make postgresql --noconfirm
~ » sudo systemctl start postgresql
~ » git clone https://github.com/ACLzz/keystore-server.git
~ » make setup
```

If all ok, binary must be at `bin` folder

## How to use it?
Keystore has an http API, so all communications goes through that API.<br/>
Keystore don't have any clients yet, but if you familiar with http and some programming language you can write it!<br/>
For a documentation about API methods go <a href="docs/API/">here</a>.

## Logs
Logs are stores in `~/.kss/logs` folder