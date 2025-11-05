# REST-сервис для агрегации подписок
- Реализован в рамках тестового задания Junior Golang Developer (Effective Mobile). 


## Конфигурация
- Выполняется с помощью yaml-файла и переменных окружения
- Переменные окружения имеют приоритет
### Все параметры
| Раздел | Параметр | YAML | ENV | Значение по умолчанию |
|---------|-----------|------------|----------------|------------------------|
| **Logger** | Уровень логирования | `level` | `LOG_LEVEL` | `info` |
|  | Формат логов | `format` | `LOG_FORMAT` | `json` |
| **HTTP** | Хост | `host` | `HTTP_HOST` | `localhost` |
|  | Порт | `port` | `HTTP_PORT` | `8080` |
|  | Таймаут чтения | `read_timeout` | `HTTP_READ_TIMEOUT` | `10s` |
|  | Таймаут записи | `write_timeout` | `HTTP_WRITE_TIMEOUT` | `10s` |
|  | Idle timeout | `idle_timeout` | `HTTP_IDLE_TIMEOUT` | `60s` |
|  | Время на корректное завершение | `shutdown_timeout` | `HTTP_SHUTDOWN_TIMEOUT` | `20s` |
| **Postgres** | Хост | `host` | `POSTGRES_HOST` | `localhost` |
|  | Порт | `port` | `POSTGRES_PORT` | `5432` |
|  | Пользователь | — | `POSTGRES_USER` | — *(обязателен)* |
|  | Пароль | — | `POSTGRES_PASSWORD` | — *(обязателен)* |
|  | Имя базы данных | — | `POSTGRES_DB` | — *(обязателен)* |
|  | Режим SSL | `ssl_mode` | `POSTGRES_SSL_MODE` | `disable` |


### Допустимые значения параметров логов  

| Параметр            | YAML ключ | ENV переменная | Допустимые значения              |
| -------------------- | ---------- | --------------- | -------------------------------- |
| Уровень логирования | level      | LOG_LEVEL       | debug, info, warn, error         |
| Формат логов        | format     | LOG_FORMAT      | json, text                       |

  

## Запуск
- Назначте обязательные переменные окружения
```
POSTGRES_USER
POSTGRES_PASSWORD
POSTGRES_DB
CONFIG_PATH
```
- Запустите с помощью docker compose
```
docker compose up -d --build
```
- После запуска swagger-документация доступна по адресу http://localhost:8080/swagger/index.html