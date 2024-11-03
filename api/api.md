# VKatuny API

## Регистрация
<details>
<summary>Работодателя</summary>

URL: `/api/v1/registration/applicant`  
Метод: `POST`  
Для регистрации пользователя принимает JSON формата  
```json
{
    "firstName": "ivan",
    "lastName": "ivanov",
    "position": "hr",
    "company": "The best company ever",
    "companyDescription": "Really the best comapny",
    "companyWebsite": "bestcompany.com",
    "email": "best@mail.com",
    "password": "12345"
}
```
Может вернуть:  
1. `400` статус с телом ответа (текст ошибки может различаться)
```json
{
    "statusCode": 400,
    "body": null,
    "error": "can't unmarshall JSON"
}
```
2. `200` с телом
```json
{
    "statusCode": 400,
    "body": {
        "userType": "employer",
        "id": 25
    },
    "error": ""
}
```
</details>

<details>
<summary>Соискателя</summary>

URL: `/api/v1/registration/applicant`  
Метод: `POST`  
На вход принимает JSON  
```json
{
    "firstName": "ivan",
    "lastName": "ivanov",
    "birthDate": "01.01.2000",
    "email": "ivan_ivanov@mail.com",
    "password": "12345"
}
```  
Может вернуть:
1. `200` с телом  
```json
{
    "statusCode": 200,
    "body": {
        "userType": "applicant",
        "id": 123,
    },
    "error": ""
}
```   
2. `400` с телом (текст ошибки может различаться)  
```json
{
    "statusCode": 400,
    "body": null,
    "error": "user's fields aren't valid"
}
``` 
3. `500` с телом
```json
{
    "statusCode": 500,
    "body": null,
    "error": "can't add applicant to db"
}
```
</details>

## Авторизация  

<details>
<summary>Логин</summary>

URL: `/api/v1/login`  
Метод: `GET`???  
На вход принимает тип пользователя, эл. почту и пароль в форме JSON
```json
{
    "userType": "employer",
    "login": "email@email.com",
    "password": "strongest password"
}
```
Может вернуть:
1. `200` в случае успешного входа вместе с кукой  
2. `400` в случае проблем со входом (разный текст ошибок)
```json
{
    "statusCode": 400,
    "body": null,
    "error": "invalid fields" 
}
```

</details>


<details>
<summary>Проверка авторизации</summary>

Позволяет проверить авторизован пользователь или нет  
URL: `/api/v1/authorized`  
Метод: `GET`  
Вытаскивает из заголовков куку и проверяет существование сессии   
1. `200` если сессия существует
```json
{
    "statusCode": 200,
    "body": {
        "userType": "employer",
        "id": 2415
    },
    "error": ""
}
```
2. `401` если не удалось проверить авторизацию либо она неуспешна  
```json
{
    "statusCode": 401,
    "body": null,
    "error": "authorization error"
}
```
</details>

<details>
<summary>Логаут</summary>

Получает куку (если есть), удаляет ее из сессии и устанавливает истекший срок  
URL: `/api/v1/logout`  
Метод: `GET?`  
Принимает JSON:
```json
{
    "userType": "applicant"
}
```
Возвращает:  
1. `200` если куки не было в хедере запроса либо такой куки нет в сессиях
```json
{
    "statusCode": 200,
    "body": null,
    "error": "client doesn't have a session"
}
```
2. `200` Успешный логаут
```json
{
    "statusCode": 200,
    "body": null,
    "error": ""
}
```
3. `400`
```json
{
    "statusCode": 400,
    "body": null,
    "error": "can't unmarshall JSON"
}
```
</details>

## Вакансии

<details>
<summary>Вакансии</summary>

URL: `/api/v1/vacancies`  
Метод: `GET`  
Query-параметры: `offset` и `num` - натуральные числа (по дефолту 0 и 10 соответственно)  
Получает: `num` вакансий с отступом `offset`  
Возвращает:  
1. `200 OK` - возвращает вакансии
```json
{
    "statusCode": 200,
    "body": [
        {
            "id": 1,
            "position": "artist",
            "description": "looking for 3d artist",
            "salary": "100500",
            "location": "Moscow",
            "employer": "X5-Retail",
            "createdAt": "02.11.2024 20:10:24",
            "logo": "image.png"
        },
        // also vacancies
    ],
    "error": ""
}
```
2. `400` некорректные query-параметры  
```json
{
    "statusCode": 400,
    "body": null,
    "error": // текст описание ошибки
}
```
3. `500` если метод не GET
```json
{
    "statusCode": 500,
    "body": null,
    "error": "http request method isn't a GET"
}
```
3
</details>
