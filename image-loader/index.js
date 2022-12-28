const fs = require('fs');
const path = require('path');
const jimp = require('jimp');
const mime = require('mime-types');

const allowedExtensions = [
  '.jpg',
  '.jpeg',
  '.png',
  '.bmp',
  '.tiff',
  '.gif'
];

const allowedSizes = [
  '8x1',
  '8x2',
  '8x3',
  '8x4',
  '8x5',
  '8x6',
  '8x7',
  '8x8',
  '8x9',
  '8x10',
  '8x11',
  '8x12',
  '8x13',
  '8x14',
  '8x15',
  '16x16'
];

function load(input, modifiers, debug=false) {
  assert(typeof input == 'string' && input.length > 3)
    .orThrow(`Image file name is a required parameter`);
  assert(!modifiers || typeof modifiers == 'string')
    .orThrow(`Expecting modifiers to be falsy or a string`);

  const filename = path.parse(input);

  assert(allowedExtensions.includes(filename.ext))
    .orThrow(`Requested image file "${input}" does not have one of the supported file name extensions`);

  assert(fs.existsSync(input))
    .orThrow(`Requested image file "${input}" doesn't exist`);

  const buffer = fs.readFileSync(input);
  const type = mime.lookup(input);
  const decoder = jimp.decoders[type];

  assert(decoder)
    .orThrow(`Requested image file "${input}" cannot be decoded`);

  const file = decoder(buffer);
  const data = bitmapToMonochrome(new Uint8Array(file.data));

  assert(file.width > 0 && file.height > 0)
    .orThrow(`Requested image file "${input}" has invalid dimensions`);

  modifiers = parseModifiers(modifiers, file.width, file.height);
  modifiers.debug ||= debug

  assert(file.width % modifiers.width == 0)
    .orThrow(`Input image size should be a multiple of requested width`);

  // Give some useful feedback in debug mode
  if ( modifiers.debug ) {
    console.info(render(data, file.width));
    console.info(`Selected sprite dimensions: ${modifiers.width}x${modifiers.height}`);
  }

  const sprites = splitIntoSprites(
    data,
    file.width,
    file.height,
    modifiers.width,
    modifiers.height,
    filename.name
  );

  if ( modifiers.debug ) {
    console.info("Resulting sprites:\n");
    sprites.forEach(sprite => console.info(render(sprite.data, modifiers.width)));
  }

  return sprites.map(sprite => `${sprite.label}\n${formatForOcto(sprite.data)}`)
                .join('');;
}

function assert(assertion) {
  return {
    orThrow(msg) {
      if ( !assertion ) throw msg;
    }
  };
}

function parseModifiers(modifiers, imageWidth, imageHeight) {
  // Determine sprite dimensions
  let width, height, spriteSize;
  // Is there a requested sprite dimension in the modifiers?
  if ( modifiers )
    spriteSize = allowedSizes.find(size =>
      modifiers.match(new RegExp(`(^|\\s)${size}(\\s|$)`, 'i'))
    )
  if ( spriteSize ) {
    [width, height] = spriteSize.split('x').map(v => +v);
  } else {
    // Otherwise, default to some sane values
    if ( imageWidth == 16 && imageHeight == 16 ) {
      width = 16;
      height = 16;
    } else {
      vertSprites = 1;
      while ( !Number.isInteger(imageHeight / vertSprites) || imageHeight / vertSprites >= 16 )
        vertSprites++;
      width = 8;
      height = imageHeight / vertSprites;
    }
  }

  return {
    width,
    height,
    debug: modifiers && modifiers.match(/debug/i)
  };
}

// Cut up bitmap in sprites of the selected resolution
function splitIntoSprites(image, inputWidth, inputHeight, spriteWidth, spriteHeight, name) {
  const sprites = [];
  for ( let y = 0; y < inputHeight; y += spriteHeight ) {
    for ( let x = 0; x < inputWidth / 8; x += spriteWidth / 8 ) {
      const index = y * inputWidth / 8 + x;
        const sprite = [];
        for ( let rows = 0; rows < spriteHeight; rows++ ) {
          for ( let cols = 0; cols < spriteWidth / 8; cols++ ) {
            const value = image[index + rows * inputWidth / 8 + cols];
            sprite.push(value == undefined ? 0 : value);
          }
        }
      sprites.push({
        label: `: ${name}-${x}-${y / spriteHeight}`,
        data: sprite
      });
    }
  }
  return sprites;
}

// Outputs the 1-bit bitmap in a format Octo understands
function formatForOcto(image, stride=16) {
  let output = "";
  let offset = 0;
  for ( let i = 0; i < image.length; i += stride ) {
    const line = image.slice(i, i + stride)
                      .map(v => '0x' + v.toString(16).padStart(2, '0'))
                      .join(' ');
    output +=  `  ${line}\n`;
  }
  return output;
}

// Take RGBA bitmap data, reduce to 1-bit bitmap
function bitmapToMonochrome(data) {
  const bitmap = [];
  let byte = 0;
  let bitmask = 1 << 7;
  for ( let i = 0; i < data.length; i += 4 ) {
    const color = (data[i] + data[i+1] + data[i+2]) / 3;
    if ( color > 128 ) byte = byte | bitmask;
    bitmask >>= 1;
    if ( bitmask == 0 ) {
      bitmap.push(byte);
      byte = 0;
      bitmask = 1 << 7;
    }
  }
  return bitmap;
}

// Visualise 1-bit bitmap
function render(image, width) {
  assert(typeof width == 'number')
    .orThrow(`Expecting a valid sprite width`);
  assert(image.length % (width/8) == 0)
    .orThrow('Expecting sprite data to be a multiple of `width`');
  assert(image.length > 0)
    .orThrow('Expecting the sprite to have some data');
  assert(image.every(v => typeof v == 'number'))
    .orThrow('Expecting the sprite to hold only values');

  let offset = 0;
  let output = '';
  const height = image.length / (width / 8);
  for ( let y = 0; y < height; y++ ) {
    output += image.slice(offset, offset + width/8)
                    .map(v =>
                      v.toString(2)
                      .padStart(8, '0')
                      .replaceAll('0', '  ')
                      .replaceAll('1', '██')
                    )
                    .join('');
    output += '\n';
    offset += width/8;
  }
  return output;
}

module.exports = {
  version: '0.2',
  allowedExtensions,
  allowedSizes,
  load,

  // For testing
  parseModifiers,
  splitIntoSprites
};
