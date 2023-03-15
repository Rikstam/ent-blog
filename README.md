# Readme

## Overview
Based on tutorial from [https://entgo.io/blog/2023/02/23/simple-cms-with-ent/](https://entgo.io/blog/2023/02/23/simple-cms-with-ent/).
- Run DB with `docker run --rm --name entdb -d -p 3306:3306 -e MYSQL_DATABASE=ent -e MYSQL_ROOT_PASSWORD=pass mysql:8` 
- Run tests with `go test ./...`
- Start the app with `go run main.go -dsn "root:pass@tcp(localhost:3306)/ent?parseTime=true"`
- App visible at `localhost:8080`
## TODO
- CSRF token
- Input validations
- add HTMX
- add authentication