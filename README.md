## JWT auth service

для тестирования можно воспользоваться [полноценно функционирующим сервисом](https://jwt-auth-4tmd.onrender.com/), размещенным на [render.com](https://render.com)

сервис бесплатный, могут быть задержки в ответе

имеется три маршрута:

* POST /register - принимает в body email и password
в ответе примерно так:
```json
{
    "guid": "3ff3e273-5dda-4075-9f74-a04f6c8bb5ad",
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6IiwxODguNDMuMTEzLjIyNiwgMTcyLjcxLjI0Ni4xMjcsIDEwLjIyMy4xNjIuNjIiLCJzdWIiOiIzZmYzZTI3My01ZGRhLTQwNzUtOWY3NC1hMDRmNmM4YmI1YWQiLCJleHAiOjE3NDUwNjE2Mzh9.vSKAM_f6NUQ-3ZcYgP2kauqRrG7od8frDnbPRANCm2HMCUBS1Twur2AolYqJnxfjAeVWUVt9Y8aqH6M9SevuqQ",
    "refresh_token": "gzuXVHn65ly+MgfJ0vPE5CwbvMH2htgnWtlGxVhl+dh7V9AYxauJWPMES0BCrrgyAKT12prarvDKhDRUIFhEt63rO7aEe5DU",
    "expires_at": "2025-04-19T11:20:38.917609199Z"
}
```
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