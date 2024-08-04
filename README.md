# Avito-Rent

![Build Status](https://github.com/delapaska/avito-rent/actions/workflows/goci.yml/badge.svg)
[![Coverage Status](https://codecov.io/gh/delapaska/avito-rent/branch/master/graph/badge.svg)](https://codecov.io/gh/delapaska/avito-rent)

# О проекте 

Сервис, предназначенный для загрузки объявлений о продаже квартир 


### Запуск проекта 

Сначала необходимо склонировать репозиторий, выполнив команду: 

```
git clone https://github.com/delapaska/avito-rent.git
```

### Конфигурация 

Конфигурация реализована через `.env` файл, в репозитории есть файл `.env.example` в котором описаны необходимые переменные (можно просто скопировать), но я сделаю небольшое описание, для переменных, которые считаю важными.
`DB_HOST` должна быть одноимённа с именем контейнера бд, в данном случае это будет `avito-db`. 
`DB_NAME` - имя для базы можно выбрать любое, так как создан скрипт `init.db`, создающий базу данных, а также автоматически производятся миграции
`PORT` - базовый порт необходимо поставить `8080`





### Запуск сервиса 

1. Если у вас установлена утилита make, то необходимо выполнить следующие команды:
    - Сборка проекта: `make build`
    - Запуск проекта: 
        - Запуск с логами докера: `make run-logs`
        - Запуск без логов: `make run`
    - Завершение работы:
        - Остановка: `make stop`
        - Удаление: `make down`

2. Если утилита make отсутствует:
    - Сборка проекта: `sudo docker-compose build`
    - Запуск проекта: 
        - Запуск с логами докера: `sudo docker-compose up`
        - Запуск без логов: `sudo docker-compose up -d`
    - Завершение работы:
        - Остановка: `sudo docker-compose stop`
        - Удаление: `sudo docker-compose down`

### Endpoints
Далее указаны примеры маршрутов при запуске через docker-compose:
- noAuth: 
    - GET `localhost:8080/dummyLogin?userType=moderator`
    - POST ` localhost:8080/register`
    - JSON: 
         ```json
        {
            "email":"email@mail.ru", 
            "password":"secretKey", 
            "userType": "moderator"
        }
        ``` 
    - POST ` localhost:8080/login`
    - JSON: 
         ```json
        {
            "id": "b540a379-94ac-4eee-8c4e-83faf2f2d508", 
            "password":"12345"
        }
        ``` 
- authOnly:
    - GET `localhost:8080/house/1`
    - POST `localhost:8080/house/1/subscribe`
    - JSON: 
         ```json
        {
            "email":"email@mail.ru"     
        }
        ``` 
    - POST `localhost:8080/flat/create`
    - JSON: 
         ```json
        {
            "house_id": 4, 
            "price": 10000,
            "rooms": 4
        }
        ``` 
- moderatorsOnly: 
    - POST `localhost:8080/house/create`
    - JSON: 
         ```json
        {
            "address":"Лесная улица, 7, Москва, 125196", 
            "year":2003, 
            "developer": "Мэрия"
        }   
        ``` 
    - POST `localhost:8080/flat/update`
    - JSON: 
         ```json
        {
            "id":3, 
            "status":"on moderation"
        }
        ``` 



### Дополнения к решению 

Так как в документации я не нашёл описания, как правильно сделать ограничение модерации над квартирой с помощью `dummyLogin`, я решил сохранять UUID пользователя в токен, а не только его тип, что в дальнейшем работает и для эндпоинтов авторизации. Также я добавил поле с UUID в таблицу для квартир.

Реализована swagger документация, чтобы открыть её, перейдите по ссылке `localhost:8080/docs/index.html`, в ней описаны все эндпоинты и модели, включая как payload модели, так и основные модели.