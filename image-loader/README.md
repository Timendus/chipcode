# CHIPCODE image loader

This tool can be used on the command line or as a plugin to
[Octopus](https://www.npmjs.com/package/@chipcode/octopus). It allows you to
convert image files into sprite data with labels and bytes, in a format that is
compatible with the [Octo
assembler](https://www.npmjs.com/package/octo-assembler).

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

## Labels and sprites

Let's say you convert or include, for example, a file called `clock.png` that is
24 pixels wide and 16 pixels high. Given the dimensions, the image loader will
guess that you want to have six 8x8 sprites. So the following labels will be
generated, which relate to the X and Y position of the sprite within the image:

  * `clock-0-0`
  * `clock-1-0`
  * `clock-2-0`
  * `clock-0-1`
  * `clock-1-1`
  * `clock-2-1`

So in this example, `clock-2-0` is the top-rightmost 8x8 sprite.

If you don't want the image loader to generate these labels, provide the
[`no-labels` modifier](#no-labels).

## Modifiers

Modifiers are keywords that you provide to the image loader in any order, that
influence how it converts your image.

### Sprite resolution

The image loader tries to make an educated guess as to what resolution sprites
you are trying to get out of the input image. It may not always guess correctly.
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

### No labels

Provide the modifier `no-labels` to suppress the sprite labels in the output.
For example if you don't need the labels because you are dynamically calculating
the offsets into the image, they just clutter up your code. And while they don't
do any harm, they do unnecessarily make your source files longer.

### Debug

Provide the word `debug` as a modifier to let the image loader output the image
to the console, as well as the selected sprite resolution and all the sprites
that it has cut from the image. This quickly and easily lets you inspect if the
conversion was a success, and everything went as you expected.
