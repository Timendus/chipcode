#!/bin/bash

# Copy fonts to this directory

echo "Building font files..."
cp ../fonts/LICENSE ./fonts
mkdir -p ./fonts/images
for THING in ../fonts/*.png
do
  FILE=`basename "$THING"`
  echo $FILE
  cp "$THING" "./fonts/images/${FILE}"
done

# Build font image to .bin files

python3 convert.py ./fonts/images/3x3-fixed-auri.png ./fonts/fixed-width/3x3-auri.bin 3 3
python3 convert.py ./fonts/images/3pix-variable-threesquare.png ./fonts/variable-width/3pix-threesquare.bin 3
python3 convert.py ./fonts/images/4pix-variable-ausgezeichnet.png ./fonts/variable-width/4pix-ausgezeichnet.bin 4
python3 convert.py ./fonts/images/4pix-variable-black-sheep.png ./fonts/variable-width/4pix-black-sheep.bin 4
python3 convert.py ./fonts/images/5pix-variable-limited-narrow.png ./fonts/variable-width/5pix-limited-narrow.bin 5
python3 convert.py ./fonts/images/5pix-variable-widewest.png ./fonts/variable-width/5pix-widewest.bin 5
python3 convert.py ./fonts/images/6pix-variable-truthful.png ./fonts/variable-width/6pix-truthful.bin 6

# Minify library
# Needs a `pip3 install python-minifier`

pyminify --no-remove-annotations --remove-literal-statements fancyFont.py > fancyFont-minified.py
