# tucows-grill-client

## Overview
This is a client that takes REST requests, and sends REST requests to the `tucows-grill-service`. 

## Endpoints
The following endpoints are supported:  
`POST /login`             login to the client account
`GET  /ingredients/{id}`  get an ingredient by the ingredient id  
`POST /ingredients`       add a new ingredient to the db  
`GET  /total-cost"`       get the total cost of all purchased items for an item id. Accepts a flag `async` which value can be `true` or anything else (false). This determines if we will calculate the total cost synchronously or asynchronously.

## Setup instructions
1. run the `tucows-grill-service` locally  
2. I used vscode, so the launch.json is already configured with a build. Run the `main.go` configuration.  
3. hit the login endpoint. Copy that JWT, and paste it in the token input box in postman using the Authorization tab, with Auth Type `Bearer`. Or you can add it as a header. All requests are authenticated via this JWT. Currently no authorization is implemented, as root has all access.

## Login credentials
This is the body for the login POST request.
```
{
    "username": "root",
    "password": "root"
}
```

## cURL commands
Here's a list of curl commands that you can use to test the endpoints. Replace the JWT with your new JWT from loggin in. In this repo is an exported postman `tucows.postman_collection.json` file which is also available for use.  

### Login
```
curl --location 'http://localhost:8081/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "root",
    "password": "root"
}'
```

### Get ingredient by id
```
curl --location 'http://localhost:8081/ingredients/1' \
    --header 'Authorization: Bearer <JWT>'
```

### Create ingredient
```
curl --location 'http://localhost:8081/ingredients' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT>' \
--data '{
    "name": "Rotting Cow",
    "quantity": 123
}'
```

### Get total cost
```
curl --location 'http://localhost:8081/total-cost?item_id=1' \
    --header 'Authorization: Bearer <JWT>'
```

### Get total cost async
```
curl --location 'http://localhost:8081/total-cost?item_id=1&async=true' \
    --header 'Authorization: Bearer <JWT>'
```
