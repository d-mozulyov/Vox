# Tray Icons

This directory contains the system tray icons for different application states.

## Available Icons

- `idle_16.png`, `idle_32.png`, `idle_64.png` - Gray microphone icon (inactive state)
- `recording_16.png`, `recording_32.png`, `recording_64.png` - Purple microphone icon (recording state)
- `processing_16.png`, `processing_32.png`, `processing_64.png` - Purple microphone with indicator dot (processing state)

## Icon Specifications

- **Format**: PNG with transparency
- **Sizes**: 16x16, 32x32, 64x64 pixels (for different DPI settings)
- **Color scheme**: 
  - Idle: Gray (#808080)
  - Active states: Purple (#8A2BE2) - matching Kiro style
  - Indicator: White (#FFFFFF)
- **Naming convention**: `{state}_{size}.png`

## Current Status

Basic placeholder icons have been generated programmatically. These are simple geometric representations:
- Microphone body (rounded rectangle)
- Microphone base (vertical line with stand)
- Processing indicator (white dot in corner)

These icons are functional but can be replaced with professionally designed icons by a designer for better visual quality.
