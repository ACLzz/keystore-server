# Configuration guide
In this document I will explain all config options

## keystore server folder
Configs and server's, clients' certificates are stored in `~/.kss` folder.<br/>
So if you need to edit configs, refer to this folder.

## config.yml and it's options
Below I will explain each option mission
- addr: Server address. It uses 127.0.0.1 by default.
- port: Choose free port for a server.
- db_host: IP address of postgres database.
- db_port: Port of postgres database.
- db_name: Name of database, by default it is "keystore".
- db_username: Username of role which will make changes in database.
- db_password: Password for role which will make changes in database. Don't forget to CHANGE it.
- timezone: Your server's timezone.
- token_lifetime: Count of seconds while token will be valid.
- allow_registration: `true` or `false` for enable or disable registration.
- salt: Additional string for hash creation, uses in tokens creation. You MUST change it to anything else.
  <br/>
