# Домашнее задание №7 "Хранение данных"


## Цель:

Модифицируйте Ваш проект, добавив следующий функционал:
- Реализуйте отдельный слой кеширования, в котором будут хранится данные чтобы уменьшить нагрузку на базу данных.
- Придумайте механизм инвалидации кеша.
- Запросы в базу данных должны использоваться в транзакциях, для поддержания целостности данных.

## Дополнительные задания (10 баллов)

- На потенциально нагруженных ручках добавьте использование memcached/redis.

### Дедлайн: 
- 20 апреля, 23:59 (сдача) / 23 апреля, 23:59 (проверка)



## GET:
curl.exe -X GET -u user:user http://localhost:9000/pvz/1
curl.exe -X GET -u user:user http://localhost:9000/pvz/2
curl.exe -X GET -u user:user http://localhost:9000/pvz/3
curl.exe -X GET -u user:user http://localhost:9000/pvz/4

## POST:
curl.exe -X POST -u user:user -H "Content-Type: application/json" -d @jsontest/post1.json http://localhost:9000/pvz 
curl.exe -X POST -u user:user -H "Content-Type: application/json" -d @jsontest/post2.json http://localhost:9000/pvz 
curl.exe -X POST -u user:user -H "Content-Type: application/json" -d @jsontest/post3.json http://localhost:9000/pvz 
curl.exe -X POST -u user:user -H "Content-Type: application/json" -d @jsontest/post4.json http://localhost:9000/pvz

## PUT:
curl.exe -X PUT -u user:user -H "Content-Type: application/json" -d @jsontest/put1.json http://localhost:9000/pvz

## DELETE
curl.exe -X DELETE -u user:user http://localhost:9000/pvz/1