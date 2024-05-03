docker run --rm -d --name pg-products -e POSTGRES_DB=productsdb -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres

docker exec -it pg-products psql -U postgres -d productsdb //to run postgres
