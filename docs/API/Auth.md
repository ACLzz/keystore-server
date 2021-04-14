# auth methods
All <span style="color:red">red</span> vars are necessary and all <span style="color:yellow">yellow</span> vars are optional

<p style="font-size: 20pt"><b>Table of contents</b></p>

- [Register](#create): Register account on keystore server.
- [Login](#login): Sign in and returns temporary token.
- [Read user](#read): Returns current user info.
- [Update user](#update): Updates current user info.
- [Delete user](#delete): Deletes current user.


## <a name="create"></a> `POST` auth/
Register account on keystore server.
### body:
- <span style="color:red">login</span>: your login (must be more than 3 symbols and less than 20 symbols. Ascii-only)
- <span style="color:red">password</span>: your password (must be more than 8 symbols and less than 25 symbols. Ascii-only)

### response:
nothing

### status code: 201

## <a name="login"></a>`POST` auth/login
Sign in and returns temporary token.
### body:
- <span style="color:red">login</span>: your login (must be more than 3 symbols and less than 20 symbols. Ascii-only)
- <span style="color:red">password</span>: your password (must be more than 8 symbols and less than 25 symbols. Ascii-only)

### response:
- token: temporary token, you will send it in many requests.

### status code: 202

## <a name="read"></a> `GET` auth/
Returns current user info.
### body:
- <span style="color:red">token</span>: your temporary token.

### response:
- registered: Datetime when account was registered.
- username: Username of account.

### status code: 200

## <a name="update"></a>`PUT` auth/
Updates current user info.
### body:
- <span style="color:red">token</span>: your temporary token.
- <span style="color:yellow">login</span>: your login (must be more than 3 symbols and less than 20 symbols. Ascii-only)
- <span style="color:yellow">password</span>: your password (must be more than 8 symbols and less than 25 symbols. Ascii-only)

### response:
nothing

### status code: 200

## <a name="delete"></a> `DELETE` auth/
Deletes current user.
### body:
- <span style="color:red">token</span>: your temporary token.

### response:
nothing

### status code: 200