# green-proxy

Реверс-прокси сервер. В качестве авторизационного сервера используется `password.berizaryad.ru`

## ЧТОБЫ ЗАПУСТИТЬ НА ВИРТУАЛЬНОЙ МАШИНЕ

Достаточно запустить ansible-скрипт. Он запустит бинарный файл `/bin/proxy-runner-lin64`; для запуска другого файла следует поменять название бинарника в скрипте.

> Важно! Никакие изменения, внесенные в файл `/env/.sh-env`, применены не будут: запускается только готовый бинарный файл

## ЧТОБЫ ЗАПУСТИТЬ ЛОКАЛЬНО

Необходимо задать переменные окружения - для этого достаточно запустить скрипт `set-env.sh`:

1. Задаем необходимые переменные в файле `/env/.sh-env`

> Важно! Для локального запуска **требуется** задать переменные в этом файле. <ез заданных переменных прокси не запустится.

2. Даем разрешение на исполнение скрипту - это достаточно сделать один (первый) раз:

```
chmod +x set-env.sh
```
3. Запускаем скрипт следующим образом:

```
eval $(./set-env.sh)
```

Далее можно запускать сам сервер - `make`

> Напоминаю, что при открытии новой вкладки терминала заданные переменные не сохранятся, и пункт 3 придется повторить

## ЧТОБЫ СКОМПИЛИРОВАТЬ

Можно самостоятельно скомпилировать бинарный файл для использования в ansible-скрипте, например. Для этого следует задать переменные окружения (см. выше) и воспользоваться Makefile. Доступные опции:

- `make` - запуск на локали
- `mac` - компиляция для MacOS
- `lin64` - компиляция для Linux 64-bit
- `lin32` - компиляция для Linux 32-bit
- `make build-all` - компиляция бинарных файлов для всех доступных архитектур

## ЧТО ЖЕ В ПЕРЕМЕННЫХ

В переменных окружения можно и нужно задать:

- `PROXY_LOCALHOST=localhost:8080` - какой сервис на локальной машине надо спрятать за прокси
- `PROXY_PORT=:3000` - на каком порту запустить сам прокси
- `PUBLIC_URL=https://superset.berizaryad.ru` - адрес ресурса, который необходимо проксировать
- `AUTH_SERVER_URL=https://password.berizaryad.ru` - адрес авторизационного сервера
- `AUTH_API_URL=https://password.berizaryad.ru/api/auth` - адрес ручки `/api/auth` для проверки валидности авторизационного токена
