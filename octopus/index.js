#!/usr/bin/env node

/**
 * This tool does pre-processing on text-based input files. It is intended for
 * use with Octo-flavoured CHIP-8 assembly code, but it can of course be used
 * with any text file.
 *
 * Features:
 *    * Include other files (`:include <filename>`)
 *    * Make code inclusion decisions based on options (`:if <option>` / `:unless <option>` / `:else` / `:end`)
 *    * Set and reset options (`:const SOME_OPTION 1` and `:const SOME_OPTION 0`)
 *    * Mark blocks as code or data, and get them re-ordered (`:segment code` / `:segment data`)
 *
 * Usage:
 *     npx octopus <input file> <ouput file> <option 1> <option 2> ...
 */

const fs      = require('fs');
const path    = require('path');
const input   = process.argv[2];
const output  = process.argv[3];
const options = process.argv.slice(4, process.argv.length);

// Check input
if ( !input || !output ) {
  console.error("Input and output file are required parameters!\n\nUsage:\n   octopus <input file> <ouput file> <option 1> <option 2> ...");
  process.exit();
}

// Can we `include` image files directly?
let imageLoader = false;
try {
  imageLoader = require("@chipcode/image-loader");
  if ( imageLoader.version != '0.2' )
    imageLoader = false;
} catch(e) {}

// Define some regular expressions to make our matchers easier to read
const capture = '([\\w\\-]+)';
const whitespace = '\\s+';

// Kick off the Octopussification!
try {
  const result = reorder(octopussify(loadFile(input), path.dirname(input)));
  fs.writeFileSync(output, result);
} catch(e) {  
  console.error('Could not Octopussify:', e);
}


// Main function doing the heavy lifting

function octopussify(file, filepath) {
  const outputting = [ true ];
  const mode = [ 'code' ];
  return ":segment code\n" +
    file.split('\n')
        .filter((line, i) => conditionals(line, outputting, i + 1))
        .map(line => loadImages(line, filepath))
        .map(line => includes(line, mode, filepath))
        .join('\n');
}

// Conditional code inclusion based on options, with :const, :if, :else et cetera

function conditionals(line, outputting, lineNr) {
  let matches;

  // Parse conditional inclusion of code
  if ( matches = match(line, `:if${whitespace}${capture}`) ) {
    outputting.push(options.includes(matches[1]));
    return false;
  }
  if ( matches = match(line, `:unless${whitespace}${capture}`) ) {
    outputting.push(!options.includes(matches[1]));
    return false;
  }
  if ( match(line, ':else') ) {
    outputting[outputting.length - 1] = !outputting[outputting.length - 1];
    return false;
  }
  if ( match(line, ':end') ) {
    outputting.pop();
    return false;
  }

  // Should we output / interpret the current line..?
  if ( !outputting.every(o => o) )
    return false;

  // Parse :const
  // Values 0 and "false" are "disable option"
  // Other values are "enable option"
  if ( (matches = match(line, `:const${whitespace}${capture}${whitespace}${capture}`)) ) {
    if ( matches[2] == '0' || matches[2].toLowerCase() == '0x0' ) { // TODO: better zero detection
      const index = options.indexOf(matches[1]);
      if ( index > -1 )
        options.splice(index, 1);
    } else {
      if ( !options.includes(matches[1]) )
        options.push(matches[1]);
    }
    return true;
  }

  // For debugging issues with your options
  if ( match(line, ':dump-options') ) {
    console.log(`Options on line ${lineNr}:`, options);
    return false;
  }

  return true;
}

// Including other source files with :include

function includes(line, mode, filepath) {
  if ( match(line, `:segment${whitespace}code`) ) mode[0] = 'code';
  if ( match(line, `:segment${whitespace}data`) ) mode[0] = 'data';
  
  let matches = match(line, `:include${whitespace}["'](.*\.8o)["']`);

  if ( !matches )
    return line;
    
  const fileToInclude = filepath + path.sep + matches[1];
  return octopussify(loadFile(fileToInclude), path.dirname(fileToInclude)) +
    (mode[0] == 'code' ? '\n:segment code' : '\n:segment data');
}

// Including image files directly

function loadImages(line, filepath) {
  const extensions = imageLoader ? imageLoader.allowedExtensions : [
    '.jpg',
    '.jpeg',
    '.png',
    '.bmp',
    '.tiff',
    '.gif'
  ];

  const matches = match(line, `:include${whitespace}["'](.*(${extensions.join('|')}))["'](${whitespace})?(.*)?`);

  if ( !matches )
    return line;

  const fileToInclude = filepath + path.sep + matches[1];
  const modifier = matches[4];

  // Check if we have image loading plugin installed in the first place
  if ( !imageLoader )
    throw `Attempt to include image "${matches[1]}" failed.\nInstall package '@chipcode/image-loader' to be able to include image files directly`;

  // Import the image
  return imageLoader.load(fileToInclude, modifier);
}

// Reordering :code and :data

function reorder(file) {
  const lines = file.split('\n');
  const code = lines.filter(selectLines(true));
  const data = lines.filter(selectLines(false));
  reportSize(code, data);
  return code.join('\n') + data.join('\n');
}

function selectLines(onlyCode = true) {
  let selected = onlyCode;
  return line => {
    if ( match(line, `:segment${whitespace}code`) ) {
      selected = onlyCode;
      return false;
    }
    if ( match(line, `:segment${whitespace}data`) ) {
      selected = !onlyCode;
      return false;
    }
    return selected;
  };
}

// Helpers

function loadFile(filename) {
  return fs.readFileSync(filename).toString();
}

function match(line, command) {
  return line.match(new RegExp(`^\\s*${command}\\s*(#.*)?$`, 'i'));
}

// Calculating and reporting # lines and # bytes

function reportSize(code, data) {
  const codeSize = calcSize(code);
  const dataSize = calcSize(data);
  const totalSize = codeSize + dataSize;
  console.log(`${code.length} lines (roughly ${codeSize} bytes) marked as code`);
  console.log(`${data.length} lines (roughly ${dataSize} bytes) marked as data`);
  if ( totalSize > 3216  ) console.warn("! Be warned: estimated program size exceeds VIP memory size");
  if ( totalSize > 3583  ) console.warn("! Be warned: estimated program size exceeds SCHIP memory size");
  if (  codeSize > 3583  ) console.warn("! Be warned: estimated code size exceeds XO-CHIP executable memory size");
  if ( totalSize > 65023 ) console.warn("! Be warned: estimated program size exceeds XO-CHIP memory size");
}

function calcSize(lines) {
  lines = lines.filter(l => l)                     // Skip empty lines
               .filter(l => !l.match(/^\s*#.*$/))  // Skip comments
               .filter(l => !l.match(/^\s*:.*$/))  // Skip labels and meta-instructions

  const numbers = /^(?:\s+[0-9a-fx]+)+\s*$/i;
  const number = /[0-9a-fx]+/gi;
  const codeLines = lines.filter(l => !l.match(numbers)).length * 2;
  const dataLines = lines.filter(l => l.match(numbers))
                         .map(l => (l.match(number) || []).length)
                         .reduce((a, v) => a + v, 0);
  return codeLines + dataLines;
}
