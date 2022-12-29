#!/usr/bin/env node

const { exec } = require('child_process');
const fs = require('fs');
const path = require('path');

require('colors');
const Diff = require('diff');

fs.readdir(__dirname, async (err, files) => {
  if ( err ) return console.error(err);
  let success = true;
  for ( const file of files ) {
    const fullPath = __dirname + path.sep + file;
    if ( fs.lstatSync(fullPath).isDirectory() ) {
      console.log(`Building "${file}"...`);
      const inputFile = `${fullPath}${path.sep}src${path.sep}index.8o`;
      const outputFile = `${fullPath}${path.sep}output.8o`;
      const expectedFile = `${fullPath}${path.sep}expected.8o`;
      await octopussify(inputFile, outputFile);
      if ( compare(outputFile, expectedFile) )
        console.log(`  ✔️ Output is as expected!`.green);
      else {
        console.log(`  ❌ Output is not as expected!`.red);
        success = false;
      }
    }
  }
  process.exit(success ? 0 : 1);
});

function octopussify(input, output) {
  return new Promise((resolve, reject) => {
    const command = __dirname + path.sep + `../index.js ${input} ${output}`;
    exec(command, (err, stdout, stderr) => {
      // We're not interested in the program's output, as long as it has run
      // without errors
      if (err) return reject(err);
      resolve();
    });
  });
}

function compare(output, expected) {
  try {
    output = fs.readFileSync(output, { encoding: 'utf8' });
  } catch(e) {
    console.error(`  🪳 Could not read 'output.8o' file: ${e}`);
    return false;
  }
  try {
    expected = fs.readFileSync(expected, { encoding: 'utf8' });
  } catch(e) {
    console.error(`  🪳 Could not read 'expected.8o' file: ${e}`);
    return false;
  }

  // We would expect all to be well here, and just return
  if ( output == expected )
    return true;

  // Otherwise, help the user to start figuring out what is wrong

  if ( !output ) {
    console.error(`  🪳 'output.8o' file is empty`);
    return false;
  }

  if ( !expected ) {
    console.error(`  🪳 'expected.8o' file is empty`);
    return false;
  }

  console.log("\n  Here's a super bad diff of the problems:\n");
  const diff = Diff.diffLines(output, expected);
  diff.forEach((part) => {
    // green for additions, red for deletions
    // grey for common parts
    if ( part.added || part.removed ) {
      const color = part.added ? 'green' :
        part.removed ? 'red' : 'grey';
      process.stdout.write(part.value[color]);
    }
  });
  console.log('\n');

  return false;
}
