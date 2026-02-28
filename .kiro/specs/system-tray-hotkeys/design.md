# Design Document: System Tray Hotkeys

## Overview

Данный документ описывает технический дизайн системного трей-агента с поддержкой глобальных горячих клавиш и индикацией состояния для приложения Vox. Фича обеспечивает базовую инфраструктуру для фонового режима работы, управления через горячие клавиши и визуально-звуковой обратной связи.

### Ключевые цели

- Обеспечить фоновую работу приложения через системный трей
- Реализовать глобальные горячие клавиши для управления записью
- Предоставить визуальную индикацию состояния через смену иконок
- Предоставить звуковую индикацию при переходах между состояниями
- Обеспечить кроссплатформенность (Windows, Linux, macOS на x64 и arm64)

### Исследование библиотек

Для реализации требуемой функциональности были исследованы следующие Go-библиотеки:

**Системный трей:**
- **getlantern/systray** - популярная кроссплатформенная библиотека для работы с системным треем
  - Поддерживает Windows, macOS, Linux
  - Простой API для создания меню и обработки событий
  - Активно поддерживается сообществом
  - Лицензия: Apache 2.0 (совместима с MIT)

**Глобальные горячие клавиши:**
- **robotn/gohook** - библиотека для перехвата глобальных событий клавиатуры и мыши
  - Поддерживает Windows, macOS, Linux
  - Нативные биндинги для каждой платформы
  - Лицензия: Apache 2.0
- **golang-design/hotkey** - минималистичная библиотека для регистрации горячих клавиш
  - Кроссплатформенная
  - Простой API
  - Лицензия: MIT

**Воспроизведение звука:**
- **faiface/beep** - универсальная библиотека для работы со звуком
  - Поддержка различных форматов (WAV, MP3, OGG)
  - Кроссплатформенная
  - Лицензия: MIT

**Рекомендуемый выбор:**
- Системный трей: **getlantern/systray**
- Горячие клавиши: **golang-design/hotkey**
- Звук: **faiface/beep**

## Architecture

### Общая архитектура

Приложение построено на основе событийно-ориентированной архитектуры с четким разделением ответственности между компонентами.

```
┌─────────────────────────────────────────────────────────────┐
│                         Main Process                         │
│                                                               │
│  ┌────────────────┐      ┌──────────────────┐               │
│  │  Tray Manager  │◄────►│  State Machine   │               │
│  └────────────────┘      └──────────────────┘               │
│         │                         │                          │
│         │                         │                          │
│         ▼                         ▼                          │
│  ┌────────────────┐      ┌──────────────────┐               │
│  │ Hotkey Manager │      │ Indicator Manager│               │
│  └────────────────┘      └──────────────────┘               │
│                                   │                          │
│                          ┌────────┴────────┐                │
│                          ▼                 ▼                 │
│                   ┌──────────┐      ┌──────────┐            │
│                   │  Visual  │      │  Audio   │            │
│                   │Indicator │      │Indicator │            │
│                   └──────────┘      └──────────┘            │
└─────────────────────────────────────────────────────────────┘
```

### Основные компоненты

1. **State Machine** - центральный компонент управления состоянием
2. **Tray Manager** - управление иконкой и меню в системном трее
3. **Hotkey Manager** - регистрация и обработка глобальных горячих клавиш
4. **Indicator Manager** - координация визуальной и звуковой индикации
5. **Visual Indicator** - управление иконками в трее
6. **Audio Indicator** - воспроизведение звуковых сигналов

### Поток данных

1. Пользователь нажимает горячую клавишу → Hotkey Manager
2. Hotkey Manager отправляет событие → State Machine
3. State Machine изменяет состояние и уведомляет → Indicator Manager
4. Indicator Manager координирует → Visual Indicator + Audio Indicator
5. Visual Indicator обновляет иконку через → Tray Manager
6. Audio Indicator воспроизводит звуковой сигнал

## Components and Interfaces

### State Machine

Центральный компонент, управляющий состоянием приложения.

**Состояния:**
- `Idle` - приложение ожидает действий пользователя
- `Recording` - идет запись голоса
- `Processing` - обработка записанного аудио

**Интерфейс:**

```go
type State int

const (
    StateIdle State = iota
    StateRecording
    StateProcessing
)

type StateMachine interface {
    // GetState returns current application state
    GetState() State
    
    // Transition attempts to transition to a new state
    // Returns error if transition is invalid
    Transition(newState State) error
    
    // Subscribe registers a callback for state changes
    Subscribe(callback func(oldState, newState State))
}
```

**Переходы состояний:**
- `Idle → Recording` - начало записи (по горячей клавише)
- `Recording → Processing` - окончание записи (по горячей клавише)
- `Processing → Idle` - завершение обработки (автоматически)

### Tray Manager

Управляет иконкой и контекстным меню в системном трее.

**Интерфейс:**

```go
type TrayManager interface {
    // Initialize creates system tray icon and menu
    Initialize() error
    
    // SetIcon updates the tray icon
    SetIcon(iconData []byte) error
    
    // SetTooltip updates the tooltip text
    SetTooltip(text string) error
    
    // Run starts the tray event loop (blocking)
    Run()
    
    // Quit removes tray icon and exits
    Quit()
}
```

**Контекстное меню:**
- "Settings" - открыть окно настроек (будущая функциональность)
- "Exit" - завершить работу приложения

### Hotkey Manager

Регистрирует и обрабатывает глобальные горячие клавиши.

**Интерфейс:**

```go
type HotkeyManager interface {
    // Register registers a global hotkey
    // Returns error if hotkey is already taken
    Register(hotkey Hotkey, callback func()) error
    
    // Unregister removes a registered hotkey
    Unregister(hotkey Hotkey) error
    
    // UnregisterAll removes all registered hotkeys
    UnregisterAll() error
}

type Hotkey struct {
    Modifiers []Modifier // Alt, Shift, Ctrl, etc.
    Key       Key        // V, A, etc.
}

type Modifier int
type Key int
```

**Обработка ошибок:**
- Если горячая клавиша занята другим приложением, логируется предупреждение
- Приложение продолжает работу без горячих клавиш
- Пользователь может настроить другую комбинацию через UI (будущая функциональность)

### Indicator Manager

Координирует визуальную и звуковую индикацию состояния.

**Интерфейс:**

```go
type IndicatorManager interface {
    // OnStateChange handles state transitions
    OnStateChange(oldState, newState State)
    
    // SetVisualIndicator sets the visual indicator implementation
    SetVisualIndicator(indicator VisualIndicator)
    
    // SetAudioIndicator sets the audio indicator implementation
    SetAudioIndicator(indicator AudioIndicator)
}
```

### Visual Indicator

Управляет визуальной индикацией через смену иконок.

**Интерфейс:**

```go
type VisualIndicator interface {
    // UpdateIcon updates the tray icon based on state
    UpdateIcon(state State) error
}
```

**Иконки:**
- `idle.png` - серая иконка микрофона (неактивное состояние)
- `recording.png` - фиолетовая иконка микрофона (запись)
- `processing.png` - фиолетовая иконка с анимацией/индикатором (обработка)

**Требования к иконкам:**
- Формат: PNG с прозрачностью
- Размеры: 16x16, 32x32, 64x64 (для разных DPI)
- Цветовая схема: фиолетовая (в стиле Kiro)

### Audio Indicator

Воспроизводит звуковые сигналы при переходах состояний.

**Интерфейс:**

```go
type AudioIndicator interface {
    // PlaySound plays an audio feedback for state transition
    PlaySound(transition StateTransition) error
}

type StateTransition struct {
    From State
    To   State
}
```

**Звуковые файлы:**
- `start_recording.wav` - начало записи (Idle → Recording)
- `stop_recording.wav` - окончание записи (Recording → Processing)
- `processing_done.wav` - завершение обработки (Processing → Idle)

**Требования к звукам:**
- Формат: WAV (для минимальных зависимостей)
- Длительность: не более 300 мс
- Характер: приятные, ненавязчивые звуки (например, мягкие клики)

## Data Models

### Application State

```go
// AppState represents the current state of the application
type AppState struct {
    Current State
    mutex   sync.RWMutex
}

// StateHistory tracks state transitions for debugging
type StateHistory struct {
    Transitions []StateTransition
    MaxSize     int
    mutex       sync.Mutex
}

type StateTransition struct {
    From      State
    To        State
    Timestamp time.Time
}
```

### Configuration

```go
// Config holds application configuration
type Config struct {
    // Hotkey configuration
    Hotkey HotkeyConfig
    
    // Audio configuration
    Audio AudioConfig
    
    // Logging configuration
    Logging LoggingConfig
}

type HotkeyConfig struct {
    Enabled   bool
    Modifiers []Modifier
    Key       Key
}

type AudioConfig struct {
    Enabled bool
    Volume  float64 // 0.0 to 1.0
}

type LoggingConfig struct {
    Level    string // debug, info, warn, error
    FilePath string
}
```

### Platform-Specific Abstractions

```go
// Platform provides platform-specific implementations
type Platform interface {
    // GetTrayManager returns platform-specific tray manager
    GetTrayManager() TrayManager
    
    // GetHotkeyManager returns platform-specific hotkey manager
    GetHotkeyManager() HotkeyManager
    
    // GetAudioPlayer returns platform-specific audio player
    GetAudioPlayer() AudioPlayer
}
```

Это позволяет абстрагировать платформенно-зависимый код и упростить тестирование.


## Error Handling

### Стратегия обработки ошибок

Приложение использует многоуровневую стратегию обработки ошибок:

1. **Критические ошибки** - приводят к завершению приложения
2. **Некритические ошибки** - логируются, функциональность деградирует
3. **Предупреждения** - логируются, работа продолжается в полном объеме

### Критические ошибки

Ошибки, при которых приложение не может продолжать работу:

- Невозможность создать иконку в системном трее
- Невозможность инициализировать State Machine
- Критические ошибки платформенных API

**Обработка:**
```go
if err := trayManager.Initialize(); err != nil {
    log.Fatalf("Failed to initialize system tray: %v", err)
    // Cleanup resources
    cleanup()
    os.Exit(1)
}
```

### Некритические ошибки

Ошибки, при которых приложение может работать с ограниченной функциональностью:

- Невозможность зарегистрировать горячую клавишу (занята другим приложением)
- Невозможность загрузить звуковые файлы
- Ошибки обновления иконки

**Обработка:**
```go
if err := hotkeyManager.Register(hotkey, callback); err != nil {
    log.Warnf("Failed to register hotkey: %v. Application will work without hotkeys.", err)
    // Continue without hotkeys
}
```

### Предупреждения

Ситуации, требующие внимания, но не влияющие на работу:

- Невозможность воспроизвести звуковой сигнал
- Задержка обновления иконки
- Проблемы с логированием

**Обработка:**
```go
if err := audioIndicator.PlaySound(transition); err != nil {
    log.Warnf("Failed to play audio feedback: %v", err)
    // Continue without audio feedback
}
```

### Освобождение ресурсов

При любом завершении работы (нормальном или аварийном) приложение должно:

1. Отменить регистрацию всех горячих клавиш
2. Удалить иконку из системного трея
3. Закрыть все открытые файлы и соединения
4. Записать финальное сообщение в лог

**Реализация через defer:**
```go
func main() {
    // Initialize components
    tray := initializeTray()
    defer tray.Quit()
    
    hotkeys := initializeHotkeys()
    defer hotkeys.UnregisterAll()
    
    // Run application
    run()
}
```

### Логирование

Все ошибки и предупреждения логируются с соответствующим уровнем:

- `ERROR` - критические ошибки, приводящие к завершению
- `WARN` - некритические ошибки и предупреждения
- `INFO` - информационные сообщения о работе приложения
- `DEBUG` - детальная информация для отладки

**Формат лога:**
```
2024-01-15 10:30:45 [ERROR] Failed to initialize system tray: access denied
2024-01-15 10:30:50 [WARN] Hotkey Alt+Shift+V is already registered by another application
2024-01-15 10:30:51 [INFO] Application started successfully (without hotkeys)
```

## Testing Strategy

### Общий подход

Тестирование включает два взаимодополняющих подхода:

1. **Unit-тесты** - проверка конкретных примеров, граничных случаев и обработки ошибок
2. **Property-based тесты** - проверка универсальных свойств на множестве входных данных

### Unit-тестирование

**Фокус unit-тестов:**
- Конкретные примеры корректного поведения
- Граничные случаи (пустые значения, nil, максимальные размеры)
- Обработка ошибок
- Интеграция между компонентами

**Примеры unit-тестов:**

```go
// Проверка конкретного перехода состояния
func TestStateTransition_IdleToRecording(t *testing.T) {
    sm := NewStateMachine()
    err := sm.Transition(StateRecording)
    assert.NoError(t, err)
    assert.Equal(t, StateRecording, sm.GetState())
}

// Проверка обработки ошибки
func TestHotkeyManager_RegisterTakenHotkey(t *testing.T) {
    manager := NewHotkeyManager()
    hotkey := Hotkey{Modifiers: []Modifier{Alt, Shift}, Key: KeyV}
    
    // First registration should succeed
    err := manager.Register(hotkey, func() {})
    assert.NoError(t, err)
    
    // Second registration should fail
    err = manager.Register(hotkey, func() {})
    assert.Error(t, err)
}

// Проверка граничного случая
func TestVisualIndicator_UpdateWithNilIcon(t *testing.T) {
    indicator := NewVisualIndicator(nil)
    err := indicator.UpdateIcon(StateIdle)
    assert.Error(t, err)
}
```

### Property-Based тестирование

**Библиотека:** Будем использовать **gopter** - популярная библиотека для property-based testing в Go.

**Конфигурация:**
- Минимум 100 итераций на каждый property-тест
- Каждый тест помечается комментарием с ссылкой на свойство из дизайна

**Фокус property-тестов:**
- Универсальные свойства, которые должны выполняться для всех входных данных
- Инварианты системы
- Свойства round-trip (сериализация/десериализация, кодирование/декодирование)

### Моки и заглушки

Для изоляции компонентов используются моки:

```go
// Mock для TrayManager
type MockTrayManager struct {
    InitializeFunc func() error
    SetIconFunc    func([]byte) error
    QuitFunc       func()
}

// Mock для HotkeyManager
type MockHotkeyManager struct {
    RegisterFunc     func(Hotkey, func()) error
    UnregisterFunc   func(Hotkey) error
}
```

### Интеграционное тестирование

Интеграционные тесты проверяют взаимодействие компонентов:

- State Machine + Indicator Manager
- Hotkey Manager + State Machine
- Tray Manager + Visual Indicator

**Подход:**
- Использование реальных компонентов (не моков)
- Проверка полного потока данных
- Тестирование на реальных платформах (через CI/CD)

### Платформенное тестирование

Каждая платформа имеет свои особенности, поэтому необходимо:

- Автоматические тесты в GitHub Actions для всех платформ
- Ручное тестирование на реальных устройствах перед релизом
- Проверка работы на разных версиях ОС

**GitHub Actions матрица:**
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
    arch: [amd64, arm64]
```

### Покрытие кода

Целевое покрытие:
- Общее покрытие: минимум 70%
- Критические компоненты (State Machine, Hotkey Manager): минимум 85%
- Платформенно-специфичный код: минимум 60% (сложнее тестировать)


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Hotkey toggles state

*For any* current application state (idle, recording, processing), when the registered global hotkey is pressed, the application state should transition to the next valid state in the cycle.

**Validates: Requirements 2.2**

### Property 2: Hotkey cleanup on shutdown

*For any* set of registered hotkeys, when the application terminates (normally or due to error), all registered hotkeys should be unregistered and no longer trigger callbacks.

**Validates: Requirements 2.4**

### Property 3: Icon reflects state

*For any* application state (idle, recording, processing), the tray icon should visually correspond to that state - idle shows inactive icon, recording shows active icon, processing shows processing icon.

**Validates: Requirements 3.1, 3.2, 3.3**

### Property 4: Icon update responsiveness

*For any* state transition, the tray icon should be updated within 100 milliseconds of the state change.

**Validates: Requirements 3.4**

### Property 5: Audio feedback duration limit

*For any* audio feedback file used for state transitions, the duration should not exceed 300 milliseconds.

**Validates: Requirements 4.4**

### Property 6: Resource cleanup on critical error

*For any* critical error that causes application termination, all system resources (tray icon, hotkeys, file handles, audio devices) should be properly released before the process exits.

**Validates: Requirements 6.4**

### Property 7: Error logging completeness

*For any* error or warning that occurs during application execution, a corresponding log entry should be written to the log file with appropriate severity level.

**Validates: Requirements 6.5**

