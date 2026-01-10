from PIL import Image, ImageDraw
import os

icons = ['home', 'home-active', 'message', 'message-active', 'order', 'order-active', 'profile', 'profile-active']
colors = {'active': (250, 44, 25), 'normal': (0, 0, 0)}

if not os.path.exists('src/assets/icons'):
    os.makedirs('src/assets/icons')

for icon in icons:
    # Check if active or normal
    color = colors['active'] if 'active' in icon else colors['normal']
    img = Image.new('RGBA', (64, 64), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    
    # Draw logic
    if 'home' in icon:
        d.polygon([(32, 10), (10, 30), (54, 30)], outline=color, width=3)
        d.rectangle([18, 30, 46, 54], outline=color, width=3)
    elif 'message' in icon:
        d.ellipse([10, 15, 54, 49], outline=color, width=3)
        d.polygon([(15, 40), (10, 50), (25, 45)], fill=color)
    elif 'order' in icon:
        d.rectangle([15, 10, 49, 54], outline=color, width=3)
        d.line([(20, 20), (44, 20)], fill=color, width=2)
        d.line([(20, 30), (44, 30)], fill=color, width=2)
    elif 'profile' in icon:
        d.ellipse([22, 10, 42, 30], outline=color, width=3)
        d.arc([10, 35, 54, 64], 180, 0, fill=color, width=3)

    if 'active' in icon:
        # Fill for active state (simplified)
        d.rectangle([25, 25, 39, 39], fill=color)

    img.save(f'src/assets/icons/{icon}.png')
print("Icons generated.")
