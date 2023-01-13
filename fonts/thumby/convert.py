# pip3 install pillow
from PIL import Image
from sys import argv, exit

# This is a (very loose) Python port of ../../image-loader/index.js

# Cut up image in sprites of the selected resolution
def splitIntoSprites(image, spriteWidth, spriteHeight):
  sprites = []
  for y in range(0, image.height - image.height % spriteHeight, spriteHeight):
    for x in range(0, image.width - image.width % spriteWidth, spriteWidth):
      box = (x, y, x+spriteWidth, y+spriteHeight)
      sprites.append(image.crop(box).tobytes())
  return sprites

# Visualise 1-bit bitmap
def render(image, width):
  offset = 0
  output = ''
  height = len(image) // (width // 8)
  for y in range(0, height):
    for byte in image[offset:(offset + width // 8)]:
      output += "{:08b}".format(byte).replace('1', '██').replace('0', '  ')
    output += '\n'
    offset += width // 8
  return output

# Program start

if len(argv) < 4:
  print("Missing required parameters.\n\nUsage:\n  python3 convert.py <input file> <output file> <character height>")
  exit()

inputFile = argv[1]
outputFile = argv[2]
characterHeight = int(argv[3])
characterWidth = 8

image = Image.open(inputFile).convert('1', dither=Image.Dither.NONE)
# print(render(image.tobytes(), image.width))

characters = splitIntoSprites(image, characterWidth, characterHeight + 1)
binary = bytearray([characterHeight, len(characters)])
for character in characters:
  # print(render(character, characterWidth))
  width = character[0]
  charImage = Image.frombytes('1', (8, 8), bytes(list(character)[1:] + [0] * 8))
  charImage = charImage.rotate(270)
  charImage = charImage.crop((0, 0, 8, width))
  binary.append(width)
  binary += charImage.tobytes()

file = open(outputFile, 'wb')
file.write(binary)
file.close()