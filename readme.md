Вот пример README.md для твоего проекта с описанием и тестовыми маршрутами:

⸻

###TODO
- add postgis db
- add GET /heat_map -> all gps point for heat map


# IP Info Server

Простой REST API-сервер на Go, который позволяет пользователю:

- Зарегистрироваться и получить токен
- Авторизоваться по логину и паролю и получить токен
- Получить информацию о своём IP
- Получить информацию о любом IP
- Получить историю запросов

## 📦 Запуск

```bash
go run cmd/ipinfo_server/main.go

По умолчанию сервер стартует на http://localhost:8080.

🗂️ Эндпоинты

{{domain}} = http://localhost:8080
{{auth_token}} = токен, полученный при логине или регистрации

⸻

🔐 POST /register

Зарегистрировать нового пользователя.
Если пользователь с таким логином уже существует, будет ошибка.

Пример запроса:

POST {{domain}}/register?username=username&password=password



⸻

🔑 GET /login

Выполнить вход и получить токен авторизации.

Пример запроса:

GET {{domain}}/login?username=username&password=password

Ответ:

{
  "token": "abc123..."
}



⸻

🌐 GET /self_ip

Получить информацию о своём IP-адресе (по заголовку запроса).

Требуется авторизация.

Пример запроса:

GET {{domain}}/self_ip
Authorization: Bearer {{auth_token}}



⸻

🌍 GET /ext_ip/{ip}

Получить информацию о внешнем IP-адресе.

Требуется авторизация.

Пример запроса:

GET {{domain}}/ext_ip/3.3.3.3
Authorization: Bearer {{auth_token}}



⸻

🕘 GET /history

Получить историю IP-запросов текущего пользователя.

Требуется авторизация.

Пример запроса:

GET {{domain}}/history
Authorization: Bearer {{auth_token}}



⸻

🛠 Зависимости
	•	Go >= 1.20
	•	SQLite
	•	chi

⸻

📁 Структура проекта

ip_info_server/
├── cmd/
│   └── ipinfo_server/
│       └── main.go         # Точка входа
├── internal/
│   ├── db/                 # Инициализация БД и функции доступа
│   ├── handlers/           # HTTP-обработчики
│   └── middleware/         # Middleware для авторизации
└── ipinfo.db               # SQLite база данных



⸻

🧪 Примеры Postman / curl

# Регистрация
curl -X POST "http://localhost:8080/register?username=user&password=12345"

# Логин
curl "http://localhost:8080/login?username=user&password=12345"

# Self IP
curl -H "Authorization: Bearer {{auth_token}}" http://localhost:8080/self_ip

# External IP
curl -H "Authorization: Bearer {{auth_token}}" http://localhost:8080/ext_ip/8.8.8.8

# History
curl -H "Authorization: Bearer {{auth_token}}" http://localhost:8080/history



⸻

📌 Примечание

Токен у пользователя сохраняется один раз при регистрации и используется повторно при логине. Повторная регистрация создаст нового пользователя с новым токеном.

