# Octo(pus) syntax highlighting

This plugin provides syntax highlighting for Visual Studio Code and VSCodium for
the Octo (*.8o) programming language. [Octo](http://octo-ide.com/) is the
primary assembly syntax used in writing CHIP-8 games. This plugin also supports
the additional instructions in
[Octopus](https://www.npmjs.com/package/@chipcode/octopus).

## Installing

The easiest way to install this plugin is to download [the VSIX
file](https://github.com/Timendus/chipcode/raw/main/vscode-syntax-highlighting/dist/octopus-syntax-highlighting.vsix)
and install it in your IDE:

  * Go to the extensions tab on the left
  * Click on the three dots at the top right of the extensions panel
  * Select "Install from VSIX..."
  * Select the downloaded file

## Release Notes

### 0.0.1

Initial release of this plugin. Expect issues 😉

### 0.0.2

Change `:data` and `:code` to `:segment data` and `:segment code`, like they
should be. Add `:dump-options` to the list.

### 0.0.3

Add support for `:monitor`, `:assert`, `:byte`, `:macro`, `:calc`,
`:stringmode`, `:pointer`, `:call` and `pitch`

## Acknowledgements

This package is an adapted subset of [Cody Hoover's
plugin](https://github.com/hoovercj/vscode-octo). Octo is a [project by John
Earnest](https://github.com/JohnEarnest/Octo).
