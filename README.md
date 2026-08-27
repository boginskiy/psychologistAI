# psychologistAI
Psychological assistance service







your-ai-agent/
├── cmd/
│   └── agent/
│       └── main.go                 # Точка входа, инициализация DI
├── internal/
│   ├── agent/                      # Ядро бизнес-логики
│   │   ├── core/
│   │   │   ├── agent.go           # Интерфейс агента
│   │   │   ├── conversation.go    # Управление диалогом
│   │   │   └── tool.go            # Интерфейс инструментов
│   │   ├── executor/
│   │   │   ├── reAct.go          # ReAct стратегия выполнения
│   │   │   └── chain.go          # Chain-of-thought
│   │   └── memory/
│   │       ├── short_term.go     # Векторная память (встраивания)
│   │       └── long_term.go      # SQLite/Postgres для истории
│   ├── llm/                       # Адаптеры для AI моделей
│   │   ├── provider/
│   │   │   ├── interface.go      # Общий интерфейс
│   │   │   ├── openai.go
│   │   │   ├── anthropic.go
│   │   │   └── ollama.go         # Локальные модели
│   │   └── streaming/
│   │       └── handler.go        # Обработка стриминга
│   ├── tools/                     # Инструменты агента
│   │   ├── builtin/
│   │   │   ├── web_search.go
│   │   │   ├── calculator.go
│   │   │   └── file_operations.go
│   │   └── registry/
│   │       └── registry.go       # Реестр инструментов
│   ├── cli/                       # CLI слой (интерфейсы)
│   │   ├── command/
│   │   │   ├── root.go
│   │   │   ├── chat.go
│   │   │   └── config.go
│   │   ├── output/
│   │   │   ├── formatter.go      # Форматирование вывода
│   │   │   └── color.go          # TTY цвета
│   │   └── input/
│   │       ├── reader.go         # Чтение stdin
│   │       └── completer.go      # Автодополнение (readline)
│   ├── config/                    # Конфигурация
│   │   ├── loader.go             # Загрузка из YAML/ENV
│   │   └── validator.go
│   └── storage/                   # Персистентность
│       ├── repository/
│       │   ├── conversation.go
│       │   └── settings.go
│       └── migrations/            # SQL миграции
├── pkg/                           # Публичные утилиты (можно вынести)
│   ├── logger/                    # Структурированный логгинг
│   ├── retry/                     # Экспоненциальные ретраи
│   └── validator/                 # Валидация
├── configs/
│   ├── config.yaml               # Дефолтная конфигурация
│   └── config.example.yaml
├── go.mod
├── go.sum
├── Makefile                      # Автоматизация сборки
└── README.md