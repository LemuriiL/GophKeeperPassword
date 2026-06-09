# GophKeeper

GophKeeper — клиент-серверный менеджер секретов.

## Что умеет

* регистрация и логин пользователя
* JWT-авторизация
* хранение секретов на сервере
* синхронизация секретов между клиентами одного пользователя
* локальное шифрование секрета на клиенте мастер-паролем
* хранение на сервере только зашифрованных данных
* типы секретов:

  * login/password
  * text
  * binary
  * card
* build info у клиента и сервера
* JSON-конфиг для клиента и сервера
* graceful shutdown сервера

## Архитектура

* `cmd/server` — HTTP-сервер
* `cmd/client` — CLI-клиент
* `internal/server` — серверная логика
* `internal/client` — клиентская логика
* `internal/secure` — шифрование и derivation key
* `internal/cfg` — загрузка JSON-конфигов
* `configs` — примеры конфигов

Сервер отвечает за регистрацию, авторизацию и хранение секретов.

Клиент отвечает за шифрование и расшифровку секретов. Мастер-пароль не отправляется на сервер. Сервер хранит только шифротекст, nonce и salt.

## Требования

* Go 1.25 или выше

## Конфиг сервера

В репозитории лежит пример конфига:

`configs/server.json`

Пример содержимого:

```json
{
  "address": ":8080",
  "db_path": "gophkeeper.db",
  "jwt_secret": "supersecret"
}
```

Поля:

* `address` — адрес запуска HTTP-сервера
* `db_path` — путь к файлу локальной базы данных
* `jwt_secret` — секрет для подписи JWT-токенов

Путь к конфигу можно передать через флаг:

```bash
./gophkeeper-server -config configs/server.json
```

или через короткий флаг:

```bash
./gophkeeper-server -c configs/server.json
```

## Конфиг клиента

В репозитории лежит пример конфига:

`configs/client.json`

Пример содержимого:

```json
{
  "server_url": "http://localhost:8080",
  "session_file": ".gophkeeper_session.json"
}
```

Поля:

* `server_url` — адрес сервера
* `session_file` — путь к файлу локальной клиентской сессии

Путь к конфигу можно передать через флаг:

```bash
./gophkeeper-client -config configs/client.json
```

или через короткий флаг:

```bash
./gophkeeper-client -c configs/client.json
```

## Сборка

### Сервер

```bash
go build -ldflags="-X main.buildVersion=1.0.0 -X main.buildDate=2026-06-01 -X main.buildCommit=local" -o gophkeeper-server ./cmd/server
```

Для Windows:

```bash
go build -ldflags="-X main.buildVersion=1.0.0 -X main.buildDate=2026-06-01 -X main.buildCommit=local" -o gophkeeper-server.exe ./cmd/server
```

### Клиент

```bash
go build -ldflags="-X main.buildVersion=1.0.0 -X main.buildDate=2026-06-01 -X main.buildCommit=local" -o gophkeeper-client ./cmd/client
```

Для Windows:

```bash
go build -ldflags="-X main.buildVersion=1.0.0 -X main.buildDate=2026-06-01 -X main.buildCommit=local" -o gophkeeper-client.exe ./cmd/client
```

## Запуск сервера

```bash
./gophkeeper-server -config configs/server.json
```

Для Windows:

```bash
.\gophkeeper-server.exe -config configs/server.json
```

По умолчанию сервер запускается на адресе из конфига:

```text
:8080
```

## CLI-клиент

Клиент поддерживает команды:

* `register` — регистрация пользователя
* `login` — вход пользователя
* `add` — добавление секрета
* `get` — получение секрета по ID
* `list` — вывод списка секретов
* `update` — обновление секрета
* `delete` — удаление секрета

Во всех примерах ниже используется конфиг:

```bash
-config configs/client.json
```

Для Windows вместо `./gophkeeper-client` можно использовать:

```bash
.\gophkeeper-client.exe
```

## Регистрация

```bash
./gophkeeper-client -config configs/client.json register --login user1 --password pass1
```

После успешной регистрации клиент сохраняет локальную сессию в файл, указанный в `session_file`.

## Логин

```bash
./gophkeeper-client -config configs/client.json login --login user1 --password pass1
```

После успешного логина клиент обновляет локальную сессию.

## Добавление секрета

Пример добавления текстового секрета:

```bash
./gophkeeper-client -config configs/client.json add --type text --title note1 --meta test --value "{\"text\":\"hello\"}" --master-password masterpass
```

После успешного добавления команда выводит ID созданного секрета.

Пример добавления пары логин/пароль:

```bash
./gophkeeper-client -config configs/client.json add --type login_password --title github --meta work --value "{\"login\":\"user1\",\"password\":\"secret\"}" --master-password masterpass
```

Пример добавления банковской карты:

```bash
./gophkeeper-client -config configs/client.json add --type card --title main-card --meta bank --value "{\"number\":\"4111111111111111\",\"holder\":\"USER TEST\",\"month\":\"12\",\"year\":\"2030\",\"cvv\":\"123\"}" --master-password masterpass
```

## Получение секрета по ID

```bash
./gophkeeper-client -config configs/client.json get --id secret-id --master-password masterpass
```

Команда выводит данные секрета и расшифрованное значение.

## Список секретов

```bash
./gophkeeper-client -config configs/client.json list --master-password masterpass
```

Команда выводит список секретов текущего пользователя.

## Обновление секрета

```bash
./gophkeeper-client -config configs/client.json update --id secret-id --type text --title note1-updated --meta updated --value "{\"text\":\"new value\"}" --master-password masterpass
```

## Удаление секрета

```bash
./gophkeeper-client -config configs/client.json delete --id secret-id
```

## Безопасность

Пароли через CLI-флаги `--password` и `--master-password` небезопасны, потому что могут остаться в истории команд и списке процессов.

В обычном сценарии клиент должен запрашивать пароль и мастер-пароль интерактивно.

Флаги можно использовать для тестов и неинтерактивных сценариев.

## Пример полного сценария

Сначала запускается сервер:

```bash
./gophkeeper-server -config configs/server.json
```

Затем пользователь регистрируется:

```bash
./gophkeeper-client -config configs/client.json register --login user1 --password pass1
```

Добавляет секрет:

```bash
./gophkeeper-client -config configs/client.json add --type text --title note1 --meta test --value "{\"text\":\"hello\"}" --master-password masterpass
```

Смотрит список секретов:

```bash
./gophkeeper-client -config configs/client.json list --master-password masterpass
```

Получает конкретный секрет по ID:

```bash
./gophkeeper-client -config configs/client.json get --id secret-id --master-password masterpass
```

## Build info

При запуске клиент и сервер выводят build info:

```text
Build version: 1.0.0
Build date: 2026-06-01
Build commit: local
```

Если значения не были переданы через `ldflags`, будут выведены значения по умолчанию:

```text
Build version: N/A
Build date: N/A
Build commit: N/A
```
