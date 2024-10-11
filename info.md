### Notes 

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

```
migrate -path ./migrations -database "postgres://username:password@localhost:5432/mydb?sslmode=disable" up
```

docker inspect dcd | grep POSTGRES_PASSWORD
