# Octopus

This command line tool can do pre-processing on text based input files. It is
intended for use with Octo-flavoured CHIP-8 assembly code, but it can be used
with any text file.

## Installing and running

Install Octopus as an NPM package, either for a project or globally:

```bash
npm install --save-dev @chipcode/octopus
```

You can then run Octopus on the command line:

```bash
npx octopus <input file> <output file> <option 1> <option 2>
```

Or you can use Octopus in your `package.json` file, for example to build a
Cosmac VIP and a SUPERCHIP version of your project:

```json
{
  "scripts": {
    "build": "npm run build-vip && npm run build-schip",
    "build-vip": "octopus ./src/index.8o ./bin/project-vip.8o VIP",
    "build-schip": "octopus ./src/index.8o ./bin/project-schip.8o SCHIP"
  }
}
```

Which you can then run using:

```bash
npm run build
```

## Features

### Conditional code inclusion

With Octopus, you can use `:if`, `:unless`, `:else` and `:end` in your code to
switch on options, like so:

```octo
: store-values
  i := target
  save v2      # This doesn't increment i on SCHIP
  :if SCHIP    # Conditionally fix that issue
    vF := 3
    i += vF
  :end
  return
```

Now the conditional code between `:if` and `:end` will be included in the target
file only if the option `SCHIP` is set.

Note that you can not use expressions, and options can only be either set (true)
or unset (false). To set an option, give it as a parameter to the Octopus
invocation or set it with `:const` (see below).

Here is a more complicated example, also showing the use of `:else` and
`:unless` (which is basically "if not").

```octo
: store-values
  i := target
  :if XOCHIP
    save v3 - v4   # This XO-CHIP opcode doesn't increment i
  :else
    v0 := v3
    v1 := v4
    save v1        # This doesn't increment i on SUPERCHIP
  :end
  :unless VIP      # So on anything but VIP, we need to increment i manually
    vF := 2
    i += vF
  :end
  return
```

### Set, reset and dump options

As described, options can be either set (true) or unset (false). You can set
options by providing them to the Octopus invocation on the command line, or in
your code with `:const`. If you set a constant to zero, it will be unset (false)
from the perspective of Octopus. Any other value will be considered set (true).

```octo
  :if VIP
    :const USE_SCROLLING 0
  :else
    :const USE_SCROLLING 1
  :end

  # ...

  :if USE_SCROLLING
    scroll-down 8
  :end
```

When you're playing with options, and you run into issues, you can use the
`:dump-options` command to instruct Octopus to output all those options that are
set at that point in the program.

```octo
  :const OPTION_1 1
  :const OPTION_2 0
  :dump-options

  :const OPTION_2 1
  :dump-options
```

This will output:

```
Options on line 3: [ 'OPTION_1' ]
Options on line 6: [ 'OPTION_1', 'OPTION_2' ]
```

### Include other files

When a project gets too large for a single source file, it is nice to be able to
split it up into more logical segments. The Octopus `:include` command allows
you to include another file into the current source file.

```octo
  :include "renderer.8o"
  :include "images/bitmaps.8o"
```

### Automatic re-ordering of code

When writing XO-CHIP code, you need to keep the code that executes in the first
3.5K of memory. Beyond that you can have another 60K of data. This is because
XO-CHIP does provide an instruction to load a 16-bit address into i (`i := long
<label>`), but no instructions to jump to or call a 16-bit label.

However, when writing code it is much more convenient to be able to keep the
code and the data that the code operates on close together. Octopus can
automatically solve this issue for you, if you annotate your code with `:segment
code` and `:segment data`.

Each file is considered to automatically start with `:segment code`, so you can
leave that out if your file does indeed start with executable code.

```octo
  i := long table
  load v4
  # Do something intelligent with data...

:segment data

: table
  0 1 2 3 4
```
