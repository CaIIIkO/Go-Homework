## Домашнее задание №6 «Логирование обращений в методы через Kafka»
### Цель:

Вывести события по обращениям в методы сервиса в консоль

### Задание:

1. При вызове любого API метода записываем событие об этом в кафку
2. Событие должно хранить время обращения, название метода и сырой запрос
3. По ходу работы приложения в консоль необходимо выводить эти события из кафки
4. События нельзя сразу выводить в консоль, минуя кафку
5. Решение должно быть покрыто unit тестами

### Дополнительное задание (за него можно получить 10 баллов):

- Написать интеграционные тесты на взаимодействие с кафкой

### Дедлайны
13 апреля, 23:59 (сдача) / 16 апреля, 23:59 (проверка)


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