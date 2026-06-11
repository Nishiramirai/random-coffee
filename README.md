# Random Coffee Bot

Telegram-бот для автоматизации нетворкинга в сообществах по принципу
«Random Coffee»: участники регистрируются, а администратор командой
`/match` запускает раунд случайного подбора пар с учётом предпочтительного
формата встречи и истории прошлых раундов.

## Стек
- Go 1.22, go-telegram-bot-api v5
- PostgreSQL 16, драйвер pgx v5, миграции golang-migrate
- Docker, docker compose

## Запуск
```
git clone https://github.com/Nishiramirai/random-coffee.git
cd random-coffee
cp .env.example .env   # заполнить BOT_TOKEN и параметры БД
docker compose up -d --build
```

## Тесты
```
go test ./...
```

## Структура
- `cmd/bot` — точка входа
- `internal/config` — конфигурация из переменных окружения
- `internal/model` — модели данных
- `internal/db` — слой доступа к данным (pgx) и миграции
- `internal/bot` — диспетчер, FSM, обработчики, алгоритм матчинга, планировщик
- `migrations` — SQL-миграции схемы базы данных
