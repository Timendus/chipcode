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

from fancyFont import FancyFont
from time import sleep_ms
import thumby
thumby.display.setFPS(0)
fancyFont = FancyFont(thumby.display.display.buffer)

AUTO = None

fonts = [
  {
    'file': '/Games/FancyBooks/3-pix.bin',
    'name': 'threesquare',
    'width': AUTO,
    'height': AUTO
  },
  {
    'file': '/Games/FancyBooks/4-pix-low.bin',
    'name': 'Splendid',
    'width': AUTO,
    'height': AUTO
  },
  {
    'file': '/Games/FancyBooks/4-pix-high.bin',
    'name': 'Black sheep',
    'width': AUTO,
    'height': AUTO
  },
  {
    'file': '/Games/FancyBooks/5-pix-narrow.bin',
    'name': 'Limited narrow',
    'width': AUTO,
    'height': AUTO
  },
  {
    'file': '/lib/font3x5.bin',
    'name': 'TC small font 3x5',
    'width': 3,
    'height': 5
  },
  {
    'file': '/Games/FancyBooks/5-pix-wide.bin',
    'name': 'BigBoy',
    'width': AUTO,
    'height': AUTO
  },
  {
    'file': '/Games/FancyBooks/6-pix.bin',
    'name': 'Truthful',
    'width': AUTO,
    'height': AUTO
  },
  {
    'file': '/lib/font5x7.bin',
    'name': 'TC medium font 5x7',
    'width': 5,
    'height': 7
  },
  {
    'file': '/lib/font8x8.bin',
    'name': 'TC large font 8x8',
    'width': 8,
    'height': 8
  }
]

books = [
  {
    'name': 'fox',
    'content': "the quick brown fox jumps over the lazy dog\nTHE QUICK BROWN FOX JUMPS OVER THE LAZY DOG\n0123456789\n!\"#$%&'()*+,-./:; <=>?@[\]^_`{|}~",
    'offset': 0
  },
  {
    'name': 'lorem',
    'content': "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam nec nisi a mi interdum placerat. Vivamus sed tincidunt risus. Fusce egestas et lectus at pretium. Donec dictum blandit libero.",
    'offset': 0
  },
  {
    'name': 'wikipedia',
    'content': "The Thumby is a small keychain sized programmable game console produced by TinyCircuits of Akron, Ohio[3][4] and funded by a Kickstarter campaign.[5][6][7] The console measures 1.2 by 0.7 by 0.3 inches (30.5 mm * 17.8 mm * 7.6 mm).",
    'offset': 0
  },
]

selectedColor = 0
selectedFont = 5
selectedBook = 0

def waitKeyRelease():
  while thumby.inputPressed():
    pass

def setSelectedFont():
  fancyFont.setFont(
    fonts[selectedFont]['file'],
    fonts[selectedFont]['width'],
    fonts[selectedFont]['height']
  )
  
setSelectedFont()

while True:
  
  offset = books[selectedBook]['offset']
  
  thumby.display.fill(selectedColor)
  pos = fancyFont.drawTextWrapped(
    books[selectedBook]['content'],
    selectedColor, # Left margin
    -1 * offset + selectedColor, # Top margin
    selectedColor ^ 1,  # Actual color
    thumby.display.width - selectedColor, # Right margin
    thumby.display.height - selectedColor # Bottom margin
  )
  # print(int(pos[0]), int(pos[1]))
  thumby.display.update()
  
  # Wait for key release
  while True:
    if thumby.buttonA.pressed():
      waitKeyRelease()
      selectedColor ^= 1
      break
    if thumby.buttonB.pressed():
      waitKeyRelease()
      selectedBook = (selectedBook + 1) % len(books)
      break
    if thumby.buttonL.pressed():
      waitKeyRelease()
      selectedFont = (selectedFont - 1) % len(fonts)
      setSelectedFont()
      break
    if thumby.buttonR.pressed():
      waitKeyRelease()
      selectedFont = (selectedFont + 1) % len(fonts)
      setSelectedFont()
      break
    if thumby.buttonU.pressed():
      if offset > 0:
        books[selectedBook]['offset'] -= 1
      break
    if thumby.buttonD.pressed():
      books[selectedBook]['offset'] += 1
      break
    
    # Sleep a bit in an attempt to be more battery efficient..?
    sleep_ms(50)
