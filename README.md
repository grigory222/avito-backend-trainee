[![Go](https://img.shields.io/badge/Go-1.24-blue?logo=go)](https://go.dev/)
[![Build](https://img.shields.io/github/actions/workflow/status/grigory222/avito-backend-trainee/go.yml?branch=main)](https://github.com/grigory222/avito-backend-trainee/actions)
[![License](https://img.shields.io/github/license/grigory222/avito-backend-trainee)](LICENSE)

# avito-backend-trainee
Решение тестового задания для стажировки в Авито

Посмотреть задание можно [тут](docs/Backend-trainee-assignment-spring-2025.md) \
Посмотреть спецификацию можно [тут](docs/swagger.yaml)

## Запуск приложения
1. Склонируйте приложение
    ```bash
   git clone https://github.com/grigory222/avito-backend-trainee.git
    ```
2. Установите зависимости
    ```bash
   go mod tidy
    ```
3. Запустите БД:
    ```bash
   sudo docker compose up -d
    ```
4. Создайте `.env` файл. Например, можно использовать текущий пример:
   ```bash
   mv .env.example .env
   ```
6. Запустите проект:
    ```bash
   ./run.sh
    ```
