# util.8o

## Table of Contents

  - Macros
    - [util.set-16](#util.set-16)
    - [util.set-32](#util.set-32)
    - [util.sleep](#util.sleep)


## Macros

### `util.set-16`

_library/lib/util.8o:5_

Set two registers to a single 16-bit value

```
 [high, low] = value
```

#### Parameters
- high
- low
- value

### `util.set-32`

_library/lib/util.8o:16_

Set four registers to a single 32-bit value

```
 [high, mid1, mid2, low] = value
```

#### Parameters
- high
- mid1
- mid2
- low
- value

### `util.sleep`

_library/lib/util.8o:31_

Wait for the delay timer to run out

Destroys: vF


