# green-proxy

Реверс-прокси сервер, обеспечивающий "заглушку" поверх свободно доступного ресурса в виде кастомной авторизации. В качестве авторизационного сервера используется `password.berizaryad.ru`

## Чтобы запустить на виртуальной машине

Достаточно запустить ansible-скрипт. Он подгрузит необходимые зависимости и запустит бинарный файл `/bin/proxy-runner-lin64` как демона через `systemd`.

> Важно! Необходимо указать верные флаги при запуске бинарного файла: при автоматическом развертывании они указываются в файле `/config/green-proxy.service`. Актуальный `apiKey` следует уточнять у администратора авторизационного сервера.

## Чтобы запустить локально

Необходимо скомпилировать прокси для своей архитектуры и запустить с корректными флагами.

Для компиляции используется утилита `make`. Опции:

- `make` - запуск приложения с дефолтными флагами на локали
- `mac` - компиляция для MacOS
- `lin64` - компиляция для Linux 64-bit
- `lin32` - компиляция для Linux 32-bit
- `make build-all` - компиляция бинарных файлов для всех доступных архитектур

## Что же в переменных

Для получения автомануалов по флагам достаточно скомпилировать прокси и вызвать как `./bin/proxy-runner-mac/ --h`.
В переменных окружения можно и нужно задать:

```
Usage: proxy-runner-mac [--host HOST] [--port PORT] [--url URL] [--apikey APIKEY] [--authurl AUTHURL] [--authapi AUTHAPI]

Options:
  --host HOST            Service to be proxied (on localhost) [default: http://localhost:8080]
  --port PORT            Port to serve proxy from [default: :3000]
  --url URL              URL of the service to be proxied [default: https://superset.berizaryad.ru]
  --apikey APIKEY        API-key for /api/auth [default: 1234]
  --authurl AUTHURL      URL of the Authorization server [default: https://password.berizaryad.ru]
  --authapi AUTHAPI      URL of the /api/auth handler [default: https://password.berizaryad.ru/api/auth]
  --help, -h             display this help and exit
```

## Принцип работы

Сервис проксирует указанный ресурс (в примере - расположенный на той же виртуальной машине), регулируя доступ. Залогиненный пользователь обладает `sessionId`-cookie на стороне браузера; значение этого токена проверяется при каждом запросе к проксируемому сервису с помощью GET-запроса /api/auth к авторизационному сервису.
Если cookie не найдена или не валидна, пользователь перенаправляется на сервис авторизации.

После авторизации сервис перенаправляет пользователя обратно на запрашиваемый ресурс с query-параметром `greenToken`, где содержится `sessionId`-токен. Прокси обрабатывает запрос, проверяет валидность токена, и при успешном ответе сервера перенаправляет на URL без `greenToken`-параметра с валидным `sessionId`-cookie в браузере.
