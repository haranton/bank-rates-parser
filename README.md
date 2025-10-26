Bank Rates Parser

Bank Rates Parser — это микросервисное приложение для парсинга курсов валют с веб-сайтов банков и последующей отправки уведомлений пользователям через gRPC-сервис.

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

Описание сервисов
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

3. Сборка и запуск проекта

Для сборки и запуска всех сервисов используйте команду:

docker compose up -d --build


После запуска будут активны следующие контейнеры:

postgres_db_parser — база данных PostgreSQL

chrome — Selenium

BankRateParser — сервис парсинга

grpc-notify — сервис уведомлений