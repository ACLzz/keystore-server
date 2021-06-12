# collection methods
All <span style="color:red">red</span> vars are necessary and all <span style="color:yellow">yellow</span> vars are optional. <span style="color:orange">*need token</span> means that you need to provide token in "Authorization" header.

<p style="font-size: 20pt"><b>Table of contents</b></p>

- [Create](#create): Creates new collection.
- [List collections](#read): Returns current user's collections list.
- [List passwords](#fetch): Returns passwords list that belongs to collection.
- [Update collection](#update): Updates collection's info.
- [Delete collection](#delete): Deletes collection.


## <a name="create"></a> `POST` collection/
<span style="color:orange">*need token</span><br/>
Creates new collection.
### body:
- <span style="color:red">title</span>: your title for collection (must be more than 2 symbols and less than 40 symbols. Ascii-only)

### response:
nothing

### status code: 201
---

## <a name="read"></a> `GET` collection/
<span style="color:orange">*need token</span><br/>
Returns current user's collections list.

### response:
array of strings

### status code: 200
---

## <a name="fetch"></a> `GET` collection/{title}
<span style="color:orange">*need token</span><br/>
Returns passwords list that belongs to collection.

### response:
array of:
- title: password's title.
- id: password's id, you will use it to do things with passwords.

### status code: 200
---

## <a name="update"></a>`PUT` collection/{title}
<span style="color:orange">*need token</span><br/>
Updates collection's info.
### body:
- <span style="color:yellow">title</span>: your title for collection (must be more than 2 symbols and less than 40 symbols. Ascii-only)

### response:
nothing

### status code: 200
---

## <a name="delete"></a> `DELETE` collection/{title}
<span style="color:orange">*need token</span><br/>
Deletes collection.

### response:
nothing

### status code: 200
