from PIL import Image

def crop_icon(mode):
    img = Image.open(f'web/static/logo/logo-{mode}.png')
    pixels = img.load()
    width, height = img.size
    
    # Find the gap between the logo and the text
    gap_start = -1
    gap_end = -1
    for x in range(int(height * 0.5), width):
        # check if column is completely transparent
        is_empty = True
        for y in range(height):
            if pixels[x, y][3] > 10: # alpha > 10
                is_empty = False
                break
        if is_empty and gap_start == -1:
            gap_start = x
        elif not is_empty and gap_start != -1:
            gap_end = x
            break
            
    print(f"{mode}: gap from {gap_start} to {gap_end}")
    if gap_start != -1:
        crop_x = gap_start
        icon = img.crop((0, 0, crop_x, height))
        icon.save(f'web/static/logo/icon-{mode}.png')
        print(f"Saved icon-{mode}.png cropped at {crop_x}")

crop_icon('light')
crop_icon('dark')
