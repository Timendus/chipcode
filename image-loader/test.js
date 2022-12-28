#!/usr/bin/env node

const imageLoader = require('./index');

// Nano test framework

function describe(msg, func) {
  console.log(`\n${msg}`);
  func((msg, func) => {
    let valid = true;
    func({
      assert(assertion) {
        if ( !assertion )
          valid = false;
      },
      assertEqual(object, values) {
        if ( !Object.keys(values).every(k => JSON.stringify(object[k]) == JSON.stringify(values[k])) ) {
          console.error(`\nNot all key/values in\n${JSON.stringify(values)}\nare the same in\n${JSON.stringify(object)}\n`);
          valid = false;
        }
      }
    });
    if ( valid )
      console.info(`  ✔ ${msg}`);
    else
      console.error(`  ❌ ${msg}`);
  });
}

// Tests

describe('Modifier parsing', it => {

  it("doesn't crash on empty input", ({ assert }) => {
    try {
      imageLoader.parseModifiers(null, 8, 8);
      imageLoader.parseModifiers('', 8, 8);
      imageLoader.parseModifiers(' ', 8, 8);
    } catch(e) {
      assert(false);
    }
  });

  it('defaults to sane sprite dimensions', ({ assertEqual }) => {
    assertEqual(imageLoader.parseModifiers('',  8,  8), { width: 8, height: 8 });
    assertEqual(imageLoader.parseModifiers('', 16, 16), { width: 16, height: 16 });
    assertEqual(imageLoader.parseModifiers('',  8, 16), { width: 8, height: 8 });
    assertEqual(imageLoader.parseModifiers('', 16,  8), { width: 8, height: 8 });
    assertEqual(imageLoader.parseModifiers('',  8, 12), { width: 8, height: 12 });
    assertEqual(imageLoader.parseModifiers('',  8, 24), { width: 8, height: 12 });
    assertEqual(imageLoader.parseModifiers('',  8, 36), { width: 8, height: 12 });
    assertEqual(imageLoader.parseModifiers('', 16, 24), { width: 8, height: 12 });
    assertEqual(imageLoader.parseModifiers('', 16, 36), { width: 8, height: 12 });
    assertEqual(imageLoader.parseModifiers('', 64, 6),  { width: 8, height: 6 });
    assertEqual(imageLoader.parseModifiers('', 32, 16), { width: 8, height: 8 });
  });

  it('accepts requested sprite dimensions', ({ assertEqual }) => {
    assertEqual(imageLoader.parseModifiers('8x8',   64,  6), { width: 8, height: 8 });
    assertEqual(imageLoader.parseModifiers('16x16', 32, 16), { width: 16, height: 16 });
    assertEqual(imageLoader.parseModifiers('8x4',   16,  8), { width: 8, height: 4 });
  });

});

describe('Splitting into sprites', it => {

  it('can split an image into just one sprite', ({ assert, assertEqual }) => {
    const sprites = imageLoader.splitIntoSprites([1, 2, 3, 4, 5, 6], 8, 6, 8, 6, 'test');
    assert(sprites.length == 1);
    assertEqual(sprites[0], {
      label: ': test-0-0',
      data: [1, 2, 3, 4, 5, 6]
    });
  });

  it('can split an image vertically', ({ assert, assertEqual }) => {
    const sprites = imageLoader.splitIntoSprites([1, 2, 3, 4, 5, 6], 8, 6, 8, 3, 'test');
    assert(sprites.length == 2);
    assertEqual(sprites[0], {
      label: ': test-0-0',
      data: [1, 2, 3]
    });
    assertEqual(sprites[1], {
      label: ': test-0-1',
      data: [4, 5, 6]
    });
  });

  it('can split an image horizontally', ({ assert, assertEqual }) => {
    const sprites = imageLoader.splitIntoSprites([1, 2, 3, 4, 5, 6], 16, 3, 8, 3, 'test');
    assert(sprites.length == 2);
    assertEqual(sprites[0], {
      label: ': test-0-0',
      data: [1, 3, 5]
    });
    assertEqual(sprites[1], {
      label: ': test-1-0',
      data: [2, 4, 6]
    });
  });

  it('can split an image both horizontally and vertically', ({ assert, assertEqual }) => {
    const sprites = imageLoader.splitIntoSprites([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12], 16, 6, 8, 3, 'test');
    assert(sprites.length == 4);
    assertEqual(sprites[0], {
      label: ': test-0-0',
      data: [1, 3, 5]
    });
    assertEqual(sprites[1], {
      label: ': test-1-0',
      data: [2, 4, 6]
    });
    assertEqual(sprites[2], {
      label: ': test-0-1',
      data: [7, 9, 11]
    });
    assertEqual(sprites[3], {
      label: ': test-1-1',
      data: [8, 10, 12]
    });
  });

});
