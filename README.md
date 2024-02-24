# book_management_api_with_gqlgen



create migration:
```
	migrate create -ext sql -dir migrations -seq create_books_table
```

migrate-up:
```
	migrate -database $(DATABASE_URL) -path migrations up
```

migrate-down:
```
	migrate -database $(DATABASE_URL) -path migrations down
```