# go-microservice
Use golang build microservice

[Reference Video](https://www.youtube.com/playlist?list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_) 

## Request Method 

GET:
```curl
curl localhost:9000 | jq
```
```json
[
  {
    "id": 1,
    "name": "Latte",
    "description": "Frothy milky coffee",
    "price": 2.45,
    "sku": "abc323"
  },
  {
    "id": 2,
    "name": "Espresso",
    "description": "Short and strong coffee without milk",
    "price": 1.99,
    "sku": "fid34"
  }
]

```
POST: 
```curl
curl localhost:9000  -d "{'id': 3, 'name': 'tea', 'descrption': 'a nice couple of tea'}"
```
PUT
```curl
curl -v localhost:9000/3 -XPUT  -d '{"id": 1, "name": "xx", "description": "a nice couple of team"}'
```
DELETE
```curl
curl -v localhost:9000/2 -XDELETE
```
