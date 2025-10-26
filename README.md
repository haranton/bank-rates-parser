Bank Rates Parser

Микросервисное приложение для парсинга курсов валют с веб-сайтов банков и отправки уведомлений пользователям через gRPC-сервис.

Архитектура проекта

Проект состоит из нескольких сервисов и вспомогательных компонентов:

.
├── bank-rates          # Сервис парсинга курсов валют
│   ├── internal/
│   ├── migrations/
│   └── cmd/main.go
│
├── grpc-notify         # Сервис уведомлений (gRPC)
│   ├── internal/
│   └── cmd/main.go
│
├── docker-compose.yml  # Оркестрация сервисов
└── README.md

Сервисы
Сервис	Назначение
parser	Парсит курсы валют с сайтов банков с использованием Selenium и сохраняет результаты в базу данных PostgreSQL.
notify	Получает данные о курсах через gRPC и отправляет уведомления пользователям (например, в Telegram).
chrome	Контейнер с Selenium Chrome Standalone, используемый для автоматического парсинга.
db	Хранилище данных на базе PostgreSQL 15.
Используемые технологии

Go 1.25

PostgreSQL 15

gRPC и Protocol Buffers

Docker и Docker Compose

Selenium

Zap Logger

sql-migrate

Быстрый запуск
1. Клонирование репозитория
git clone https://github.com/<username>/bank-rates-parser.git
cd bank-rates-parser

2. Настройка конфигурации

Конфигурационные файлы находятся в директориях:

bank-rates/config/local.yaml
grpc-notify/config/local.yaml


Пример файла bank-rates/config/local.yaml:

server:
  port: 8085

database:
  host: db
  port: 5432
  user: db
  password: db
  name: parser_bank_rate_db
  sslmode: disable

selenium:
  url: http://chrome:4444/wd/hub

notify:
  address: grpc-notify:50051

3. Запуск проекта

Для сборки и запуска всех сервисов используйте команду:

docker compose up -d --build


После запуска будут активны следующие контейнеры:

postgres_db_parser — PostgreSQL

chrome — Selenium

BankRateParser — сервис парсинга

grpc-notify — сервис уведомлений

Проверить состояние контейнеров:

docker compose ps


Просмотр логов конкретного сервиса:

docker compose logs -f parser

Работа с базой данных

Миграции базы данных применяются автоматически при запуске контейнера parser.
Файлы миграций располагаются в директории bank-rates/migrations/.

Применение миграций вручную:

docker exec -it BankRateParser ./main migrate up


Подключение к базе данных:

docker exec -it postgres_db_parser psql -U db -d parser_bank_rate_db

Взаимодействие сервисов

Сервисы взаимодействуют через gRPC по заранее описанному протоколу, определённому в файле:

bank-rates/proto/notification.proto


Общий поток данных:

Selenium → Parser → PostgreSQL → gRPC Notify → Telegram

Структура проекта (основные компоненты)
bank-rates/
├── cmd/
│   └── main.go             # Точка входа
├── internal/
│   ├── app/                # Инициализация приложения
│   ├── db/                 # Работа с базой данных
│   ├── scraper/            # Логика парсинга
│   ├── sender/             # Отправка данных в сервис уведомлений
│   ├── pb/                 # gRPC сгенерированные файлы
│   └── logger/             # Настройка логирования
├── migrations/             # SQL-миграции
└── config/local.yaml       # Конфигурационный файл

Полезные команды
Команда	Описание
docker compose up -d	Запуск проекта в фоновом режиме
docker compose down -v	Остановка проекта и удаление данных
docker compose logs -f parser	Просмотр логов сервиса парсера
docker exec -it postgres_db_parser psql -U db -d parser_bank_rate_db	Подключение к базе данных
