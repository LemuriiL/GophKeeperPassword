# GophKeeper

GophKeeper — клиент-серверный менеджер секретов.

## Возможности

* регистрация пользователя
* вход пользователя
* JWT-авторизация
* хранение секретов на сервере
* синхронизация секретов между клиентами одного пользователя
* локальное шифрование секретов на клиенте мастер-паролем
* хранение на сервере только зашифрованных данных
* TLS для соединения клиента и сервера
* JSON-конфиги клиента и сервера
* graceful shutdown сервера
* build info у клиента и сервера

Поддерживаемые типы секретов:

* `login_password`
* `text`
* `binary`
* `card`

## Архитектура

* `cmd/server` — серверное приложение
* `cmd/client` — CLI-клиент
* `internal/server` — серверная логика
* `internal/client` — клиентская логика
* `internal/secure` — шифрование и derivation key
* `internal/cfg` — загрузка JSON-конфигов
* `configs` — примеры конфигов

Сервер отвечает за регистрацию, авторизацию и хранение зашифрованных секретов.

Клиент отвечает за шифрование и расшифровку секретов. Мастер-пароль не отправляется на сервер. Сервер хранит только шифротекст, nonce и salt.

## Требования

* Go 1.25 или выше

## Конфиг сервера

В репозитории лежит пример конфига:

```text
configs/server.example.json
```

Пример содержимого:

```json
{
  "address": ":8080",
  "db_path": "gophkeeper.db",
  "jwt_secret": "change-me-for-local-development",
  "tls_cert_file": "certs/server.crt",
  "tls_key_file": "certs/server.key"
}
```

Поля:

* `address` — адрес запуска сервера
* `db_path` — путь к файлу SQLite
* `jwt_secret` — секрет для подписи JWT-токенов
* `tls_cert_file` — путь к TLS-сертификату
* `tls_key_file` — путь к TLS-ключу

`jwt_secret` обязательно должен быть задан через конфиг или переменную окружения `JWT_SECRET`.

Если `tls_cert_file` и `tls_key_file` заданы, сервер запускается через TLS. Если оба поля пустые, сервер запускается без TLS.

Для локального запуска можно скопировать пример:

```bash
cp configs/server.example.json configs/server.json
```


Путь к конфигу можно передать через флаг:

```bash
./gophkeeper-server -config configs/server.json
```

или через короткий флаг:

```bash
./gophkeeper-server -c configs/server.json
```

Также поддерживаются переменные окружения:

* `CONFIG`
* `ADDRESS`
* `DB_PATH`
* `JWT_SECRET`
* `TLS_CERT_FILE`
* `TLS_KEY_FILE`

## Конфиг клиента

В репозитории лежит пример конфига:

```text
configs/client.example.json
```

Пример содержимого:

```json
{
  "server_url": "https://localhost:8080",
  "session_file": ".gophkeeper_session.json",
  "insecure_skip_verify": true
}
```

Поля:

* `server_url` — адрес сервера
* `session_file` — путь к файлу локальной клиентской сессии
* `insecure_skip_verify` — отключение проверки TLS-сертификата для локального самоподписанного сертификата

Для локального запуска можно скопировать пример:

```bash
cp configs/client.example.json configs/client.json
```

Путь к конфигу можно передать через флаг:

```bash
./gophkeeper-client -config configs/client.json
```

или через короткий флаг:

```bash
./gophkeeper-client -c configs/client.json
```

Также поддерживаются переменные окружения:

* `CONFIG`
* `SERVER_URL`
* `SESSION_FILE`
* `INSECURE_SKIP_VERIFY`

## TLS-сертификат для локального запуска
Сгенерировать самоподписанный сертификат можно через OpenSSL:

```bash
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -sha256 -days 365 -nodes -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
```

Для локального самоподписанного сертификата в клиентском конфиге можно оставить:

```json
"insecure_skip_verify": true
```

Для реального сертификата это значение должно быть false.

## Сборка

Сервер:

```bash
go build -ldflags="-X main.buildVersion=1.0.0 -X main.buildDate=2026-06-01 -X main.buildCommit=local" -o gophkeeper-server ./cmd/server
```

Клиент:

```bash
go build -ldflags="-X main.buildVersion=1.0.0 -X main.buildDate=2026-06-01 -X main.buildCommit=local" -o gophkeeper-client ./cmd/client
```


## Запуск сервера

```bash
./gophkeeper-server -config configs/server.json
```

## CLI-клиент

Клиент поддерживает команды:

* `register` — регистрация пользователя
* `login` — вход пользователя
* `add` — добавление секрета
* `get` — получение секрета по ID
* `list` — список секретов
* `update` — обновление секрета
* `delete` — удаление секрета

Во всех примерах используется конфиг:

```bash
-config configs/client.json
```

## Регистрация

```bash
./gophkeeper-client -config configs/client.json register --login user1
```

## Логин

```bash
./gophkeeper-client -config configs/client.json login --login user1
```

## Добавление текстового секрета

```bash
./gophkeeper-client -config configs/client.json add --type text --title note1 --meta test --value "{\"text\":\"hello\"}"
```

## Добавление пары логин-пароль

```bash
./gophkeeper-client -config configs/client.json add --type login_password --title github --meta work --value "{\"login\":\"user1\",\"password\":\"secret\"}"
```

## Добавление банковской карты

```bash
./gophkeeper-client -config configs/client.json add --type card --title main-card --meta bank --value "{\"number\":\"4111111111111111\",\"holder\":\"USER TEST\",\"month\":\"12\",\"year\":\"2030\",\"cvv\":\"123\"}"
```

## Получение секрета по ID

```bash
./gophkeeper-client -config configs/client.json get --id secret-id
```

## Список секретов

```bash
./gophkeeper-client -config configs/client.json list
```

## Обновление секрета

```bash
./gophkeeper-client -config configs/client.json update --id secret-id --type text --title note1-updated --meta updated --value "{\"text\":\"new value\"}"
```


## Удаление секрета

```bash
./gophkeeper-client -config configs/client.json delete --id secret-id
```

## Проверка

```bash
go test ./...
```
