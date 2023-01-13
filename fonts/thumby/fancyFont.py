# Note: characters can't be larger than 8 x 8 pixels

# TODO:
#  * Stability issues (on wide words..?)

from thumbyGraphics import display
from os import stat

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

  # Load a fixed width font. Accepts "regular Thumby fonts"
  @micropython.native
  def setFixedWidthFont(self, fontPath, width:int, height:int, space:int = 1):
    self.characterWidth = width
    self.characterHeight = height
    self.characterMarginWidth = space
    self.fontFile = open(fontPath)
    self.characterBuffer = bytearray(8)
    self.numCharactersInFont = stat(fontPath)[6] // self.characterWidth

  # Load a variable width font. Accepts fonts in "FancyFont format"
  @micropython.native
  def setVariableWidthFont(self, fontPath):
    self.characterWidth = VARIABLE_WIDTH
    self.characterMarginWidth = 0
    self.fontFile = open(fontPath)
    self.characterBuffer = bytearray(9)
    self.fontFile.readinto(self.characterBuffer)
    self.characterHeight = self.characterBuffer[0]
    self.numCharactersInFont = self.characterBuffer[1]
    self._collectCharacterIndices()

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

  # Draw a string within the square defined by (xPos, yPos) and (xMax, yMax), in
  # the given color with word wrapping
  @micropython.native
  def drawTextWrapped(self, string, xPos:int, yPos:int, color:int = 1, xMax:int = display.width, yMax:int = display.height):
    wrappedString = self._wrapText(string, len(string), xPos, yPos, xMax, yMax)
    return self._drawText(wrappedString, len(wrappedString), xPos, yPos, color, xMax, yMax)

  # Draw a string within the square defined by (xPos, yPos) and (xMax, yMax), in
  # the given color. This wrapper function is here because viper functions can't
  # have a variable number of arguments
  @micropython.native
  def drawText(self, string, xPos:int, yPos:int, color:int = 1, xMax:int = display.width, yMax:int = display.height):
    return self._drawText(string, len(string), xPos, yPos, color, xMax, yMax)

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
    displayBuffer:ptr8       = ptr8(display.display.buffer)
    fontFile                 = self.fontFile
    characterBuffer          = self.characterBuffer
    characterBufferPtr:ptr8  = ptr8(characterBuffer)
    screenWidth:int          = int(display.width)
    characterWidth:int       = int(self.characterWidth)
    characterHeight:int      = int(self.characterHeight)
    characterMarginWidth:int = int(self.characterMarginWidth)
    numCharactersInFont:int  = int(self.numCharactersInFont)
    dispBufSize:int          = int(len(display.display.buffer))

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
          fontByte = characterBufferPtr[characterOffset + x]
          if 0 <= dispBufIndex + x < dispBufSize:
            displayBuffer[dispBufIndex + x] &= 0xFF ^ (fontByte << vertOffset)
          if 0 <= dispBufIndex + x + screenWidth < dispBufSize:
            displayBuffer[dispBufIndex + x + screenWidth] &= 0xFF ^ (fontByte >> (8 - vertOffset))
      else:
        for x in range(0, blitWidth):
          fontByte = characterBufferPtr[characterOffset + x]
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

# Make an instance of the class available for importing
fancyFont = FancyFont()
