Для миграции используется goose.

Для деплоя:

Установить goose: go install github.com/pressly/goose/v3/cmd/goose@latest
make start
make migrate-up
Postman коллекция: test-note.postman_collection.json
