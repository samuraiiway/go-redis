# go-redis

## Features
- Store key-value pair to memory
- Indexing key-value column
- Get value by key
- Get value by indexing
- Stremaing listener for new data

## Example

### Store data with indexing
**id can be null to auto generated with UUID schema.**

POST http://localhost:8000/redis/user
```json
{
    "id": "1234567890",
    "name": "song",
    "password": "12345678",
    "role": "admin",
    "index": ["name", "password", "role"]
}
```

### Get data by key
GET http://localhost:8000/redis/user/1234567890
```json
{
    "id": "1234567890",
    "name": "song",
    "password": "12345678",
    "role": "admin"
}
```

### Get data by index
GET http://localhost:8000/redis/user/role/admin
```json
[
    {
        "id": "1234567890",
        "name": "song",
        "password": "12345678",
        "role": "admin"
    }
]
```

### Stream listener
GET http://localhost:8000/redis/listen/user
```json
data: {"id":"1234567890","name":"song","password":"12345678","role":"admin"}

data: {"id":"1234567890","name":"song","password":"12345678","role":"admin"}
```
