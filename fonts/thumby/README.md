# FancyFont

This is a font library for the Thumby playable keychain console. It's in this
repository because I didn't think it really needed its own repository, and so it
can share font definitions with the CHIP-8 version of the library.

FancyFont is a drop-in replacement to the `setFont()` and `drawText()` functions in
the Thumby API. It has a couple of improvements over those functions:

  * Support for word wrapping with `drawTextWrapped()`
  * Support for clipping
  * Support for newline character (`\n`)
  * Support for variable width fonts
  * Can render to any graphics buffer, so "should work" with grayscale too

## Visual examples

TODO: Nice pictures here

## Code example

To render this to the display using FancyFont:

![The visual result of the code example below](./pictures/code-example.png)

You need to do something like this:

```python
# Fix import path so we can import files in our game directory
from sys import path
path.insert(0, '/Games/FancyBooks')

# Load either `thumby` or `thumbyGraphics` library
import thumbyGraphics

# Load and initialize `FancyFont` library
from fancyFont import FancyFont
fancyFont = FancyFont(thumbyGraphics.display.display.buffer)

# Clear the screen
thumbyGraphics.display.fill(0)

# Select a TinyCircuits font and draw some text
fancyFont.setFont('/lib/font5x7.bin', 5, 7)      # Gap may be specified as in
                                                 # Thumby API, but defaults to 1
fancyFont.drawText('FancyFont\n---------', 9, 4) # Color may be specified but
                                                 # defaults to white

# Select a FancyFont font and draw some wrapping text
fancyFont.setFont('5-pix-wide.bin')   # No need to specify character dimensions
                                      # for FancyFont fonts. Also: filename may
                                      # be relative to import path above.
fancyFont.drawTextWrapped('Hello there, Thumby user!', 8, 22)

# Show the result
thumbyGraphics.display.update()
```

## API documentation

For the full API documentation, see the very complete docstrings in FancyFont by
opening [fancyFont.py](./fancyFont.py) or by entering this in the Thumby REPL:

```python
from fancyFont import FancyFont
help(FancyFont)
```
