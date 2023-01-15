"""
Font rendering for fixed width and variable width fonts, with optional
word-wrapping.

Classes:

    FancyFont
"""

# TODO:
#  * Wrapping consistency issue (off by one on the width somewhere?)
#  * Idea: allow loading fonts as bytearray too?

from os import stat
from sys import path

VARIABLE_WIDTH = const(0)
NEWLINE        = const(10)
SPACE          = const(32)

# For debugging

# def dump(barray):
#   result = ''
#   for i in barray:
#     result += "0x%02x" % i + ' '
#   return result

class FancyFont:
  """
  A container class that holds functions for font rendering for fixed width and
  variable width fonts, with optional word-wrapping.

  Methods:

      __init__
      setFont
      drawText
      drawTextWrapped

  Attributes:

      None that should be manipulated by the user
  """

  @micropython.native
  def __init__(self, displayBuffer, displayWidth = 72, displayHeight = 40):
    """
    Constructor function to initialize the FancyFont class.

    Parameters:

        displayBuffer : object
            The display buffer to draw to. Usually
            thumbyGraphics.display.display.buffer.

        displayWidth : int
            The width of the display to draw to, in pixels. Usually
            thumbyGraphics.display.width. Defaults to 72.

        displayHeight : int
            The height of the display to draw to, in pixels. Usually
            thumbyGraphics.display.height. Defaults to 40.
    """
    self.displayBuffer = displayBuffer
    self.displayWidth = int(displayWidth)
    self.displayHeight = int(displayHeight)

  @micropython.native
  def setFont(self, fontPath, width:int = None, height:int = None, space:int = 1):
    """
    Set the font file at `fontPath` as the current font to be used for all
    subsequent drawText commands.

    Parameters:

        fontPath : string
            A path to a file that contains a font in either the TinyCircuits
            fixed width font file format or a FancyFont variable width font
            file. Path may be absolute or relative to any entry in `sys.path`.

        width : int
            The character width of the font, if the font is fixed width. Omit or
            supply `None` for variable width. Note that characters with a width
            of more than 8 pixels are *not supported*.

        height : int
            The character height of the font, if the font is fixed width. Omit
            or supply `None` for variable width font files. The character height
            will then be read from the font file. Note that characters with a
            height of more than 8 pixels are *not supported*.

        space : int
            The margin between characters for fixed width fonts. Defaults to 1.
            Ignored for variable width fonts.
    """
    fontPath = self._findFile(fontPath)
    self.fontFile = open(fontPath, 'rb')
    self.characterBuffer = bytearray(9)

    if width == None and height == None:
      # Assume variable width FancyFont file format
      self.characterWidth = VARIABLE_WIDTH
      self.characterMarginWidth = 0
      self.fontFile.readinto(self.characterBuffer)
      self.characterHeight = self.characterBuffer[0]
      self.numCharactersInFont = self.characterBuffer[1]
      self._collectCharacterIndices()
    else:
      # Assume fixed width TinyCircuits font file format
      self.characterWidth = width
      self.characterHeight = height
      self.characterMarginWidth = space
      self.numCharactersInFont = stat(fontPath)[6] // self.characterWidth

  # This wrapper function is here because viper functions can't have a variable
  # number of arguments
  @micropython.native
  def drawText(self, string, xPos:int, yPos:int, color:int = 1, xMax:int = None, yMax:int = None):
    """
    Draw a string within the square defined by (xPos, yPos) and (xMax, yMax), in
    the given color.

    Parameters:

        string : string
            The string to draw to the screen.

        xPos : int
            The X coordinate to start drawing from, counting from the left side
            of the screen.

        yPos : int
            The Y coordinate to start drawing from, counting from the top of the
            screen.

        color : int
            The color to draw the string in: either 1 (white) or 0 (black).

        xMax : int
            The X coordinate to stop drawing from, counting from the left side
            of the screen. Any text wider than xMax - xPos will be clipped.
            Defaults to the display width supplied to the constructor.

        yMax : int
            The Y coordinate to stop drawing from, counting from the top of the
            screen. Any line of text higher than yMax - yPos will be clipped.
            Defaults to the display height supplied to the constructor.
    """
    return self._drawText(
      string,
      len(string),
      xPos,
      yPos,
      color,
      xMax or self.displayWidth,
      yMax or self.displayHeight
    )

  @micropython.native
  def drawTextWrapped(self, string, xPos:int, yPos:int, color:int = 1, xMax:int = None, yMax:int = None):
    """
    Draw a string within the square defined by (xPos, yPos) and (xMax, yMax), in
    the given color with word wrapping.

    Parameters:

        string : string
            The string to draw to the screen.

        xPos : int
            The X coordinate to start drawing from, counting from the left side
            of the screen.

        yPos : int
            The Y coordinate to start drawing from, counting from the top of the
            screen.

        color : int
            The color to draw the string in: either 1 (white) or 0 (black).

        xMax : int
            The X coordinate to wrap the text at, counting from the left side of
            the screen. Defaults to the display width supplied to the
            constructor.

        yMax : int
            The Y coordinate to stop drawing from, counting from the top of the
            screen. Any text higher than yMax - yPos will be clipped. Defaults
            to the display height supplied to the constructor.
    """
    wrappedString = self._wrapText(
      string,
      len(string),
      xPos,
      yPos,
      xMax or self.displayWidth,
      yMax or self.displayHeight
    )
    return self._drawText(
      wrappedString,
      len(wrappedString),
      xPos,
      yPos,
      color,
      xMax or self.displayWidth,
      yMax or self.displayHeight
    )

  # Search for the requested file in all directories that are in `path`,
  # starting with interpreting it as an absolute path.
  def _findFile(self, filePath):
    try:
      stat(filePath)
      return filePath
    except OSError:
      pass
    for p in path:
      try:
        fullPath = (p + '/' if p and not p.endswith('/') else p) + filePath
        stat(fullPath)
        return fullPath
      except OSError:
        pass
    raise OSError('Font file not found')

  # Read through the file and cache the starting indices for all the characters
  @micropython.viper
  def _collectCharacterIndices(self):
    self.characterIndices = []
    currentIndex:int = 2
    for i in range(int(self.numCharactersInFont)):
      self.characterIndices.append(currentIndex)
      self.fontFile.seek(currentIndex)
      self.fontFile.readinto(self.characterBuffer)
      currentIndex += int(self.characterBuffer[0]) + 1

  @micropython.viper
  def _drawText(self, string:ptr8, strLen:int, xStart:int, yPos:int, color:int, xMax:int, yMax:int):
    # Cast the input parameters (because type-hints are really just for show...)
    string:ptr8 = ptr8(string)
    strLen:int = int(strLen)
    xStart:int = int(xStart)
    yPos:int = int(yPos)
    xMax:int = int(xMax)
    yMax:int = int(yMax)

    # Define variables up front so we have a stable memory profile in the loop
    stringIndex:int     = 0
    characterOffset:int = 0
    character:int       = 0
    currentWidth:int    = 0
    vertOffset:int      = 0
    dispBufIndex:int    = 0
    fontByte:int        = 0
    blitWidth:int       = 0
    heightMask:int      = 0
    xPos:int            = xStart

    # Track the rightmost and lowest pixel that we draw to for return value
    bottom:int          = xPos
    right:int           = yPos

    # Look up and cast all variables up front so we're faster in the loop
    displayBuffer:ptr8       = ptr8(self.displayBuffer)
    fontFile                 = self.fontFile
    characterBuffer          = self.characterBuffer
    characterBufferPtr:ptr8  = ptr8(characterBuffer)
    screenWidth:int          = int(self.displayWidth)
    characterWidth:int       = int(self.characterWidth)
    characterHeight:int      = int(self.characterHeight)
    characterMarginWidth:int = int(self.characterMarginWidth)
    numCharactersInFont:int  = int(self.numCharactersInFont)
    dispBufSize:int          = int(len(self.displayBuffer))

    if characterWidth == VARIABLE_WIDTH:
      characterIndices  = self.characterIndices

    while stringIndex < strLen:

      # Fetch character from string
      character = string[stringIndex]
      stringIndex += 1

      # Is this a newline?
      if character == NEWLINE:
        yPos += characterHeight + 1
        xPos = xStart
        continue

      # Convert ascii character to "index in the font"
      character = character - SPACE

      # Just ignore unprintable characters
      if not 0 <= character < numCharactersInFont:
        continue

      # Seek the font file to where the character is stored
      if characterWidth == VARIABLE_WIDTH:
        fontFile.seek(int(characterIndices[character]))
      else:
        fontFile.seek(character * characterWidth)

      # Load the character bitmap from file into our buffer
      fontFile.readinto(characterBuffer)
      
      # Load the width of this character (for variable width fonts)
      if characterWidth == VARIABLE_WIDTH:
        currentWidth = characterBufferPtr[0]
        characterOffset = 1
      else:
        currentWidth = characterWidth
        characterOffset = 0

      # Are we fully outside screen bounds? Then don't draw anything, just count
      if xPos + currentWidth <= 0 or xPos >= xMax or yPos + characterHeight <= 0:
        xPos += currentWidth + characterMarginWidth
        continue
      # Are we fully below where we're allowed to draw? Then we're really done
      if yPos >= yMax:
        break;

      # Can we draw the full character width?
      blitWidth = xMax - xPos        # What space do we have available?
      if blitWidth > currentWidth:
          blitWidth = currentWidth   # It's more than we need

      # Can we draw the full character height?
      heightMask = 0xFF
      if yPos + characterHeight > yMax:
        heightMask >>= (8 - (yMax - yPos))   # Nope; just the top part

      # Update drawn bounds
      right = int(max(right, xPos + blitWidth - 1))
      bottom = int(max(bottom, min(yPos + characterHeight - 1, yMax - 1)))

      # Blit the character bitmap to the display buffer in the right place. Note
      # that either the top or bottom byte to blit to may be outside the screen
      # bounds, so we have to check for that before blitting to make sure we
      # don't overwrite some random data (yay viper!)
      vertOffset = yPos & 7 # y % 8
      dispBufIndex = (yPos >> 3) * screenWidth + xPos
      if color == 0:
        for x in range(0, blitWidth):
          fontByte = characterBufferPtr[characterOffset + x] & heightMask
          if 0 <= dispBufIndex + x < dispBufSize:
            displayBuffer[dispBufIndex + x] &= 0xFF ^ (fontByte << vertOffset)
          if 0 <= dispBufIndex + x + screenWidth < dispBufSize:
            displayBuffer[dispBufIndex + x + screenWidth] &= 0xFF ^ (fontByte >> (8 - vertOffset))
      else:
        for x in range(0, blitWidth):
          fontByte = characterBufferPtr[characterOffset + x] & heightMask
          if 0 <= dispBufIndex + x < dispBufSize:
            displayBuffer[dispBufIndex + x] |= fontByte << vertOffset
          if 0 <= dispBufIndex + x + screenWidth < dispBufSize:
            displayBuffer[dispBufIndex + x + screenWidth] |= fontByte >> (8 - vertOffset)

      # Set the stage for the next character
      xPos += currentWidth + characterMarginWidth

    # Return the actual dimensions of the rendered text
    return bytearray([right, bottom])

  @micropython.viper
  def _wrapText(self, string:ptr8, strLen:int, xStart:int, yPos:int, xMax:int, yMax:int):
    # Cast the input parameters (because type-hints are really just for show...)
    string:ptr8 = ptr8(string)
    strLen:int = int(strLen)
    xStart:int = int(xStart)
    yPos:int = int(yPos)
    xMax:int = int(xMax)
    yMax:int = int(yMax)

    # Define variables up front so we have a stable memory profile in the loop
    stringIndex:int     = 0
    character:int       = 0
    currentWidth:int    = 0
    xPos:int            = xStart

    # Look up and cast all variables up front so we're faster in the loop
    fontFile                 = self.fontFile
    characterBuffer          = self.characterBuffer
    characterBufferPtr:ptr8  = ptr8(characterBuffer)
    characterWidth:int       = int(self.characterWidth)
    characterHeight:int      = int(self.characterHeight)
    characterMarginWidth:int = int(self.characterMarginWidth)
    numCharactersInFont:int  = int(self.numCharactersInFont)

    if characterWidth == VARIABLE_WIDTH:
      characterIndices  = self.characterIndices

    # Write result to a new string, so we don't mess with the original
    output = bytearray(strLen)
    outputPtr:ptr8 = ptr8(output)

    while stringIndex < strLen:

      # Fetch character from string
      character = string[stringIndex]
      outputPtr[stringIndex] = character
      stringIndex += 1

      # Is this a newline?
      if character == NEWLINE:
        yPos += characterHeight + 1
        xPos = xStart
        continue

      # Convert ascii character to "index in the font"
      character = character - SPACE

      # Just ignore unprintable characters
      if not 0 <= character < numCharactersInFont:
        continue

      # Seek the font file to where the character is stored
      if characterWidth == VARIABLE_WIDTH:
        fontFile.seek(characterIndices[character])
      else:
        fontFile.seek(character * characterWidth)

      # Load the character bitmap from file into our buffer
      fontFile.readinto(characterBuffer)

      # Load the width of this character (for variable width fonts)
      currentWidth = characterWidth
      if characterWidth == VARIABLE_WIDTH:
        currentWidth = characterBufferPtr[0]

      # Are we below where we're allowed to draw? Then we're really done
      if yPos >= yMax:
        break;

      # Are we overflowing the horizontal bounds? Then add in a NEWLINE
      if xPos + currentWidth >= xMax:
        stringIndex -= 1 # Undo stringIndex increase after loading character
        problemIndex = stringIndex
        while stringIndex > 0 and string[stringIndex] != SPACE:
          stringIndex -= 1
        # Did we go all the way back to the previous break, or can we add a NEWLINE?
        if stringIndex == 0 or outputPtr[stringIndex] == NEWLINE:
          stringIndex = problemIndex
        else:
          # Just as a guard:
          if 0 > stringIndex >= strLen:
            raise ValueError("Attempted to write outside string bounds: this should never happen!")
          outputPtr[stringIndex] = NEWLINE
        stringIndex += 1
        yPos += characterHeight + 1
        xPos = xStart
        continue

      xPos += currentWidth + characterMarginWidth

    return output
