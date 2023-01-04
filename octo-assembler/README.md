# Octo assembler

This package is just a very thin wrapper around
[John Earnest](https://github.com/JohnEarnest)'s excellent
[Octo](https://github.com/JohnEarnest/Octo)-flavoured CHIP-8 assembler and
disassembler. I made this wrapper just because it bugged me not to be able to
depend on it through NPM.

## How to use

```bash
npm install @chipcode/octo-assembler
```

And then you can run:
```bash
npx octo <input file> <output file>
```

Or use it in your package.json file:
```json
{
  "name": "example",
  "scripts": {
    "assemble": "octo input.8o output.ch8",
    "disassemble": "octo --decompile input.ch8 output.8o",
    "assemble-all": "echo 'Assembling all *.8o files...'; for file in `find . -type f -name \"*.8o\"`; do echo \"  * $file\"; target=${file%.8o}; octo $file ${target}.ch8; done"
  }
}
```

See https://github.com/JohnEarnest/Octo#command-line-mode for available command
line options to the assembler.
