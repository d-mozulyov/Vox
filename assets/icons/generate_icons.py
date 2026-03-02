#!/usr/bin/env python3
from PIL import Image, ImageDraw
import os, subprocess

def create_microphone_icon(size, color, is_recording=False):
    img = Image.new('RGBA', size, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    width, height = size

    # Body (large rectangle)
    body_w = width * 0.6
    body_h = height * 0.65
    body_x = (width - body_w) / 2
    draw.rectangle([body_x, 0, body_x + body_w, body_h], fill=color)

    # Stand (vertical rectangle)
    stand_w = max(3, width // 6)
    stand_h = height * 0.25
    stand_x = (width - stand_w) / 2
    draw.rectangle([stand_x, body_h, stand_x + stand_w, body_h + stand_h], fill=color)

    # Base (horizontal rectangle)
    base_w = width * 0.8
    base_h = max(3, height // 10)
    base_x = (width - base_w) / 2
    base_y = height - base_h
    draw.rectangle([base_x, base_y, base_x + base_w, height], fill=color)

    # Red indicator
    if is_recording:
        ind_size = max(4, width // 4)
        draw.ellipse([width-ind_size-2, 2, width-2, 2+ind_size], fill=(255,0,0,255))

    return img


def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    GRAY = (128, 128, 128, 255)
    PURPLE = (138, 43, 226, 255)
    sizes = [16, 22, 24, 32, 44, 48]

    print("Generating PNG icons...")
    for size in sizes:
        create_microphone_icon((size, size), GRAY, False).save(
            os.path.join(script_dir, f'idle_{size}.png'))
        create_microphone_icon((size, size), PURPLE, True).save(
            os.path.join(script_dir, f'recording_{size}.png'))
        print(f"  {size}x{size}")

    print("\nCalling convert_to_ico.py...")
    subprocess.run(['python', os.path.join(script_dir, 'convert_to_ico.py')])
    print("\nDone!")

if __name__ == '__main__':
    main()
