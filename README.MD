# Board
Приложение состоит из двух частей
1. Бэкенд на GO
2. Фронт на Vue3

# Билд
1. ```docker compose exec -it node bun run build``` - Билд фронта
2.
```docker compose exec -it app go run cmd/board/main.go``` - Исполнение приложения

или

```docker compose exec -it app go build cmd/board/main.go ``` - Билд бинарного файла