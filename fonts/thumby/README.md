# FancyFont

_A MicroPython font library for the Thumby playable keychain console_

FancyFont is a drop-in replacement to the `setFont()` and `drawText()` functions in
the Thumby API. It has a couple of improvements over those functions:

  * Support for word wrapping with `drawTextWrapped()`
  * Support for clipping (and not just at the edge of the screen)
  * Support for the newline character (`\n`) in strings
  * Support for variable width fonts
  * Can render to any bytearray buffer, so "should work" with grayscale too

This sub-project is in this repository because I didn't think it really needed
its own repository. And also so it can share font definitions with the CHIP-8
version of the library.

## Visual examples

TODO: Nice pictures here

## Explanation by example

To render this to the display using FancyFont:

![The visual result of the code example below](./pictures/code-example.png)

You need to do something like this:

```python
# Fix import path so we can import files in our game directory
from sys import path
path.insert(0, '/Games/MyGame')

# Load either `thumby` or `thumbyGraphics` library
import thumbyGraphics

# Load and initialize `FancyFont` library
from fancyFont import FancyFont
fancyFont = FancyFont(thumbyGraphics.display.display.buffer)

# Clear the screen
thumbyGraphics.display.fill(0)

# Select a TinyCircuits font and draw some text
fancyFont.setFont('font5x7.bin', 5, 7)
fancyFont.drawText('FancyFont\n---------', 9, 4)

# Select a FancyFont font and draw some wrapping text
fancyFont.setFont('5-pix-wide.bin')
fancyFont.drawTextWrapped('Hello there, Thumby user!', 8, 22)

# Show the result
thumbyGraphics.display.update()
```

There are a few interesting differences with the regular Thumby API to note in
the example above:

### `setFont`

`setFont` will actually search the system path if the given string is not an
absolute path. This means you can just call `setFont('font5x7.bin', 5, 7)` and
it will find this file in the `/lib` folder, which is in the default path. After
adding our project's directory to the path (in lines 2 and 3 of the example) we
can also load a font that is stored in our project directory, with
`setFont('5-pix-wide.bin')`.

Something else to notice is that `setFont` can take three more parameters
(character width, character height and gap width) for compatibility with the
Thumby API. But the gap width (or "space") will now default to `1`, which is
usually what it's set to anyway.

When loading variable width fonts in the FancyFont format, you must not give
`setFont` a character width and height. That will tell it that you're loading a
variable width font, and load the dimensions from the file. This is what happens
in the line that says `setFont('5-pix-wide.bin')`.

### `drawText`

The `drawText` function seems to behave similarly to the Thumby API version at
first, but there are some differences. First, we can omit the color parameter if
we want to draw in white. Second, it will do the right thing if you give it a
string with a newline character (`\n`) as shown in the example on line 17.
Third, this version of `drawText` accepts two more parameters to tell it to clip
the string at some row and column.

### `drawTextWrapped`

`drawTextWrapped` accepts the same inputs as `drawText`, the difference being
that `drawTextwrapped` will treat the clipping column as the point to word-wrap
the string instead.

### Constructor

Finally, both `drawText` and `drawTextWrapped` draw to the buffer that you
specify in the FancyFont constructor. You must give the constructor a buffer
(generally a bytearray) and optionally a width and a height, which default to 72
and 40. For normal operation this mean you just do:

```python
fancyFont = FancyFont(thumbyGraphics.display.display.buffer)
  # or
fancyFont = FancyFont(thumby.display.display.buffer)
```

But it gives you some flexibility to draw into another buffer, for example for
use with the grayscale library. It also probably means that this library is not
really a Thumby library, but rather a generic MicroPython library, and could
probably be used for projects on other MicroPython based platforms too.

## Minified version

Note that there is also a file called
[`fancyFont-minified.py`](./fancyFont-minified.py) in this repository, which is
functionally the same library, but without all the comments and documentation:

```
   15K  fancyFont.py
  4.1K  fancyFont-minified.py
```

It makes quite a bit of difference, so you may want to ship your project with
the minified version of the library.

## API documentation

### Class FancyFont

A container class that holds functions for font rendering for fixed width and
variable width fonts, with optional word-wrapping.

Methods:

  * `__init__`
  * `setFont`
  * `drawText`
  * `drawTextWrapped`

Attributes:

  * None that should be manipulated by the user

#### `__init__(displayBuffer, displayWidth = 72, displayHeight = 40)`

Constructor function to initialize the FancyFont class.

Parameters:

  * `displayBuffer` : `object`
    * The display buffer to draw to. Usually
      thumbyGraphics.display.display.buffer.

  * `displayWidth` : `int`
    * The width of the display to draw to, in pixels. Usually
      thumbyGraphics.display.width. Defaults to 72.

  * `displayHeight` : `int`
    * The height of the display to draw to, in pixels. Usually
      thumbyGraphics.display.height. Defaults to 40.

#### `setFont(fontPath, width:int = None, height:int = None, space:int = 1)`

Set the font file at `fontPath` as the current font to be used for all
subsequent drawText commands.

Parameters:

  * `fontPath` : `string`
      * A path to a file that contains a font in either the TinyCircuits
        fixed width font file format or a FancyFont variable width font
        file.

  * `width` : `int`
      * The character width of the font, if the font is fixed width. Omit or
        supply `None` for variable width. Note that characters with a width
        of more than 8 pixels are *not supported*.

  * `height` : `int`
      * The character height of the font, if the font is fixed width. Omit
        or supply `None` for variable width font files. The character height
        will then be read from the font file. Note that characters with a
        height of more than 8 pixels are *not supported*.

  * `space` : `int`
      * The margin between characters for fixed width fonts. Defaults to 1.
        Ignored for variable width fonts.

#### `drawText(string, xPos:int, yPos:int, color:int = 1, xMax:int = None, yMax:int = None)`

Draw a string within the square defined by (xPos, yPos) and (xMax, yMax), in
the given color.

Parameters:

  * `string` : `string`
      * The string to draw to the screen.

  * `xPos` : `int`
      * The X coordinate to start drawing from, counting from the left side
        of the screen.

  * `yPos` : `int`
      * The Y coordinate to start drawing from, counting from the top of the
        screen.

  * `color` : `int`
      * The color to draw the string in: either 1 (white) or 0 (black).

  * `xMax` : `int`
      * The X coordinate to stop drawing from, counting from the left side
        of the screen. Any text wider than xMax - xPos will be clipped.
        Defaults to the display width supplied to the constructor.

  * `yMax` : `int`
      * The Y coordinate to stop drawing from, counting from the top of the
        screen. Any line of text higher than yMax - yPos will be clipped.
        Defaults to the display height supplied to the constructor.


#### `drawTextWrapped(string, xPos:int, yPos:int, color:int = 1, xMax:int = None, yMax:int = None)`

Draw a string within the square defined by (xPos, yPos) and (xMax, yMax), in
the given color with word wrapping.

Parameters:

  * `string` : `string`
      * The string to draw to the screen.

  * `xPos` : `int`
      * The X coordinate to start drawing from, counting from the left side
        of the screen.

  * `yPos` : `int`
      * The Y coordinate to start drawing from, counting from the top of the
        screen.

  * `color` : `int`
      * The color to draw the string in: either 1 (white) or 0 (black).

  * `xMax` : `int`
      * The X coordinate to wrap the text at, counting from the left side of
        the screen. Defaults to the display width supplied to the
        constructor.

  * `yMax` : `int`
      * The Y coordinate to stop drawing from, counting from the top of the
        screen. Any text higher than yMax - yPos will be clipped. Defaults
        to the display height supplied to the constructor.
