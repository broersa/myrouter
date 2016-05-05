#! /bin/bash
docker run --name myrouterpg -e POSTGRES_PASSWORD=myrouterpg -d postgres

until docker run -e PGPASSWORD=myrouterpg --rm --link myrouterpg:postgres postgres psql -h postgres -U postgres; do
    echo "Waiting for janky-postgres container..."
    sleep 0.5
done

docker run --rm --link myrouterpg:postgres -v `pwd`/sql:/flyway/sql shouldbee/flyway -url=jdbc:postgresql://postgres/postgres -user=postgres -password=myrouterpg migrate

docker build -t myrouter .

docke run -rm -e MYLOGLEVEL=TRACE -e MYROUTER_DB="postgres://postgres:myrouterpg@host?sslmode=disable" --link myrouterpg:host -publish 1700:1700 -name myrouter -d myrouter
