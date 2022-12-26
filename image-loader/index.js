const DEBUG = false;
const jimp = require('jimp');

module.exports = {

  allowedSizes: [
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
  ],

  load: (input, width, height, name) => {
    return jimp.read(input)
      .then(image => {
        // Check dimensions
        if ( image.bitmap.width % width != 0 ) {
          console.error("Input image size should be a multiple of requested width");
          process.exit(128);
        }

        // Load in bitmap
        const data = bitmapToMonochrome(new Uint8Array(image.bitmap.data));

        // Give some useful feedback in debug mode
        if ( DEBUG ) {
          console.log(`Selected sprite dimensions: ${width}x${height}`);
          console.log(render(data, image.bitmap.width));
          console.log("Sprites:");
        }

        // Chop it up in sprites and output to stdout
        return splitIntoSprites(data, image.bitmap.width, image.bitmap.height, width, height, name)
                  .map(sprite => `${sprite.label}\n${formatForOcto(sprite.data)}`)
                  .join('\n');
      })
      .catch(err => {
        console.error(`Could not read image ${input}:`, err);
        process.exit(1);
      });
  }

};

function assert(condition, message) {
  if (!condition) throw `Assertion failed: ${message}`;
}

// Cut up bitmap in sprites of the selected resolution
function splitIntoSprites(plane, inputWidth, inputHeight, spriteWidth, spriteHeight, name) {
  const sprites = [];
  for ( let y = 0; y < inputHeight; y += spriteHeight ) {
    for ( let x = 0; x < inputWidth / 8; x += spriteWidth / 8 ) {
      const index = y * inputWidth / 8 + x;
        const myPlane = [];
        for ( let rows = 0; rows < spriteHeight; rows++ ) {
          for ( let cols = 0; cols < spriteWidth / 8; cols++ ) {
            const value = plane[index + rows * inputWidth / 8 + cols];
            myPlane.push(value == undefined ? 0 : value);
          }
        }
      sprites.push({
        label: `: ${name}-${x}-${y / spriteHeight}`,
        data: myPlane
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
  assert(typeof width == 'number', `Expecting a valid image width`);
  assert(image.length % (width/8) == 0, 'Expecting image data to be a multiple of `width`');
  assert(image.length > 0, 'Expecting the image to have some data');
  assert(image.every(v => typeof v == 'number'), 'Expecting the image to hold only values');
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
