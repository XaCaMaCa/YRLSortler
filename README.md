# YRLSortler
Учебный **сервис сокращения ссылок** на **Go**: короткий алиас ведёт на исходный URL.
**API (HTTP):**
- `POST /url/` — создать короткую ссылку (**Basic Auth**)
- `GET /{alias}` — редирект на оригинальный URL
- `DELETE /{alias}` — удалить алиас

**Стек:** [Chi](https://github.com/go-chi/chi), SQLite, `slog`, конфиг через env/config, интеграционные тесты (**httpexpect**), CI в `.github/workflows`.
**Зачем:** потренировать типичный небольшой backend-сервис (маршруты, middleware, хранилище, авторизация на запись) без лишней магии.
