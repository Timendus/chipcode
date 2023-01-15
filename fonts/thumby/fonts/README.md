# Font images

The images in the [`images`](./images) subdirectory are copied there from the
[CHIP-8 project](../../fonts) by the [build script](../build.sh). They are then
converted to `.bin` files in the [`fixed-width`](./fixed-width) and
[`variable-width`](./variable-width) directories using the
[`convert.py`](../convert.py) script.

## How to build

To convert all images in this directory, take a look at the
[`build.sh`](../build.sh) script.

To convert one variable width font image to a `.bin` file, run:

```bash
python3 convert.py <font image> <bin file> <font height>
```

To convert one fixed width font image to a `.bin` file, run:

```bash
python3 convert.py <font image> <bin file> <font height> <font width>
```

## Copyright

The fonts in these directories are released under the [Creative Commons
Attribution-ShareAlike
license](https://creativecommons.org/licenses/by-sa/4.0/).

This means that when you use these fonts in your own projects, you are required
to give credit to the author and that derivative works may only be distributed
under the same or a compatible license. This way we make sure the body of
publicly available pixel art fonts will only grow over time.

The 3x3 fixed width font called "Auri" was designed by Thumby Discord regular
Auri, with feedback from the Thumby community. You may give him credit as such.

The other fonts have been designed by myself
([Timendus](https://github.com/Timendus/)). Please credit me by linking to [this
font library](../), so others may find it and the fonts contained in it.
