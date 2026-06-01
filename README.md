# GophKeeper

GophKeeper — клиент-серверный менеджер секретов.

## Что умеет

- регистрация и логин
- JWT-авторизация
- хранение и синхронизация секретов через сервер
- локальное шифрование на клиенте мастер-паролем
- типы секретов:
    - login/password
    - text
    - binary
    - card
- build info в клиенте и сервере
- JSON-конфиг
- graceful shutdown сервера

## Архитектура

- `cmd/server` — HTTP-сервер
- `cmd/client` — CLI-клиент
- сервер хранит только шифротекст
- клиент шифрует и расшифровывает секреты локально
- мастер-пароль не отправляется на сервер

## Сборка

### Сервер

```bash
go build -ldflags="-X main.buildVersion=1.0.0 -X main.buildDate=2026-06-01 -X main.buildCommit=local" -o gophkeeper-server ./cmd/server