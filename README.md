## JWT auth service

для тестирования можно воспользоваться [полноценно функционирующим сервисом](https://jwt-auth-4tmd.onrender.com/swagger), размещенным на [render.com](https://render.com)

сервис бесплатный, могут быть задержки в ответе

имеется три маршрута:

* POST /register - принимает в body email и password
в ответе примерно так:
```json
{
    "guid": "89cb6314-5adf-4255-8baa-8bd3999d4623",
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6ImV4YW1wbGUgaXAgYWRkcmVzcyIsInN1YiI6Ijg5Y2I2MzE0LTVhZGYtNDI1NS04YmFhLThiZDM5OTlkNDYyMyIsImV4cCI6MTc0NTA3MTA5NX0.G9DOOsy8qlTdH4WW9cMThq_r9RRj71cd-SYWeCXb7BGXyWvTGj6WYPFuriCGljq9EI3VPHsI33HMUfDO6n_pEQ",
    "refresh_token": "KJWGq0/n8Og9639c2NRhCZBRZYu9J5M7Be+lZmFTiju9KW3HvwPCx9PW0qz9tXR5mXGNOhP1iLkKlhbIR7vZnsuaeSUXnZPs",
    "expires_at": "2025-04-19T13:58:15.605241076Z"
}
```

<img width="1154" alt="Снимок экрана 2025-04-19 в 15 59 16" src="https://github.com/user-attachments/assets/add75079-9cf4-4684-8925-c4b51c95f486" />

expires_at: время до которого действителен access token

* GET /tokens - принимает в параметре запрос guid и выдает новую пару access и refresh токенов

>ответ выглядит аналогично

* POST /refresh - принимает в параметре bosy пару access и refresh токенов и выдает новую пару access и refresh токенов

>ответ выглядит аналогично

### Тестовое задание на позицию Junior Backend Developer

<details>
<summary>смотреть подробнее описание задачи</summary>

Используемые технологии:

* Go
* JWT
* PostgreSQL

Задание:

Написать часть сервиса аутентификации.

Два REST маршрута:

* Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
* Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

Требования:

Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.

Refresh токен тип произвольный, формат передачи base64, хранится в базе исключительно в виде bcrypt хеша, должен быть защищен от изменения на стороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним.

Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан. В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные).

Результат:

Результат выполнения задания нужно предоставить в виде исходного кода на Github. Будет плюсом, если получится использовать Docker и покрыть код тестами.

P.S. Друзья! Задания, выполненные полностью или частично с использованием chatGPT видно сразу. Если вы не готовы самостоятельно решать это тестовое задание, то пожалуйста, давайте будем ценить время друг друга и даже не будем пытаться :)
</details>