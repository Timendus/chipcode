#!/bin/bash

echo '' > build.log

# Copy library files

echo "Copying library files..."
mkdir -p ./dist
cp ./src/font-header.8o ./src/font-library.8o ./dist/

# Build font files

echo "Building font files..."
mkdir -p ./dist/fonts
for FONT in ./src/fonts/*
do
  FILE=`basename "$FONT"`
  npx octopus ${FONT} ./dist/fonts/${FILE} >> build.log 2>> build.log
done

# Build demos

echo "Building demo programs..."
mkdir -p ./demos

npx octopus ./src/test.8o ./demos/test.8o XOCHIP >> build.log 2>> build.log
npx octo --options octo-config.json ./demos/test.8o ./demos/test.html >> build.log 2>> build.log
npx octo ./demos/test.8o ./demos/test.ch8 >> build.log 2>> build.log

npx octopus ./src/demo-xochip.8o ./demos/demo-xochip.8o XOCHIP >> build.log 2>> build.log
npx octo --options octo-config.json ./demos/demo-xochip.8o ./demos/demo-xochip.html >> build.log 2>> build.log
npx octo ./demos/demo-xochip.8o ./demos/demo-xochip.ch8 >> build.log 2>> build.log

npx octopus ./src/demo-vip.8o ./demos/demo-vip.8o VIP >> build.log 2>> build.log
npx octo --options octo-config.json ./demos/demo-vip.8o ./demos/demo-vip.html >> build.log 2>> build.log
npx octo ./demos/demo-vip.8o ./demos/demo-vip.ch8 >> build.log 2>> build.log

# Calculate binary sizes

echo "Calculating file sizes..."

printf '' > font-sizes.md
for FILE in ./dist/fonts/*
do
  echo "  * $FILE"
  echo ": main" > temp.8o
  cat "$FILE" >> temp.8o
  npx octo temp.8o temp.ch8 >> build.log 2>> build.log
  FONT=`basename $FILE .8o`
  printf "| $FONT |" >> font-sizes.md
  wc -c < "temp.ch8" >> font-sizes.md
done

printf '' > library-sizes.md
for PLATFORM in SUPERCHIP VIP XOCHIP
do
  for WRAPPING in WRAP NOWRAP
  do
    echo "  * $PLATFORM $WRAPPING"
    echo ": main" > temp.8o
    echo ":const $PLATFORM 1" >> temp.8o
    echo ":const FONTLIB-${WRAPPING} 1" >> temp.8o
    cat ./dist/font-header.8o >> temp.8o
    cat ./dist/font-library.8o >> temp.8o
    npx octopus temp.8o temp.8o >> build.log 2>> build.log
    npx octo temp.8o temp.ch8 >> build.log 2>> build.log
    printf "| $PLATFORM $WRAPPING |" >> library-sizes.md
    wc -c < "temp.ch8" >> library-sizes.md
  done
done

rm temp.8o temp.ch8

# Build README file

echo "Building README file..."
sed -e '/<library-sizes-table>/ {' -e 'r library-sizes.md' -e 'd' -e '}' -e '/<font-sizes-table>/ {' -e 'r font-sizes.md' -e 'd' -e '}' README.template.md > ./README.md
rm font-sizes.md library-sizes.md
