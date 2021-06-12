# password methods
All <span style="color:red">red</span> vars are necessary and all <span style="color:yellow">yellow</span> vars are optional. <span style="color:orange">*need token</span> means that you need to provide token in "Authorization" header.

<p style="font-size: 20pt"><b>Table of contents</b></p>

- [Create](#create): Creates new password.
- [Get password](#read): Returns password's info.
- [Update password](#update): Updates password's info.
- [Delete password](#delete): Deletes password.


## <a name="create"></a> `POST` collection/{title}/
<span style="color:orange">*need token</span><br/>
Creates new password.
### body:
- <span style="color:red">title</span>: your title for password (must be more than 2 symbols and less than 25 symbols)
- <span style="color:red">login</span>: your login for password (must be more than 1 symbol and less than 128 symbols)
- <span style="color:red">password</span>: your password (must be more than 1 symbol and less than 2048 symbols)
- <span style="color:yellow">email</span>: your title for password (must be less than 64 symbols)

### response:
nothing

### status code: 201
---

## <a name="read"></a> `GET` collection/{title}/{id}
<span style="color:orange">*need token</span><br/>
Returns password info.

### response:
- title: Title for password.
- login: Login for password.
- email: Email for password.
- password: Password.

### status code: 200
---

## <a name="update"></a>`PUT` collection/{title}/{id}
<span style="color:orange">*need token</span><br/>
Updates password's info.
### body:
- <span style="color:yellow">title</span>: your title for password (must be more than 2 symbols and less than 25 symbols)
- <span style="color:yellow">login</span>: your login for password (must be more than 1 symbol and less than 128 symbols)
- <span style="color:yellow">password</span>: your password (must be more than 1 symbol and less than 2048 symbols)
- <span style="color:yellow">email</span>: your title for password (must be less than 64 symbols)

### response:
nothing

### status code: 200
---

## <a name="delete"></a> `DELETE` collection/{title}/{id}
<span style="color:orange">*need token</span><br/>
Deletes password.

### response:
nothing

### status code: 200
