# Media Service

Проект на Go с использованием Docker и PostgreSQL.

## Настройка перед первым запуском

Перед первым запуском проекта необходимо создать файл с переменными окружения.

1.  Создайте в корне проекта файл с именем `.env`.
2.  Скопируйте в него следующее содержимое:

    ```env
    # PostgreSQL settings
    POSTGRES_USER=user
    POSTGRES_PASSWORD=password
    POSTGRES_DB=media_db

    # Application connection settings
    # The host is the service name in docker-compose.yml
    DATABASE_HOST=db
    DATABASE_PORT=5432
    ```

Этот файл игнорируется системой контроля версий (`.gitignore`), поэтому ваши секреты останутся в безопасности.

## Запуск проекта

Для запуска проекта убедитесь, что у вас установлен Docker и Docker Compose.

Выполните следующую команду в корневой директории проекта:

```bash
docker-compose up -d --build
```

Эта команда соберет образ вашего Go-приложения, запустит контейнеры для приложения и базы данных в фоновом режиме.

-   Веб-сервер будет доступен по адресу: [http://localhost:8080](http://localhost:8080)
-   База данных PostgreSQL будет доступна по адресу `127.0.0.1:5432`.

## Остановка проекта

Для остановки и удаления контейнеров выполните команду:

```bash
docker-compose down
```