FROM gomicro/goose:2.6.1 as goose

ENV MIGRATION_PATH /migration

WORKDIR /migrations

ADD migrations /migrations

CMD goose postgres "postgres://$DB_USER:$DB_PASS@$DB_HOST/$DB_NAME?sslmode=disable" up
