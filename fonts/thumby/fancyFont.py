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
  def drawTextWrapped(self, string:ptr8, xPos:int, yPos:int, color:int = 1, xMax:int = display.width, yMax:int = display.height):
    wrappedString = self._wrapText(string, xPos, yPos, xMax, yMax)
    self._drawText(wrappedString, xPos, yPos, color, xMax, yMax)

  # Draw a string within the square defined by (xPos, yPos) and (xMax, yMax), in
  # the given color. This wrapper function is here because viper functions can't
  # have a variable number of arguments
  @micropython.native
  def drawText(self, string:ptr8, xPos:int, yPos:int, color:int = 1, xMax:int = display.width, yMax:int = display.height):
    self._drawText(string, xPos, yPos, color, xMax, yMax)

  @micropython.viper
  def _drawText(self, string:ptr8, xStart:int, yPos:int, color:int, xMax:int, yMax:int):
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

    if characterWidth == VARIABLE_WIDTH:
      characterIndices  = self.characterIndices

    while string[stringIndex] != 0:

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

      # Are we outside screen bounds? Then don't draw anything, just count
      if xPos + currentWidth <= 0 or xPos >= xMax or yPos + characterHeight <= 0:
        xPos += currentWidth + characterMarginWidth
        continue
      # Are we below where we're allowed to draw? Then we're really done
      if yPos >= yMax:
        break;

      # Can we draw the full character height?
      heightMask = 0xFF
      if yPos + characterHeight > yMax:
        heightMask >>= (8 - (yMax - yPos))   # Nope; just the top part

      # Can we draw the full character width?
      if xPos + currentWidth > xMax:
        blitWidth = xMax - xPos     # Nope; just the left part
      else:
        blitWidth = currentWidth

      # Blit the character bitmap to the display buffer in the right place
      vertOffset = yPos & 7 # y % 8
      dispBufIndex = (yPos >> 3) * screenWidth + xPos
      if color == 0:
        for x in range(0, blitWidth):
          fontByte = characterBufferPtr[characterOffset + x] & heightMask
          displayBuffer[dispBufIndex + x] &= 0xFF ^ (fontByte << vertOffset)
          displayBuffer[dispBufIndex + x + screenWidth] &= 0xFF ^ (fontByte >> (8 - vertOffset))
      else:
        for x in range(0, blitWidth):
          fontByte = characterBufferPtr[characterOffset + x] & heightMask
          displayBuffer[dispBufIndex + x] |= fontByte << vertOffset
          displayBuffer[dispBufIndex + x + screenWidth] |= fontByte >> (8 - vertOffset)

      xPos += currentWidth + characterMarginWidth

  @micropython.viper
  def _wrapText(self, string, xStart:int, yPos:int, xMax:int, yMax:int):
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
    output = bytearray(int(len(string)))
    outputPtr:ptr8 = ptr8(output)
    stringPtr:ptr8 = ptr8(string)

    while stringPtr[stringIndex] != 0:

      # Fetch character from string
      character = stringPtr[stringIndex]
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
        while stringIndex > 0 and stringPtr[stringIndex] != SPACE:
          stringIndex -= 1
        # Did we go all the way back to the previous break, or can we add a NEWLINE?
        if stringIndex == 0 or stringPtr[stringIndex] == NEWLINE:
            stringIndex = problemIndex
        else:
            outputPtr[stringIndex] = NEWLINE
        stringIndex += 1
        yPos += characterHeight + 1
        xPos = xStart
        continue

      xPos += currentWidth + characterMarginWidth

    return output

# Make an instance of the class available for importing
fancyFont = FancyFont()
