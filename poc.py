from PIL import Image, ImageGrab
import ctypes
import math
import time

REGION_SIZE = 300
RATE_RED_MIN = 55
RATE_RED_MAX = 70
RATE_GREEN_MAX = int((100 - RATE_RED_MIN) / 2)
RATE_BLUE_MAX = RATE_GREEN_MAX
DISTANCE_MAX_ZONE = 15
NB_ITERATIONS = 700
DEBUG_OUTPUT = True

if DEBUG_OUTPUT:
    NB_ITERATIONS = 1

def compute_region_positions():
    image = ImageGrab.grab()
    width, height = image.size
    center_x, center_y = int(width / 2), int(height / 2)
    offset = REGION_SIZE / 2
    return (int(center_x - offset), int(center_y - offset), int(center_x + offset), int(center_y + offset))

def get_image(bbox):
    return ImageGrab.grab(bbox=bbox)

def load_image(path):
    return Image.open(path)

def is_near_zone(zone, y, x):
    for point in zone:
        if abs(point[0] - y) <= DISTANCE_MAX_ZONE and abs(point[1] - x) <= DISTANCE_MAX_ZONE:
            return True
    return False

def analyse_image(image, debug_path="assets/python.png"):
    image_rgb = image.convert('RGB')
    best_x, best_y = REGION_SIZE / 2, REGION_SIZE / 2
    best_distance = None
    zones = []

    for y in range(REGION_SIZE):
        for x in range(REGION_SIZE):
            r, g, b = image_rgb.getpixel((x, y))
            rate_red = 0 if r < 100 else ((r * 100) / (r + g + b))
            rate_green = 0 if rate_red < RATE_RED_MIN else ((g * 100) / (r + g + b))
            rate_blue = 0 if rate_red < RATE_RED_MIN else ((b * 100) / (r + g + b))

            if rate_red >= RATE_RED_MIN and rate_red <= RATE_RED_MAX and rate_green <= RATE_GREEN_MAX and rate_blue < RATE_BLUE_MAX:
                added = False
                """
                if len(zones):
                    zones[0].append((y, x))
                else:
                    zones.append([(y, x)])
                """
                for zone in zones:
                    if is_near_zone(zone, y, x):
                        added = True
                        zone.append((y, x))
                        break
                if not added:
                    zones.append([(y, x)])

    for zone in zones:
        if len(zone) > 1:
            sum_x = 0
            sum_y = 0

            for point in zone:
                sum_x += point[1]
                sum_y += point[0]

                if DEBUG_OUTPUT:
                    image_rgb.putpixel((point[1], point[0]), (20, 20, 255))

            x = int(sum_x / len(zone))
            y = int(sum_y / len(zone))
            dist = math.sqrt(abs(x - REGION_SIZE / 2)**2 + (y - REGION_SIZE / 2)**2)

            if best_distance is None or best_distance > dist:
                best_distance = dist
                best_x = x
                best_y = y

    if DEBUG_OUTPUT:
        for tmp in range(max(0, best_x - 3), min(best_x + 4, REGION_SIZE - 1)):
            image_rgb.putpixel((tmp, best_y), (20, 255, 20))
        for tmp in range(max(0, best_y - 3), min(best_y + 4, REGION_SIZE - 1)):
            image_rgb.putpixel((best_x, tmp), (20, 255, 20))
        image_rgb.save(debug_path)

    return int(best_x), int(best_y)

if __name__ == '__main__':
    # New stuff
    img = load_image("assets/sample-crop.png")
    analyse_image(img)

    # bbox = compute_region_positions()
    # for _ in range(NB_ITERATIONS):
    #     image = get_image(bbox)
    #     x, y = analyse_image(image)
    #     ctypes.windll.user32.SetCursorPos(bbox[0] + x, bbox[1] + y)