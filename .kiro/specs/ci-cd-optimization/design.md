# Дизайн оптимизации CI/CD

## Архитектура решения

### 1. Docker образ для сборки

**Dockerfile** будет содержать все необходимые зависимости:

```dockerfile
FROM ubuntu:22.04

# Установка базовых инструментов
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    git \
    build-essential \
    musl-tools \
    musl-dev \
    libasound2-dev \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libgl1-mesa-dev \
    libayatana-appindicator3-dev \
    xvfb \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Установка Go
RUN wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz && \
    rm go1.21.6.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Установка musl cross-compiler для arm64
RUN wget https://musl.cc/aarch64-linux-musl-cross.tgz && \
    tar -xzf aarch64-linux-musl-cross.tgz && \
    mv aarch64-linux-musl-cross /opt/ && \
    rm aarch64-linux-musl-cross.tgz

ENV PATH="/opt/aarch64-linux-musl-cross/bin:${PATH}"

WORKDIR /workspace
```

**Расположение**: `docker/builder/Dockerfile`

**Сборка и публикация**:
```bash
# Локально на Windows с WSL
docker build -t ghcr.io/<username>/vox-builder:latest -f docker/builder/Dockerfile .
docker push ghcr.io/<username>/vox-builder:latest
```

### 2. Оптимизированный GitHub Actions workflow

**Ключевые изменения**:

1. **Использование Docker образа**:
```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    container:
      image: ghcr.io/<username>/vox-builder:latest
```

2. **Кеширование Go**:
```yaml
- name: Cache Go modules
  uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

3. **Упрощенные шаги**:
- Убрать установку зависимостей
- Убрать скачивание musl cross-compiler
- Оставить только логику сборки и тестирования

### 3. Структура файлов

```
.github/
  workflows/
    build.yml          # Оптимизированный workflow
    docker-build.yml   # Workflow для обновления Docker образа (опционально)
docker/
  builder/
    Dockerfile         # Образ для сборки
    README.md          # Инструкции по сборке и публикации
```

## Процесс работы

### Первоначальная настройка

1. Создать Dockerfile
2. Собрать образ локально (WSL на Windows)
3. Опубликовать в GitHub Container Registry
4. Обновить workflow для использования образа

### Регулярная работа

1. Push кода → GitHub Actions
2. Pull Docker образа (быстро, ~10 сек)
3. Восстановление Go кеша (если есть)
4. Запуск тестов
5. Сборка для всех платформ
6. Публикация релиза

### Обновление зависимостей

Если нужно добавить новую зависимость:
1. Обновить Dockerfile
2. Пересобрать образ локально
3. Опубликовать новую версию
4. Workflow автоматически использует новый образ

## Оценка производительности

### До оптимизации
- Установка зависимостей: ~5-7 минут
- Тесты: ~1-2 минуты
- Сборка: ~3-5 минут
- **Итого**: ~10-15 минут

### После оптимизации
- Pull Docker образа: ~10-20 секунд
- Восстановление кеша: ~5-10 секунд
- Тесты: ~1-2 минуты (с кешем быстрее)
- Сборка: ~2-3 минуты (с кешем быстрее)
- **Итого**: ~3-5 минут

**Ускорение**: в 3-4 раза

## Дополнительные преимущества

1. **Воспроизводимость**: Одинаковое окружение везде
2. **Локальная разработка**: Можно использовать тот же образ локально
3. **Упрощение**: Меньше кода в workflow
4. **Надежность**: Нет риска, что apt-get update сломается
5. **Версионирование**: Можно иметь разные версии образа для разных веток

## Альтернативные подходы

### Вариант 1: GitHub Actions Cache только
- Кешировать установленные пакеты apt
- Проще, но медленнее чем Docker
- Не решает проблему полностью

### Вариант 2: Использовать готовый образ
- Искать существующий образ с Go + musl
- Может не иметь всех нужных зависимостей
- Меньше контроля

### Выбранный вариант: Собственный Docker образ
- Полный контроль над зависимостями
- Максимальная скорость
- Воспроизводимость
- **Рекомендуется**

## Риски и митигация

| Риск | Вероятность | Митигация |
|------|-------------|-----------|
| Образ слишком большой | Средняя | Использовать alpine или минимизировать слои |
| Проблемы с публикацией | Низкая | Документировать процесс, использовать GitHub CR |
| Несовместимость версий | Низкая | Фиксировать версии в Dockerfile |
| Сложность обновления | Низкая | Автоматизировать через отдельный workflow |

## План реализации

1. Создать структуру директорий
2. Написать Dockerfile
3. Создать инструкции по сборке
4. Собрать и опубликовать образ
5. Обновить workflow
6. Протестировать на тестовой ветке
7. Задокументировать процесс
