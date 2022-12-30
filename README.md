# CHIPCODE — Tools for writing CHIP-8 code

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
