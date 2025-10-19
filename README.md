# Kart Challenge
## Steps To Run
```sh
# build the project includes server, database and migration
make run

# seed the database with some basic data for the endpoints to work
make seed

# WIP duplicate keys within the same file are causing batches to fail
make process-coupons
```

## Changes
### Incrementing IDs
Using UUIDs instead of sequential IDs

### API Key 
There is too little context around what the api_key aims to achive. If it is front end facing it should be a bearer token. If it is for Third party services no valid or invalid keys have been provided

### Invalid Coupon
Returns error on invalid coupon.

### Idempotency
Very basic implementation to avoid accidental duplicate orders

## Notes
Due to time constraints several things have not been implemented

### Missing Integration Test
Ideally would have used testcontainers to test the store layer

### Pagination
Ideally pagination would be included for the ListProducts endpoint

## Requests

GetProductByID
```sh
curl --header "Content-Type: application/json" \
http://localhost:8080/api/v1/product/00000000-0000-0000-0000-000000000001

```

GetProduct
```sh
curl --header "Content-Type: application/json" \
http://localhost:8080/api/v1/product

```

CreateOrder
```sh
curl http://localhost:8080/api/v1/order \
  --request POST \
  --header 'Content-Type: application/json' \
  --header 'api_key: YOUR_SECRET_TOKEN' \
  --header 'Idempotency-Key: YOUR_ORDER_KEY' \
  --data '{
  "coupon_code": "FIFTYOFF",
  "items": [
    {
      "product_id": "00000000-0000-0000-0000-000000000001",
      "quantity": 1
    },
    {
      "product_id": "00000000-0000-0000-0000-000000000002",
      "quantity": 2
    }
  ]
}'
```
