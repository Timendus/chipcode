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

See README.md for valid modifiers`;

try {
  console.log(imageLoader.load(input, modifiers));
} catch(e) {
  console.error(`Error: ${e}\n${usageString}`);
}
