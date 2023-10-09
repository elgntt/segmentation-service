# Сервис динамического сегментирования пользователей

Сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент)

Используемые технологии:
- PostgreSQL (хранилище данных)
- Docker (запуск сервиса)
- Swagger (для документации API)
- Gin (веб-фреймворк)
- pgx (драйвер работы с PostgreSQL)
- golang/mock (для тестирования)

Сервис написан на Clean Architecture, добавление функционала и тестирование его не должно приводить к каким-либо проблемам

# 🔧Getting Started
Перед запуском нужно создать файл .env и заполнить его по шаблону .env.example

# 🚀Запуск 

Запуск сервиса осуществляется использованием команды `make compose-up`

После запуска по пути `http://localhost:8080/swagger/index.html` доступен swagger, где описаны ручки(если port не был изменён)

Запуск тестов командой `make test`, запуск тестов с покрытием `make cover` и для получения отчёта в html формате `make cover-html` 

## Examples

### Создание сегмента

Добавление сегмента в базу данных:

```curl
curl --location --request POST 'localhost:8080/segment' \
--header 'Content-Type: application/json' \
--data '{
    "slug": "VOICE_MESSAGE",
    "autoJoinProcent": 100
}'
```

Пример ответа: http-статус код: 201(Created)

### Удаление сегмента

Удаление сегмента из базы данных

```curl
curl --location --request DELETE 'localhost:8080/segment' \
--header 'Content-Type: application/json' \
--data '{
    "slug": "DISCOUNT_20"
}'
```
Пример ответа: http-статус код: 200(OK)

### Получение активных сегментов пользователя

Получение активных сегментов пользователя(id пользователя передаётся в URL(userId))

```curl
curl --location --request GET 'localhost:8080/user/segment/active?userId=347'
```

Пример ответа:
```json
{
    "userId": 31,
    "segments": [
        "TECH",
        "DISCOUNT_15",
        "DISCOUNT_11",
        "DISCOUNT_12",
        "DISCOUNT_13"
    ]
}
```

### Метод добавления и удаления юзера из сегмента

Метод добавляет и удаляет для юзера переданные в массиве сегменты. Если сегментов в базе не существует, отправится ошибка с массивом ошибочных сегментов
```curl
curl --location --request POST 'localhost:8080/user/segment/action' \
--header 'Content-Type: application/json' \
--data '{
    "userId": 347,
    "segmentsToAdd": ["DISCOUNT_12"],
    "segmentsToRemove": ["DISCOUNT_15"]
    
}'
```

Также принимает необязательный параметр в виде даты истечения времени(формате ISO 8601) жизни для всех добавляемых сегментов.
Пример запроса с указанием времени истечения:
```curl
curl --location --request POST 'localhost:8080/user/segment/action' \
--header 'Content-Type: application/json' \
--data '{
    "userId": 347,
    "segmentsToAdd": ["DISCOUNT_12"],
    "segmentsToRemove": ["DISCOUNT_15"],
    "expirationTime": "2023-08-31T22:18:10+03:00"
}'
```

Пример ответа: http-статус код: 200(OK)

### Метод получения истории по одному юзеру по указанному месяцу и году

Метод возвращает ссылку на сгенерированный CSV-файл. В URL передаются месяц, год и id юзера
```curl
curl --location --request GET 'localhost:8080/history/file?month=8&year=2023&userId=32123'
```

Пример ответа: 
```json
{
    "url": "http://localhost:8080/assets/csv_reports/0e666515-c657-4e49-b195-431c682563f7.csv"
}
```