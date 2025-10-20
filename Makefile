run:
	docker compose up

process-coupons:
	go run ./cmd/couponreader/... ;

seed:
	go run ./cmd/seeder/... ;

sql-gen:
	sqlc generate

mock-gen:
	go generate ./... ;

test:
	go test ./... ;
