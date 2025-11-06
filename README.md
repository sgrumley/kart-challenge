# Kart Challenge
## Steps To Run
```sh
# build the project includes server, database and migration
make run

# seed the database with some basic data for the endpoints to work
make seed

# the coupon files need to be unzipped into a root folder named coupons/
# preprocess the data by removing duplicates within each file
# in ./coupons/ run:
sort -u couponbase1 > ucouponbase1
sort -u couponbase2 > ucouponbase2
sort -u couponbase3 > ucouponbase3

# then from root dir
make process-coupons
```

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
