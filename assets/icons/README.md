# Tray Icons

This directory contains the system tray icons for different application states.

## Icon Design

- **Simple rounded rectangle** - maximizes visibility in system tray
- **Edge-to-edge design** - uses 95% of canvas for maximum impact
- **No shadows or complex details** - clean, bold appearance
- **Color scheme**:
  - Idle: Gray (#808080)
  - Recording: Purple (#8A2BE2) - matching Kiro style
  - Recording indicator: Red dot in top-right corner

## Platform-Specific Icons

### Windows
- **Files**: `idle_32.ico`, `recording_32.ico`
- **Format**: ICO with embedded sizes: 16x16, 24x24, 32x32, 48x48
- Windows automatically selects the appropriate size based on DPI

### macOS
- **Normal DPI**: `idle_22.png`, `recording_22.png` (22x22)
- **Retina**: `idle_44.png`, `recording_44.png` (44x44)

### Linux
- **Files**: `idle_24.png`, `recording_24.png` (24x24)

## Available Sizes

All PNG sizes generated: 16, 22, 24, 32, 44, 48

## Regenerating Icons

To regenerate all icons:
```bash
python generate_icons.py
```

This will:
1. Generate all PNG files for all sizes
2. Create Windows ICO files with multiple embedded sizes

You can manually edit PNG files and re-run the script to regenerate ICO files from your edited PNGs.

## Icon Generation Script

The `generate_icons.py` script:
- Creates simple, bold microphone icons
- Generates all required sizes for all platforms
- Automatically creates Windows ICO files with multiple layers
- Can be re-run after manual PNG edits to update ICO files
