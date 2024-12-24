# gobny
E-commerce API in go

# Getting started

clone the project

```bash
git clone https://github.com/musllim/gobny.git
```
Generate source code from SQL

```bash

sqlc generate

```
Run the project

Create database schemas
1. With docker

```bash
cat internal/database/schema.sql | docker container exec -i {container_name} psql -U postgres
```
provide required environmental variables

```bash
JWT_SECRET={secret} DB_URL={uri} go run cmd/gobny/main.go
```
