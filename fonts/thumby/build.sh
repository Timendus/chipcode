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

python3 convert.py ./fonts/5-pix-wide-ascii-layout.png ./fonts/5-pix-wide.bin 5
# python3 convert.py ./fonts/5-pix-wide-ascii-layout.png ./fonts/5-pix-wide.bin 5
# Etc
