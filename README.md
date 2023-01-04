# CHIPCODE — Tools for writing CHIP-8 code

## Octo assembler

Just a very thin wrapper around [John Earnest](https://github.com/JohnEarnest)'s
excellent [Octo](https://github.com/JohnEarnest/Octo)-flavoured CHIP-8 assembler
and disassembler. Just because it bugged me not to be able to depend on it
through NPM.

  * [Octo assembler on NPM](https://www.npmjs.com/package/@chipcode/octo-assembler)
  * [Octo assembler on Github](./octo-assembler)

## Octopus

Preprocessing for Octo-flavoured CHIP-8 source files. This command line tool can
make conditional code inclusion decisions, merge multiple source files and
re-order your source to put the executable stuff first.

  * [Octopus on NPM](https://www.npmjs.com/package/@chipcode/octopus)
  * [Octopus on Github](./octopus)

## Image loader

Conversion of images into sprite data with labels and bytes, in a format that is
compatible with Octo. This tool can be used on the command line or as a plugin
to [Octopus](https://www.npmjs.com/package/@chipcode/octopus).

  * [Image loader on NPM](https://www.npmjs.com/package/@chipcode/image-loader)
  * [Image loader on Github](./image-loader)

## Syntax highlighting for VSCode

Syntax highlighting for Visual Studio Code and VSCodium for the
[Octo](http://octo-ide.com/) (*.8o) programming language. This plugin also
supports the additional instructions in
[Octopus](https://www.npmjs.com/package/@chipcode/octopus).

  * [Installation instructions](./vscode-syntax-highlighting#readme)

## CHIPCODE Fonts (text rendering)

The long-missing CHIP-8 text rendering library you've been wishing existed
already so you didn't have to write it! 😄

  * [CHIPCODE fonts on NPM](https://www.npmjs.com/package/@chipcode/fonts)
  * [CHIPCODE fonts on Github](./fonts)
