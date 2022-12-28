#!/usr/bin/env node

/**
 * This tool can convert image files into octo-ready sprites
 * 
 * Usage:
 *  picturething <input file> <modifiers>
 */

const imageLoader = require('./index');
const input = process.argv[2];
const modifiers = process.argv.slice(3).join(' ');

const usageString = `
Usage:
  picturething <input file> <modifiers>

Where modifiers can be any of the following:
  * A sprite resolution
  * The word 'debug' to see what the tool is doing

Sprite resolution is any of the following:
  * ${imageLoader.allowedSizes.join('\n  * ')}
`;

try {
  console.log(imageLoader.load(input, modifiers));
} catch(e) {
  console.error(`Error: ${e}\n${usageString}`);
}
