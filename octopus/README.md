## Octopus

This tool does pre-processing on text-based input files. It is intended for use
with Octo-flavoured CHIP-8 assembly code, but it can of course be used with any
text file.

Features:
* Include other files (`#include <filename>`)
* Make code inclusion decisions (`#if <option>` / `#else` / `#end`)
* Mark blocks as code or data, and get them re-ordered (`#code` / `#data`)

### Example

```octo
: main
  i := pointer
  load v2
  # Manipulate v0 - v2 here, or whatever

  # The value of i now depends on the interpreter running the binary. On Super
  # CHIP it still points to `pointer`, whereas in most other interpreters it
  # will point to `pointer` plus three. Let's fix that with a conditional:
  #if SCHIP
    vF := 3
    i += vF
  #end
  save v2

  # ...

# When writing programs for XO-CHIP, you want your data to end up at the end of
# memory and your code at the beginning. Because only the first 3.5K of memory
# is executable. But we don't always want our source files to be structured in
# that way. Using Octopus, we can mark a section as data, and it will
# automatically be moved down:

#data
: table
  0x00 0x01
  0x02 0x03
#code

# If your project grows, you don't want to keep everything in a single file
# anymore. Octopus allows you to recursively include files:

#include "map-rendering.8o"

# Maybe, depending on some variable, you want to load a different version of a
# file. That's possible too:

#if SCHIP
  #include "decompress-schip.8o"
#else
  #include "decompress-vip.8o"
#end
```

The general syntax to run Octopus is:
```bash
npx octopus <input file> <ouput file> <option 1> <option 2> ...
```

So given the example above, to produce the two different versions of your code:

```bash
npx octopus main.8o mygame.8o
npx octopus main.8o mygame-schip.8o SCHIP
```
