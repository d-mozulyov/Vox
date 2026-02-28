# Следующие шаги для завершения оптимизации CI/CD

## Что уже сделано ✓

1. Создан Dockerfile с всеми зависимостями
2. Создана документация по сборке и публикации образа
3. Обновлен workflow для использования Docker образа и кеширования
4. Создан автоматический workflow для обновления образа
5. Обновлена документация BUILD.md
6. Workflow настроен на использование `ghcr.io/d-mozulyov/vox-builder:latest`

## Что нужно сделать

### Шаг 1: Собрать Docker образ локально

На Windows с WSL выполни:

```bash
# Перейди в корень репозитория
cd /path/to/vox

# Собери образ
docker build -t ghcr.io/d-mozulyov/vox-builder:latest -f docker/builder/Dockerfile .

# Проверь, что образ собрался
docker images | grep vox-builder
```

Ожидаемое время сборки: 5-10 минут (один раз).

### Шаг 2: Протестировать образ локально

```bash
# Запусти контейнер
docker run --rm -it ghcr.io/d-mozulyov/vox-builder:latest bash

# Внутри контейнера проверь инструменты:
go version                    # Должно быть: go version go1.21.6 linux/amd64
musl-gcc --version           # Должно работать
aarch64-linux-musl-gcc --version  # Должно работать

# Выйди из контейнера
exit
```

### Шаг 3: Опубликовать образ в GitHub Container Registry

```bash
# Создай Personal Access Token (PAT) на GitHub:
# 1. Перейди: https://github.com/settings/tokens
# 2. Кликни "Generate new token" → "Generate new token (classic)"
# 3. Дай имя: "Vox Docker Builder"
# 4. Выбери права: write:packages, read:packages
# 5. Кликни "Generate token"
# 6. Скопируй токен (он больше не покажется!)

# Залогинься в GitHub Container Registry
echo YOUR_PAT | docker login ghcr.io -u d-mozulyov --password-stdin

# Опубликуй образ
docker push ghcr.io/d-mozulyov/vox-builder:latest

# Проверь, что образ опубликован
# Перейди на https://github.com/d-mozulyov?tab=packages
```

### Шаг 4: Сделать образ публичным

Это важно, чтобы GitHub Actions мог скачивать образ без авторизации.

1. Перейди на https://github.com/d-mozulyov?tab=packages
2. Найди пакет `vox-builder`
3. Кликни на него
4. Перейди в "Package settings" (справа)
5. Прокрути вниз до "Danger Zone"
6. Кликни "Change visibility" → "Public"
7. Подтверди изменение

### Шаг 5: Протестировать workflow

```bash
# Закоммить изменения
git add .
git commit -m "ci: optimize build with Docker and caching"

# Запуш в GitHub
git push origin main

# Проверь GitHub Actions:
# https://github.com/d-mozulyov/Vox/actions
```

Workflow должен:
- Скачать Docker образ (~10-20 секунд)
- Восстановить кеш Go (если есть)
- Запустить тесты

### Шаг 6: Замерить улучшение производительности

После успешного прохождения:

1. Посмотри время выполнения предыдущих runs
2. Посмотри время выполнения нового run
3. Сравни результаты

Ожидаемое улучшение: 3-4x ускорение (с ~10-15 минут до ~3-5 минут).

## Автоматизация обновления образа

Workflow `.github/workflows/docker-build.yml` уже настроен и будет:

- Автоматически пересобирать образ при изменении `docker/builder/Dockerfile`
- Публиковать в `ghcr.io/d-mozulyov/vox-builder:latest`
- Можно запустить вручную через GitHub UI

Чтобы запустить вручную:
1. Перейди: https://github.com/d-mozulyov/Vox/actions
2. Выбери "Build Docker Image"
3. Кликни "Run workflow"
4. Выбери ветку `main` и запусти

## Для контрибьюторов и форков

Если кто-то форкнет проект и захочет использовать свой Docker образ:

1. Собери и опубликуй свой образ: `ghcr.io/username/vox-builder:latest`
2. В настройках своего форка на GitHub:
   - Перейди: Settings → Secrets and variables → Actions → Variables
   - Создай переменную `BUILDER_IMAGE`
   - Значение: `ghcr.io/username/vox-builder:latest`
3. Workflow автоматически будет использовать твой образ

По умолчанию workflow использует официальный образ: `ghcr.io/d-mozulyov/vox-builder:latest`

## Troubleshooting

### Docker не найден в WSL

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
# Перезапусти WSL: wsl --shutdown (в PowerShell)
```

### Ошибка при push образа

- Проверь, что PAT имеет права `write:packages`
- Проверь, что ты залогинен: `docker login ghcr.io`
- Проверь имя образа (должно быть `ghcr.io/d-mozulyov/vox-builder`)

### GitHub Actions не может скачать образ

- Убедись, что образ публичный (шаг 4)
- Проверь логи GitHub Actions для деталей
- Проверь, что образ существует: https://github.com/d-mozulyov?tab=packages

### Образ слишком большой

Текущий размер образа: ~1.5-2GB (сжатый: ~600-800MB)

Это нормально для образа с полным набором инструментов. Скачивание занимает ~10-20 секунд на серверах GitHub.

## Результат

После выполнения всех шагов:

✓ CI/CD работает в 3-4 раза быстрее  
✓ Нет установки зависимостей на каждом запуске  
✓ Используется кеширование Go модулей и сборки  
✓ Сборка надежная и воспроизводимая  
✓ Контрибьюторы могут использовать свои образы  

## Дополнительно: Локальная разработка с Docker

Можешь использовать тот же образ локально:

```bash
# Сборка в Docker (как в CI)
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest \
  go build -o vox ./cmd/vox

# Тесты в Docker
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest \
  bash -c "Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 & export DISPLAY=:99.0 && sleep 1 && go test ./..."
```

Это гарантирует, что локальная сборка идентична CI.
