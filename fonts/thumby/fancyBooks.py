# TODO:
#  * Actually load in text files
#  * Allow selecting book more elegantly
#  * Maybe just an overall settings screen
#  * Allow changing contrast
#  * Figure out how to make offset not font size dependent
#  * Store settings per book

# Fix import path so it finds our modules above all else
import sys
sys.path.insert(0, '/Games/FancyBooks')

from fancyFont import fancyFont
from time import sleep_ms
import thumby
thumby.display.setFPS(0)

fonts = {
  "BigBoy": '/Games/FancyBooks/5-pix-wide.bin'
}

books = {
  "fox": "the quick brown fox jumped over the lazy dog\nTHE QUICK BROWN FOX JUMPED OVER THE LAZY DOG\n0123456789\n,./:&'-!?",
  "lorem": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam nec nisi a mi interdum placerat. Vivamus sed tincidunt risus. Fusce egestas et lectus at pretium. Donec dictum blandit libero.",
  "wikipedia": "CHIP-8 is an interpreted programming language, developed by Joseph Weisbecker made on his 1802 Microprocessor. It was initially used on the COSMAC VIP and Telmac 1800 8-bit microcomputers in the mid-1970's."
}

offsets = {
  "fox": 0,
  "lorem": 0,
  "wikipedia": 0
}

selectedColor = 0
selectedFont = "BigBoy"
selectedBook = "fox"

def waitKeyRelease():
  while thumby.inputPressed():
    pass

while True:
  
  offset = offsets[selectedBook]
  
  thumby.display.fill(selectedColor)
  fancyFont.setVariableWidthFont(fonts[selectedFont])
  # fancyFont.setFixedWidthFont("/lib/font3x5.bin", 3, 5)
  fancyFont.drawTextWrapped(
    books[selectedBook],
    selectedColor, # Left margin
    -1 * offset + selectedColor, # Top margin
    selectedColor ^ 1,  # Actual color
    thumby.display.width - selectedColor, # Right margin
    thumby.display.height - selectedColor # Bottom margin
  )
  thumby.display.update()
  
  # Wait for key release
  while True:
    if thumby.buttonA.pressed():
      waitKeyRelease()
      selectedColor ^= 1
      break
    if thumby.buttonB.pressed():
      waitKeyRelease()
      titles = list(books.keys())
      chosen = titles.index(selectedBook)
      selectedBook = titles[(chosen + 1) % len(books)]
      break
    if thumby.buttonL.pressed():
      waitKeyRelease()
      fontNames = list(fonts.keys())
      chosen = fontNames.index(selectedFont)
      selectedFont = fontNames[(chosen - 1) % len(fonts)]
      break
    if thumby.buttonR.pressed():
      waitKeyRelease()
      fontNames = list(fonts.keys())
      chosen = fontNames.index(selectedFont)
      selectedFont = fontNames[(chosen + 1) % len(fonts)]
      break
    if thumby.buttonU.pressed():
      if offset > 0:
          offsets[selectedBook] -= 1
      break
    if thumby.buttonD.pressed():
      offsets[selectedBook] += 1
      break
    
    # Sleep a bit in an attempt to be more battery efficient..?
    sleep_ms(50)
