Записка по реализации [тут](https://github.com/sayanli/calculator/note.md)

- сервис поднимается через docker-compose;
- тестами покрыто 27% кода (сервисный слой на 98%, httpserver на 57%);
- реализован grpc и http endpoints;
- есть swagger документация (можно посмотреть по адресу http://localhost:8080/swagger/index.html);


**Команды для запуска**

```
make run // локальный запуск без докера

make compose-build // сборка проекта
make compose-run // запуск проекта

make test // запустить тесты
make cover-html // запустить тесты с покрытием и html отчет
make cover // запустить тесты с покрытием

make proto-gen // генерация proto файлов
make swag // генерация swagger документации
```


**Примеры запросов**

Запрос
```
curl --location --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '[{"type": "calc", "op": "+", "var": "x", "left": 1, "right": 2 }, \
{"type": "print", "var": "x" }]'
```
Ответ
```json
{
  "item": [
    {"Var": "x", "Value": 3}
  ]
}
```

Запрос с ошибкой

```
curl --location --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '[{"type": "calc", "op": "+", "var": "x", "left": "unknown", "right": 2 }, \
{"type": "print", "var": "x" }]'
```
Ответ
```json
{
  "Code": 400,
  "Message": "Unknown variable",
  "Timestamp": "2025-03-20T12:34:56Z"
}
```
