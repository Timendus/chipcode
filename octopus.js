#!/usr/bin/env node

/**
 * This tool does pre-processing on text-based input files. It is intended for
 * use with Octo-flavoured CHIP-8 assembly code, but it can of course be used
 * with any text file.
 *
 * Features:
 *    * Include other files (`#include <filename>`)
 *    * Make code inclusion decisions (`#if <option>` / `#else` / `#end`)
 *    * Mark blocks as code or data, and get them re-ordered (`#code` / `#data`)
 *
 * Usage:
 *     ./octopus.js <input file> <ouput file> <option 1> <option 2> ...
 */

const fs      = require('fs');
const path    = require('path');
const input   = process.argv[2];
const output  = process.argv[3];
const options = process.argv.slice(4, process.argv.length);
const dirname = path.dirname(input);

const result = reorder(octopussify(loadFile(input), dirname));
fs.writeFileSync(output, result);

function loadFile(filename) {
  return fs.readFileSync(filename).toString();
}

function octopussify(file, filepath) {
  let outputting = [ true ];
  return file.split('\n')
             .filter(line => conditionals(line, outputting))
             .map(line => includes(line, filepath))
             .join('\n')
             + "#code\n";
}

function conditionals(line, outputting) {
  let matches;
  if ( matches = line.match(/^\s*#if (.*)\s*$/i) ) {
    outputting.push(options.includes(matches[1]));
    return false;
  }
  if ( line.match(/^\s*#else\s*$/i) ) {
    outputting[outputting.length - 1] = !outputting[outputting.length - 1];
    return false;
  }
  if ( line.match(/^\s*#end\s*$/i) ) {
    outputting.pop();
    return false;
  }
  return outputting.every(o => o);
}

function includes(line, filepath) {
  let matches;
  if ( matches = line.match(/^\s*#include "(.*)"\s*$/i) ) {
    const fileToInclude = filepath + path.sep + matches[1];
    return octopussify(loadFile(fileToInclude), path.dirname(fileToInclude));
  } else {
    return line;
  }
}

function reorder(file) {
  const lines = file.split('\n');
  let selected = true;
  const code = lines.filter(line => {
    if ( line.match(/^\s*#code\s*$/i) ) {
      selected = true;
      return false;
    }
    if ( line.match(/^\s*#data\s*$/i) ) {
      selected = false;
      return false;
    }
    return selected;
  });
  selected = false;
  const data = lines.filter(line => {
    if ( line.match(/^\s*#data\s*$/i) ) {
      selected = true;
      return false;
    }
    if ( line.match(/^\s*#code\s*$/i) ) {
      selected = false;
      return false;
    }
    return selected;
  });
  reportSize(code, data);
  return code.join('\n') + data.join('\n');
}

function reportSize(code, data) {
  const codeSize = calcSize(code);
  const dataSize = calcSize(data);
  const totalSize = codeSize + dataSize;
  console.log(`${code.length} lines (roughly ${codeSize} bytes) of code`);
  console.log(`${data.length} lines (roughly ${dataSize} bytes) of data`);
  if ( totalSize > 3216  ) console.warn("! Be warned: estimated program size exceeds VIP memory size");
  if ( totalSize > 3583  ) console.warn("! Be warned: estimated program size exceeds SCHIP memory size");
  if (  codeSize > 3583  ) console.warn("! Be warned: estimated code size exceeds XO-CHIP executable memory size");
  if ( totalSize > 65023 ) console.warn("! Be warned: estimated program size exceeds XO-CHIP memory size");
}

function calcSize(lines) {
  lines = lines.filter(l => l)                     // Skip empty lines
               .filter(l => !l.match(/^\s*#.*$/))  // Skip comments
               .filter(l => !l.match(/^\s*:.*$/))  // Skip labels

  const numbers = /^(?:\s+[0-9a-fx]+)+\s*$/i;
  const number = /[0-9a-fx]+/gi;
  const codeLines = lines.filter(l => !l.match(numbers)).length * 2;
  const dataLines = lines.filter(l => l.match(numbers))
                         .map(l => (l.match(number) || []).length)
                         .reduce((a, v) => a + v, 0);
  return codeLines + dataLines;
}
