# Локальное тестирование Docker-образа

Эта инструкция описывает, как тестировать Docker-образ для сборки перед его публикацией на GitHub Container Registry.

## Предварительные требования

- Docker установлен в WSL2 (если работаешь в Windows)
- Проект доступен в WSL (например, `/mnt/c/Projects/Vox` или `~/vox`)

## Установка Docker в WSL2 (если еще не установлен)

```bash
# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Добавление пользователя в группу docker
sudo usermod -aG docker $USER

# Перезапуск WSL или выход/вход для применения изменений
```

## Шаг 1: Сборка Docker-образа

```bash
# Перейди в директорию проекта
cd /mnt/c/Projects/Vox  # или твой путь

# Собери Docker-образ
docker build -t vox-builder-test:local -f docker/builder/Dockerfile .
```

Это создаст локальный образ с тегом `vox-builder-test:local`.

## Шаг 2: Проверка образа

```bash
# Запусти контейнер в интерактивном режиме
docker run --rm -it vox-builder-test:local bash

# Внутри контейнера проверь установленные инструменты
go version          # Должна быть версия Go 1.23
make --version      # Должен быть установлен Make
git --version       # Должен быть установлен Git

# Выйди из контейнера
exit
```

## Шаг 3: Тестирование сборки проекта

### Полная сборка всех платформ

```bash
# Запусти сборку всех платформ
docker run --rm -v $(pwd):/workspace vox-builder-test:local make all

# Проверь результаты
ls -lh dist/
```

Должны появиться 6 бинарников:
- `vox-windows-amd64.exe`
- `vox-windows-arm64.exe`
- `vox-linux-amd64`
- `vox-linux-arm64`
- `vox-darwin-amd64`
- `vox-darwin-arm64`

### Сборка конкретной платформы

```bash
# Очисти предыдущие артефакты
docker run --rm -v $(pwd):/workspace vox-builder-test:local make clean

# Собери только для Linux amd64
docker run --rm -v $(pwd):/workspace vox-builder-test:local make linux-amd64

# Проверь бинарник
file dist/vox-linux-amd64
```

## Шаг 4: Тестирование запуска тестов

```bash
# Запусти тесты с виртуальным дисплеем
docker run --rm -v $(pwd):/workspace vox-builder-test:local \
  bash -c "Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 & export DISPLAY=:99.0 && sleep 1 && go test ./..."
```

Или через Makefile:

```bash
docker run --rm -v $(pwd):/workspace vox-builder-test:local make test
```

## Шаг 5: Интерактивная отладка

Если что-то пошло не так, запусти контейнер в интерактивном режиме:

```bash
# Запусти bash в контейнере с примонтированным проектом
docker run --rm -it -v $(pwd):/workspace vox-builder-test:local bash

# Внутри контейнера можешь выполнять команды вручную
cd /workspace
go mod download
make linux-amd64
go test ./...

# Выйди когда закончишь
exit
```

## Шаг 6: Публикация образа (после успешного тестирования)

Если все тесты прошли успешно, можешь опубликовать образ:

```bash
# Перетегируй образ для публикации
docker tag vox-builder-test:local ghcr.io/d-mozulyov/vox-builder:latest

# Войди в GitHub Container Registry
echo YOUR_PAT | docker login ghcr.io -u d-mozulyov --password-stdin

# Опубликуй образ
docker push ghcr.io/d-mozulyov/vox-builder:latest
```

## Полезные команды

### Очистка Docker

```bash
# Удалить тестовый образ
docker rmi vox-builder-test:local

# Удалить все неиспользуемые образы
docker image prune -a

# Посмотреть все образы
docker images
```

### Проверка размера образа

```bash
# Посмотреть размер образа
docker images vox-builder-test:local
```

Ожидаемый размер: ~400-500MB

### Быстрая проверка всего процесса

```bash
# Одна команда для полной проверки
docker build -t vox-builder-test:local -f docker/builder/Dockerfile . && \
docker run --rm -v $(pwd):/workspace vox-builder-test:local make all && \
ls -lh dist/
```

## Типичные проблемы

### Docker не найден в WSL

```bash
# Проверь, запущен ли Docker daemon
sudo service docker start

# Или установи Docker Desktop для Windows с интеграцией WSL2
```

### Ошибка "permission denied" при монтировании

```bash
# Убедись, что ты в директории проекта
pwd

# Используй полный путь
docker run --rm -v /mnt/c/Projects/Vox:/workspace vox-builder-test:local make all
```

### Сборка не находит go.mod

```bash
# Убедись, что монтируешь корень проекта, где находится go.mod
ls go.mod  # Должен существовать

# Проверь, что рабочая директория в контейнере правильная
docker run --rm -v $(pwd):/workspace vox-builder-test:local ls -la /workspace
```

## Заключение

После успешного прохождения всех тестов можешь быть уверен, что образ работает корректно и готов к публикации. Это гарантирует, что CI/CD на GitHub Actions также будет работать без проблем.
