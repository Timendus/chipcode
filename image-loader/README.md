# Chipcode image loader

This tool can be used on the command line or as a plugin to
[Octopus](https://www.npmjs.com/package/@chipcode/octopus). It allows you to
convert image files into labels and bytes, in a format that is compatible with
the [Octo assembler](https://www.npmjs.com/package/octo-assembler).

## Installing and running

Install as an NPM package, either for a project or globally:

```bash
npm install --save-dev @chipcode/image-loader
```

You can then convert images on the command line:

```bash
npx image-loader path/to/some-file.png <optional modifiers>
```

Which will output the assembly code to stdout.

Or you can use it through Octopus like so:

```octo
:include "path/to/some-file.png" <optional modifiers>
```

## Supported file types

This image loader supports all file types that are supported by
[JIMP](https://www.npmjs.com/package/jimp), which is used for the file loading
under the hood. Currently those file types are:

  * JPEG
  * PNG
  * BMP
  * TIFF
  * GIF

## Modifiers

Modifiers are keywords that you provide to the image loader in any order, that
influence how it converts your image. Currently it knows only two types: a
`debug` flag or the target sprite resolution.

### Sprite resolution

Provide one of the following strings as a modifier, and the image loader will
cut the image into sprites of the requested dimensions:

  * `8x1`
  * `8x2`
  * `8x3`
  * `8x4`
  * `8x5`
  * `8x6`
  * `8x7`
  * `8x8`
  * `8x9`
  * `8x10`
  * `8x11`
  * `8x12`
  * `8x13`
  * `8x14`
  * `8x15`
  * `16x16`

### Debug

Provide the word `debug` as a modifier to let the image loader output the image
to the console, as well as the selected sprite resolution and all the sprites
that it has cut from the image. This quickly and easily lets you inspect if the
conversion was a success, and everything went as you expected.
