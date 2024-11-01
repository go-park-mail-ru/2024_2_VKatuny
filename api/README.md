### Как запустить генератор Swagger'а

Установка зависимостей
```
npm install express
```
```
npm install swagger
```

В корне проекта выполнить команду, в директории /api появится swagger файл
```
swag init --parseInternal --pd --dir cmd/myapp/,delivery/handler/ --output api/
```
Запустить интерфейс 
```
node server.js
```
