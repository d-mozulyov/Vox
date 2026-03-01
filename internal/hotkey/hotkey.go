package hotkey

import (
	"fmt"
	"sync"

	"github.com/d-mozulyov/vox/internal/platform"
	"golang.design/x/hotkey"
)

// Modifier represents a keyboard modifier key (Alt, Shift, Ctrl, etc.)
type Modifier int

const (
	// ModAlt represents the Alt modifier key
	ModAlt Modifier = iota
	// ModShift represents the Shift modifier key
	ModShift
	// ModCtrl represents the Ctrl modifier key
	ModCtrl
	// ModWin represents the Windows/Command modifier key
	ModWin
)

// Key represents a keyboard key
type Key int

const (
	// KeyA represents the A key
	KeyA Key = iota
	// KeyB represents the B key
	KeyB
	// KeyC represents the C key
	KeyC
	// KeyD represents the D key
	KeyD
	// KeyE represents the E key
	KeyE
	// KeyF represents the F key
	KeyF
	// KeyG represents the G key
	KeyG
	// KeyH represents the H key
	KeyH
	// KeyI represents the I key
	KeyI
	// KeyJ represents the J key
	KeyJ
	// KeyK represents the K key
	KeyK
	// KeyL represents the L key
	KeyL
	// KeyM represents the M key
	KeyM
	// KeyN represents the N key
	KeyN
	// KeyO represents the O key
	KeyO
	// KeyP represents the P key
	KeyP
	// KeyQ represents the Q key
	KeyQ
	// KeyR represents the R key
	KeyR
	// KeyS represents the S key
	KeyS
	// KeyT represents the T key
	KeyT
	// KeyU represents the U key
	KeyU
	// KeyV represents the V key
	KeyV
	// KeyW represents the W key
	KeyW
	// KeyX represents the X key
	KeyX
	// KeyY represents the Y key
	KeyY
	// KeyZ represents the Z key
	KeyZ
)

// Hotkey represents a global hotkey combination
type Hotkey struct {
	Modifiers []Modifier
	Key       Key
}

// String returns a string representation of the hotkey
func (h Hotkey) String() string {
	result := ""
	for i, mod := range h.Modifiers {
		if i > 0 {
			result += "+"
		}
		switch mod {
		case ModAlt:
			result += "Alt"
		case ModShift:
			result += "Shift"
		case ModCtrl:
			result += "Ctrl"
		case ModWin:
			result += "Win"
		}
	}
	if len(h.Modifiers) > 0 {
		result += "+"
	}
	result += string(rune('A' + int(h.Key)))
	return result
}

// HotkeyManager defines the interface for managing global hotkeys
type HotkeyManager interface {
	// Register registers a global hotkey with a callback
	// Returns error if the hotkey is already taken by another application
	Register(hk Hotkey, callback func()) error

	// Unregister removes a registered hotkey
	Unregister(hk Hotkey) error

	// UnregisterAll removes all registered hotkeys
	UnregisterAll() error
}

// hotkeyManager implements the HotkeyManager interface
type hotkeyManager struct {
	hotkeys map[string]*hotkeyEntry
	mutex   sync.RWMutex
}

// hotkeyEntry stores a registered hotkey and its associated resources
type hotkeyEntry struct {
	hk       *hotkey.Hotkey
	callback func()
	stopChan chan struct{}
}

// NewHotkeyManager creates a new hotkey manager
func NewHotkeyManager() HotkeyManager {
	return &hotkeyManager{
		hotkeys: make(map[string]*hotkeyEntry),
	}
}

// Register registers a global hotkey with a callback
func (hm *hotkeyManager) Register(hk Hotkey, callback func()) error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	logger := platform.GetLogger()
	key := hk.String()

	// Check if hotkey is already registered
	if _, exists := hm.hotkeys[key]; exists {
		err := fmt.Errorf("hotkey %s is already registered", key)
		logger.Warn("Attempted to register already registered hotkey: %s", key)
		return err
	}

	// Convert our types to golang-design/hotkey types
	mods := convertModifiers(hk.Modifiers)
	k := convertKey(hk.Key)

	// Create the hotkey
	nativeHotkey := hotkey.New(mods, k)

	// Try to register the hotkey
	err := nativeHotkey.Register()
	if err != nil {
		// If hotkey is already registered, try to unregister and register again
		logger.Warn("Hotkey %s appears to be already registered, attempting to unregister first", key)

		// Try to unregister (this might fail if it's registered by another process)
		_ = nativeHotkey.Unregister()

		// Try to register again
		err = nativeHotkey.Register()
		if err != nil {
			logger.Error("Failed to register hotkey %s: %v", key, err)
			return fmt.Errorf("failed to register hotkey %s: %w", key, err)
		}
	}

	// Create stop channel for the listener goroutine
	stopChan := make(chan struct{})

	// Store the hotkey entry
	hm.hotkeys[key] = &hotkeyEntry{
		hk:       nativeHotkey,
		callback: callback,
		stopChan: stopChan,
	}

	logger.Info("Hotkey registered successfully: %s", key)

	// Start listening for hotkey events in a goroutine
	go func() {
		for {
			select {
			case <-nativeHotkey.Keydown():
				callback()
			case <-stopChan:
				return
			}
		}
	}()

	return nil
}

// Unregister removes a registered hotkey
func (hm *hotkeyManager) Unregister(hk Hotkey) error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	logger := platform.GetLogger()
	key := hk.String()

	entry, exists := hm.hotkeys[key]
	if !exists {
		err := fmt.Errorf("hotkey %s is not registered", key)
		logger.Warn("Attempted to unregister non-existent hotkey: %s", key)
		return err
	}

	// Stop the listener goroutine
	close(entry.stopChan)

	// Unregister the hotkey
	if err := entry.hk.Unregister(); err != nil {
		logger.Error("Failed to unregister hotkey %s: %v", key, err)
		return fmt.Errorf("failed to unregister hotkey %s: %w", key, err)
	}

	// Remove from map
	delete(hm.hotkeys, key)

	logger.Info("Hotkey unregistered successfully: %s", key)

	return nil
}

// UnregisterAll removes all registered hotkeys
func (hm *hotkeyManager) UnregisterAll() error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	logger := platform.GetLogger()

	if len(hm.hotkeys) == 0 {
		logger.Info("No hotkeys to unregister")
		return nil
	}

	logger.Info("Unregistering all hotkeys (%d total)", len(hm.hotkeys))

	var errors []error

	for key, entry := range hm.hotkeys {
		// Stop the listener goroutine
		close(entry.stopChan)

		// Unregister the hotkey
		if err := entry.hk.Unregister(); err != nil {
			logger.Error("Failed to unregister hotkey %s: %v", key, err)
			errors = append(errors, fmt.Errorf("failed to unregister hotkey %s: %w", key, err))
		} else {
			logger.Info("Hotkey unregistered: %s", key)
		}
	}

	// Clear the map
	hm.hotkeys = make(map[string]*hotkeyEntry)

	if len(errors) > 0 {
		return fmt.Errorf("errors during unregister all: %v", errors)
	}

	logger.Info("All hotkeys unregistered successfully")

	return nil
}

// convertModifiers converts our Modifier type to golang-design/hotkey modifiers
// Uses numeric values directly for cross-platform compatibility
func convertModifiers(mods []Modifier) []hotkey.Modifier {
	result := make([]hotkey.Modifier, 0, len(mods))
	for _, mod := range mods {
		switch mod {
		case ModAlt:
			// Alt: 0x1 on Windows/macOS, Mod1 (1<<3=8) on Linux
			// We use the common value that works across platforms
			result = append(result, hotkey.Modifier(1<<3)) // Mod1 on Linux, works as ModAlt elsewhere
		case ModShift:
			result = append(result, hotkey.Modifier(1<<0)) // ModShift: 1<<0 = 1
		case ModCtrl:
			result = append(result, hotkey.Modifier(1<<2)) // ModCtrl: 1<<2 = 4
		case ModWin:
			// Win/Super: 0x8 on Windows, Mod4 (1<<6=64) on Linux, ModCmd on macOS
			result = append(result, hotkey.Modifier(1<<6)) // Mod4 on Linux, works as ModWin elsewhere
		}
	}
	return result
}

// convertKey converts our Key type to golang-design/hotkey key
func convertKey(k Key) hotkey.Key {
	// golang-design/hotkey uses the same key codes as our Key type
	// We just need to add the offset to get to the correct key code
	return hotkey.Key(hotkey.KeyA + hotkey.Key(k))
}
