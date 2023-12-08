в своем проекте переименовать `.env copy` в `.env`

`docker-compose up -d` либо создать локально постгрес например через pgAdmin4 и в env указать свои параметры

Если без докера то удалить файл `docker-compose`

`go run main.go`

Используем структурированные логи `https://pkg.go.dev/go.uber.org/zap`

адаптер к postgresql `pgx/v5`

библиотека для чтения переменных из `.env` godotenv


# curl запросы для сервиса

# Необходимо вводить все поля
`curl --request POST \
  --url http://0.0.0.0:9000/create \
  --data '{"name":"Maxim","last_name":"Korolev","middle_name":"Alexevich","address":"Baumanskaja","phone":"8 800 535 35 35"}'`

# Минимум телефон, остальное  сделает поиск уже
`curl --request POST \
  --url http://0.0.0.0:9000/get \
  --data '{"name":"Maxim","last_name":"Korolev","middle_name":"Alexevich","address":"Baumanskaja","phone":"8 800 535 35 35"}'`

# Меняет данные по номеру телефона
`curl --request POST \
  --url http://0.0.0.0:9000/update \
  --data '{"name":"Maxim","last_name":"Korolev","middle_name":"Alexevich","address":"Malaia Pochtovaia 5/12 kv. 159","phone":"+79190771930"}'`

# Удаляет данные по номеру телефона
`curl --request POST \
  --url http://0.0.0.0:9000/delete \
  --data '{"phone":"8800535 35 35"}'`