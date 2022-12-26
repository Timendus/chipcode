#!/usr/bin/env node

/**
 * This tool can convert image files into octo-ready sprites
 * 
 * Usage:
 *  picturething <input file> <sprite size>
 */

const imageLoader = require('./index');
const path = require('path');

const input = process.argv[2];
const size = process.argv[3] || "8x8";
const [width, height] = size.split('x').map(v => +v);

const usageString = `
Usage:
  picturething <input file> <sprite size>

Sprite size is any of the following:
  * ${imageLoader.allowedSizes.join('\n  * ')}
`;

// Check input
if ( !input ) {
  console.error(`Input file is a required parameter!\n${usageString}`);
  process.exit(128);
}
if ( !imageLoader.allowedSizes.includes(size) || isNaN(width) || isNaN(height) ) {
  console.error(`Size should be numeric\n${usageString}`);
  process.exit(128);
}

const name = path.parse(input).name;

(async function() {
    output = await imageLoader.load(input, width, height, name);
    console.log(output);
})();
