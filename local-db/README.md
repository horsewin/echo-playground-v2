# Local database for development

```shell
docker compose up -d
```

## Access

Replace `${MYSQL_USER_NAME}` and `${MYSQL_USER_PASSWORD}` with the values from `docker-compose.yml`.

```shell
MYSQL_USER_NAME=$(yq e ".services.db.environment.MYSQL_USER" docker-compose.yml)
MYSQL_USER_PASSWORD=$(yq e ".services.db.environment.MYSQL_PASSWORD" docker-compose.yml)
docker compose exec db mysql -u ${MYSQL_USER_NAME} -p${MYSQL_USER_PASSWORD}
```