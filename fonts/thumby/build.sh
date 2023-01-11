#!/bin/bash

# Copy fonts to this directory

echo "Building font files..."
for THING in ../fonts/*
do
  FILE=`basename "$THING"`
  echo $FILE
  cp "$THING" "./fonts/${FILE}"
done

# Build font image to .bin files

python3 convert.py ./fonts/3-pix-ascii-layout.png ./fonts/3-pix.bin 3
python3 convert.py ./fonts/4-pix-low-ascii-layout.png ./fonts/4-pix-low.bin 4
python3 convert.py ./fonts/4-pix-high-ascii-layout.png ./fonts/4-pix-high.bin 4
python3 convert.py ./fonts/5-pix-narrow-ascii-layout.png ./fonts/5-pix-narrow.bin 5
python3 convert.py ./fonts/5-pix-wide-ascii-layout.png ./fonts/5-pix-wide.bin 5
python3 convert.py ./fonts/6-pix-ascii-layout.png ./fonts/6-pix.bin 6
