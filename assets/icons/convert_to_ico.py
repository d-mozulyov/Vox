#!/usr/bin/env python3
"""
Convert PNG icons to Windows ICO format with multiple sizes embedded.
This script can be run independently after manual PNG editing.
"""

from PIL import Image
import os

def create_ico_files():
    """Create Windows ICO files with multiple layers from existing PNG files."""
    script_dir = os.path.dirname(os.path.abspath(__file__))

    # Windows ICO should contain: 16, 24, 32, 48
    sizes = [16, 24, 32, 48]

    print("Creating Windows ICO files from PNG sources...")

    for state in ['idle', 'recording']:
        # Load all available PNG sizes
        images = []
        missing_sizes = []

        for size in sizes:
            png_path = os.path.join(script_dir, f'{state}_{size}.png')
            if not os.path.exists(png_path):
                missing_sizes.append(size)
                continue

            img = Image.open(png_path)
            if img.mode != 'RGBA':
                img = img.convert('RGBA')
            images.append(img)

        if missing_sizes:
            print(f"  Warning: Missing PNG files for {state}: {missing_sizes}")

        if not images:
            print(f"  Error: No PNG files found for {state}, skipping")
            continue

        # Save ICO with all loaded images
        ico_path = os.path.join(script_dir, f'{state}.ico')
        # The first image is used as base, but all images are embedded
        images[0].save(ico_path, format='ICO', append_images=images[1:])
        print(f"  Created {ico_path} with {len(images)} sizes from individual PNG files")

    print("\nICO file creation complete!")

if __name__ == '__main__':
    create_ico_files()
