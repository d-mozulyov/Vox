# Следующие шаги для завершения оптимизации CI/CD

## Что уже сделано ✓

1. Создан Dockerfile с всеми зависимостями
2. Создана документация по сборке и публикации образа
3. Обновлен workflow для использования Docker образа и кеширования
4. Создан автоматический workflow для обновления образа
5. Обновлена документация BUILD.md

## Что нужно сделать

### Шаг 1: Заменить placeholder в workflow

В файле `.github/workflows/build.yml` замени `REPLACE_WITH_USERNAME` на свой GitHub username:

```yaml
# Было:
image: ghcr.io/REPLACE_WITH_USERNAME/vox-builder:latest

# Должно быть (например):
image: ghcr.io/yourusername/vox-builder:latest
```

Это нужно сделать в двух местах:
- В job `test`
- В job `build`

### Шаг 2: Собрать Docker образ локально

На Windows с WSL выполни:

```bash
# Перейди в корень репозитория
cd /path/to/vox

# Собери образ (замени yourusername на свой GitHub username)
docker build -t ghcr.io/yourusername/vox-builder:latest -f docker/builder/Dockerfile .

# Проверь, что образ собрался
docker images | grep vox-builder
```

### Шаг 3: Протестировать образ локально

```bash
# Запусти контейнер
docker run --rm -it ghcr.io/yourusername/vox-builder:latest bash

# Внутри контейнера проверь инструменты:
go version
musl-gcc --version
aarch64-linux-musl-gcc --version

# Выйди из контейнера
exit
```

### Шаг 4: Опубликовать образ в GitHub Container Registry

```bash
# Создай Personal Access Token (PAT) на GitHub:
# Settings → Developer settings → Personal access tokens → Tokens (classic)
# Создай новый токен с правами: write:packages, read:packages

# Залогинься в GitHub Container Registry
echo YOUR_PAT | docker login ghcr.io -u yourusername --password-stdin

# Опубликуй образ
docker push ghcr.io/yourusername/vox-builder:latest

# Проверь, что образ опубликован
# Перейди на https://github.com/yourusername?tab=packages
```

### Шаг 5: Сделать образ публичным

1. Перейди на https://github.com/yourusername?tab=packages
2. Найди пакет `vox-builder`
3. Кликни на него
4. Перейди в "Package settings"
5. Прокрути вниз до "Danger Zone"
6. Кликни "Change visibility" → "Public"

Это важно, чтобы GitHub Actions мог скачивать образ без авторизации.

### Шаг 6: Протестировать workflow

```bash
# Создай тестовую ветку
git checkout -b test/ci-optimization

# Закоммить изменения
git add .
git commit -m "ci: optimize build with Docker and caching"

# Запуш в GitHub
git push origin test/ci-optimization

# Создай Pull Request и проверь, что тесты проходят
```

### Шаг 7: Замерить улучшение производительности

После успешного прохождения тестов:

1. Посмотри время выполнения старого workflow (в main)
2. Посмотри время выполнения нового workflow (в PR)
3. Сравни результаты

Ожидаемое улучшение: 3-4x ускорение.

### Шаг 8: Смержить изменения

Если всё работает:

```bash
# Смержи PR в main
# Через GitHub UI или:
git checkout main
git merge test/ci-optimization
git push origin main
```

## Опционально: Автоматизация обновления образа

Workflow `.github/workflows/docker-build.yml` уже создан. Он автоматически:

- Запускается при изменении Dockerfile
- Может быть запущен вручную через GitHub UI (Actions → Build Docker Image → Run workflow)

Чтобы запустить вручную:
1. Перейди в GitHub → Actions
2. Выбери "Build Docker Image"
3. Кликни "Run workflow"
4. Выбери ветку и запусти

## Troubleshooting

### Docker не найден в WSL

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
# Перезапусти WSL
```

### Ошибка при push образа

- Проверь, что PAT имеет права `write:packages`
- Проверь, что ты залогинен: `docker login ghcr.io`
- Проверь имя образа (должно быть `ghcr.io/username/vox-builder`)

### GitHub Actions не может скачать образ

- Убедись, что образ публичный (шаг 5)
- Проверь правильность имени образа в workflow
- Проверь, что образ существует: https://github.com/username?tab=packages

## Результат

После выполнения всех шагов:

- CI/CD будет работать в 3-4 раза быстрее
- Не будет установки зависимостей на каждом запуске
- Будет использоваться кеширование Go
- Сборка станет более надежной и воспроизводимой
