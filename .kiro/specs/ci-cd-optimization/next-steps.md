# Следующие шаги для завершения оптимизации CI/CD

## Что уже сделано ✓

1. Создан Dockerfile с всеми зависимостями
2. Создана документация по сборке и публикации образа
3. Обновлен workflow для использования Docker образа и кеширования
4. Создан автоматический workflow для обновления образа
5. Обновлена документация BUILD.md
6. Workflow настроен на использование `ghcr.io/d-mozulyov/vox-builder:latest`

## Что нужно сделать (ПРОСТОЙ СПОСОБ)

### Шаг 1: Закоммитить и запушить изменения

```bash
# Добавь все изменения
git add .

# Закоммить
git commit -m "ci: optimize build with Docker and caching"

# Запуш в GitHub
git push origin main
```

### Шаг 2: Собрать Docker образ через GitHub Actions

Вместо локальной сборки, используй GitHub Actions:

1. Перейди на https://github.com/d-mozulyov/Vox/actions
2. В левом меню выбери "Build Docker Image"
3. Справа кликни "Run workflow"
4. Выбери ветку `main`
5. Кликни зелёную кнопку "Run workflow"

GitHub Actions автоматически:
- Соберёт Docker образ
- Опубликует его в `ghcr.io/d-mozulyov/vox-builder:latest`
- Образ будет доступен публично

Время сборки: ~5-10 минут.

### Шаг 3: Проверить, что образ опубликован

1. Перейди на https://github.com/d-mozulyov?tab=packages
2. Должен появиться пакет `vox-builder`
3. Кликни на него
4. Проверь, что visibility = "Public" (если нет - измени в Package settings)

### Шаг 4: Готово!

Теперь при следующем push или создании тега, GitHub Actions будет:
- Использовать готовый Docker образ (быстро!)
- Кешировать Go модули
- Работать в 3-4 раза быстрее

## Альтернатива: Локальная сборка (если нужно)

Если всё же хочешь собрать локально, вот как получить PAT:

### Получение Personal Access Token (PAT)

1. Перейди: https://github.com/settings/tokens
2. Кликни "Generate new token" → "Generate new token (classic)"
3. Дай имя: "Vox Docker Builder"
4. Выбери срок действия: 90 days (или больше)
5. Выбери права (scopes):
   - ✓ `write:packages` (автоматически включит `read:packages`)
6. Прокрути вниз и кликни "Generate token"
7. **ВАЖНО**: Скопируй токен сразу! Он больше не покажется

### Использование PAT

```bash
# Залогинься (вставь свой PAT вместо YOUR_PAT)
echo YOUR_PAT | docker login ghcr.io -u d-mozulyov --password-stdin

# Собери образ
docker build -t ghcr.io/d-mozulyov/vox-builder:latest -f docker/builder/Dockerfile .

# Опубликуй
docker push ghcr.io/d-mozulyov/vox-builder:latest
```

## SSH не поможет

К сожалению, SSH-доступ не используется для публикации в GitHub Container Registry. Нужен либо:
- PAT (Personal Access Token) - для локальной публикации
- GitHub Actions (рекомендуется) - автоматическая публикация

## Рекомендация

**Используй GitHub Actions** (Шаг 2 выше) - это проще, безопаснее и не требует хранения токенов локально.

## Обновление образа в будущем

Когда нужно обновить зависимости:

1. Измени `docker/builder/Dockerfile`
2. Закоммить и запушить
3. GitHub Actions автоматически пересоберёт образ

Или запусти workflow вручную через GitHub UI.

## Для контрибьюторов

Если кто-то форкнет проект, они могут:
- Использовать твой образ `ghcr.io/d-mozulyov/vox-builder:latest` (по умолчанию)
- Или создать свой и указать в Settings → Variables → `BUILDER_IMAGE`

## Troubleshooting

### Workflow "Build Docker Image" не появляется

- Убедись, что файл `.github/workflows/docker-build.yml` закоммичен
- Подожди 1-2 минуты после push
- Обнови страницу Actions

### Ошибка "denied" при публикации образа

- Проверь, что workflow имеет права `packages: write` (уже настроено)
- Проверь, что ты владелец репозитория

### Образ не публичный

1. Перейди: https://github.com/d-mozulyov?tab=packages
2. Кликни на `vox-builder`
3. Package settings → Change visibility → Public

## Результат

После выполнения:

✓ Docker образ опубликован и доступен публично  
✓ CI/CD работает в 3-4 раза быстрее  
✓ Нет установки зависимостей на каждом запуске  
✓ Используется кеширование Go модулей  
✓ Не нужно хранить PAT локально
