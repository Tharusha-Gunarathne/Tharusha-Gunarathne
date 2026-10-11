#!/usr/bin/env python3
"""Turn a pencil-sketch image into data/portrait.txt.

Thin pencil lines disappear when a 2000px image is squashed into 60 columns,
so this script thickens the lines and crops to the drawing before converting.

Usage:
    python3 scripts/portrait.py sketch.png                    # sensible defaults
    python3 scripts/portrait.py sketch.png --width 70 --thicken 9
    python3 scripts/portrait.py sketch.png --crop 0.15,0.05,0.85,0.75   # x0,y0,x1,y1 as fractions
    python3 scripts/portrait.py sketch.png --invert           # for light-on-dark drawings

Needs Pillow:  nix-shell -p "python3.withPackages(p:[p.pillow])"
"""
import argparse
from PIL import Image, ImageEnhance, ImageFilter, ImageOps

RAMP = " .`':;!*#WM$@"  # light -> dark

ap = argparse.ArgumentParser()
ap.add_argument("image")
ap.add_argument("--width", type=int, default=64, help="columns (60-80 works well)")
ap.add_argument("--thicken", type=int, default=-1,
                help="line-thickening filter size in px (odd number, 0 = off, default = auto)")
ap.add_argument("--contrast", type=float, default=2.0)
ap.add_argument("--gamma", type=float, default=1.6, help=">1 darkens mid-greys (stronger lines)")
ap.add_argument("--crop", help="manual crop x0,y0,x1,y1 as fractions of the image")
ap.add_argument("--no-autocrop", action="store_true")
ap.add_argument("--aspect", type=float, default=0.54,
                help="char cell width/height; the SVG card uses 7.22/13.4 = 0.54. "
                     "Raise it (0.6, 0.65) if the portrait still looks squashed, lower it if stretched")
ap.add_argument("--invert", action="store_true")
ap.add_argument("--out", default="data/portrait.txt")
a = ap.parse_args()

img = Image.open(a.image).convert("L")
if a.invert:
    img = ImageOps.invert(img)

if a.crop:
    x0, y0, x1, y1 = (float(v) for v in a.crop.split(","))
    W, H = img.size
    img = img.crop((int(x0 * W), int(y0 * H), int(x1 * W), int(y1 * H)))
elif not a.no_autocrop:
    # crop to the drawn area (anything noticeably darker than paper) + small margin
    mask = img.point(lambda p: 255 if p < 200 else 0)
    box = mask.getbbox()
    if box:
        m = int(0.03 * max(img.size))
        img = img.crop((max(0, box[0] - m), max(0, box[1] - m),
                        min(img.width, box[2] + m), min(img.height, box[3] + m)))

# thicken dark lines so they survive downsampling: one output column covers
# img.width / width source pixels, so lines need to be a good fraction of that
cell = img.width / a.width
k = a.thicken if a.thicken >= 0 else max(3, int(cell * 0.5)) | 1
if k > 1:
    img = img.filter(ImageFilter.MinFilter(k | 1))

img = ImageOps.autocontrast(img, cutoff=1)
img = ImageEnhance.Contrast(img).enhance(a.contrast)

# SVG card cells are 7.22px wide x 13.4px tall -> squash height by that ratio
w = a.width
h = max(1, round(img.height / img.width * w * a.aspect))
img = img.resize((w, h), Image.BOX)  # BOX averages = keeps line darkness

px = img.load()
lines = []
for y in range(h):
    row = ""
    for x in range(w):
        darkness = 1 - px[x, y] / 255
        darkness = darkness ** (1 / a.gamma)
        row += RAMP[min(len(RAMP) - 1, int(darkness * len(RAMP)))]
    lines.append(row.rstrip())
while lines and not lines[0].strip():
    lines.pop(0)
while lines and not lines[-1].strip():
    lines.pop()
open(a.out, "w").write("\n".join(lines) + "\n")
print(f"wrote {a.out}: {w} cols x {len(lines)} rows (thicken={k})")
print("\n".join(lines))
