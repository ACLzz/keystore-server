# collection methods
All <span style="color:red">red</span> vars are necessary and all <span style="color:yellow">yellow</span> vars are optional

<p style="font-size: 20pt"><b>Table of contents</b></p>

- [Create](#create): Creates new collection.
- [List collections](#read): Returns current user's collections list.
- [List passwords](#fetch): Returns passwords list that belongs to collection.
- [Update collection](#update): Updates collection's info.
- [Delete collection](#delete): Deletes collection.


## <a name="create"></a> `POST` collection/
Creates new collection.
### body:
- <span style="color:red">title</span>: your title for collection (must be more than 2 symbols and less than 40 symbols. Ascii-only)
- <span style="color:red">token</span>:your temporary token.

### response:
nothing

### status code: 201

## <a name="read"></a> `GET` collection/
Returns current user's collections list.
### body:
- <span style="color:red">token</span>: your temporary token.

### response:
array of strings

### status code: 200

## <a name="fetch"></a> `GET` collection/{title}
Returns passwords list that belongs to collection.
### body:
- <span style="color:red">token</span>: your temporary token.

### response:
array of:
- title: password's title.
- id: password's id, you will use it to do things with passwords.

### status code: 200

## <a name="update"></a>`PUT` collection/{title}
Updates collection's info.
### body:
- <span style="color:red">token</span>: your temporary token.
- <span style="color:yellow">title</span>: your title for collection (must be more than 2 symbols and less than 40 symbols. Ascii-only)

### response:
nothing

### status code: 200

## <a name="delete"></a> `DELETE` collection/{title}
Deletes collection.
### body:
- <span style="color:red">token</span>: your temporary token.

### response:
nothing

### status code: 200