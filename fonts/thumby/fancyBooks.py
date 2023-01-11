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
  "threesquare": '/Games/FancyBooks/3-pix.bin',
  "Splendid": '/Games/FancyBooks/4-pix-low.bin',
  "Black sheep": '/Games/FancyBooks/4-pix-high.bin',
  "Limited narrow": '/Games/FancyBooks/5-pix-narrow.bin',
  "BigBoy": '/Games/FancyBooks/5-pix-wide.bin',
  "Truthful": '/Games/FancyBooks/6-pix.bin'
}

books = {
  "fox": "the quick brown fox jumped over the lazy dog\nTHE QUICK BROWN FOX JUMPED OVER THE LAZY DOG\n0123456789\n!\"#$%&'()*+,-./:; <=>?@[\]^_`{|}~",
  "lorem": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam nec nisi a mi interdum placerat. Vivamus sed tincidunt risus. Fusce egestas et lectus at pretium. Donec dictum blandit libero.",
  "wikipedia": "The Thumby is a small keychain sized programmable game console produced by TinyCircuits of Akron, Ohio[3][4] and funded by a Kickstarter campaign.[5][6][7] The console measures 1.2 by 0.7 by 0.3 inches (30.5 mm * 17.8 mm * 7.6 mm)."
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
