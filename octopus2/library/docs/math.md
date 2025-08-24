# Math library

A set of macros to generate mathematical routines for addition, subtraction,
multiplication and division. Supports 8, 16, 24 and 32 bits (for some
operations for now; it's not 100% complete yet).

Each macro accepts a list of registers as parameters, which the routine is
allowed to use. As inputs, as outputs or as temporary registers. Usage
example using [mul-8-8](#mathmul-8-8):

```
:include "std/math"

# Multiplies v0 with v1, result in [v2, v3]. Destroys v0, v1, v4, vF.
: multiply
    math.mul-8-8 v0 v1 v2 v3 v4
    return
```

## Table of Contents

- Division
  - Macros
    - [math.div-16-8](#mathdiv-16-8)
    - [math.div-8-8](#mathdiv-8-8)
- Multiplication
  - Macros
    - [math.mul-8-8](#mathmul-8-8)
- Addition
  - Macros
    - [math.add-16-8](#mathadd-16-8)
    - [math.add-16-16](#mathadd-16-16)
    - [math.add-32-8](#mathadd-32-8)
    - [math.add-32-16](#mathadd-32-16)
    - [math.add-32-32](#mathadd-32-32)
    - [math.add-16-to-i](#mathadd-16-to-i)
- Subtraction
  - Macros
    - [math.sub-8-16](#mathsub-8-16)
    - [math.sub-8-32](#mathsub-8-32)
    - [math.sub-16-8](#mathsub-16-8)
    - [math.sub-16-16](#mathsub-16-16)
    - [math.sub-16-32](#mathsub-16-32)
    - [math.sub-32-8](#mathsub-32-8)
    - [math.sub-32-16](#mathsub-32-16)
    - [math.sub-32-32](#mathsub-32-32)


# Division

## Macros

### `math.div-16-8`

_library/lib/math.8o:29_

Divide 16-bit value by 8-bit value, resulting in 8-bit result

```
 numerator-low = [numerator-high, numerator-low] / divisor
 remainder = [numerator-high, numerator-low] % divisor
```

Destroys: numerator-high, temp1, temp2, vF

#### Parameters
- numerator-high
- numerator-low
- divisor
- remainder
- temp1
- temp2

### `math.div-8-8`

_library/lib/math.8o:67_

Divide 8-bit value by 8-bit value, resulting in 8-bit result

```
 numerator = numerator / divisor
 remainder = numerator % divisor
```

Destroys: temp1, temp2, vF

#### Parameters
- numerator
- divisor
- remainder
- temp1
- temp2


# Multiplication

## Macros

### `math.mul-8-8`

_library/lib/math.8o:103_

Multiply two 8-bit values, resulting in 16-bit result

```
 [result-high, result-low] = multiplicand * multiplier
```

Destroys: multiplicand, multiplier, temp, vF

#### Parameters
- multiplicand
- multiplier
- result-high
- result-low
- temp


# Addition

## Macros

### `math.add-16-8`

_library/lib/math.8o:181_

Add an 8-bit and a 16-bit value, resulting in 16-bit result + vF = 1 if overflow

```
 [value1-high, value1-low] += value2
```

Destroys: value1-high, value1-low, vF

#### Parameters
- value1-high
- value1-low
- value2

### `math.add-16-16`

_library/lib/math.8o:192_

Add two 16-bit values, resulting in 16-bit result + overflow = 1 if overflow

```
 [value1-high, value1-low] += [value2-high, value2-low]
```

Destroys: value1-high, value1-low, vF

#### Parameters
- value1-high
- value1-low
- value2-high
- value2-low
- overflow

### `math.add-32-8`

_library/lib/math.8o:206_

Add an 8-bit and a 32-bit value, resulting in 32-bit result + vF = 1 if overflow

```
 [value1-high, value1-mid1, value1-mid2, value1-low] += value2
```

Destroys: value1-high, value1-mid1, value1-mid2, value1-low, vF

#### Parameters
- value1-high
- value1-mid1
- value1-mid2
- value1-low
- value2

### `math.add-32-16`

_library/lib/math.8o:219_

Add a 16-bit and a 32-bit value, resulting in 32-bit result + overflow = 1 if overflow

```
 [value1-high, value1-mid1, value1-mid2, value1-low] += value2
```

Destroys: value1-high, value1-mid1, value1-mid2, value1-low, overflow, vF

#### Parameters
- value1-high
- value1-mid1
- value1-mid2
- value1-low
- value2-high
- value2-low
- overflow

### `math.add-32-32`

_library/lib/math.8o:230_

Add two 32-bit values, resulting in 32-bit result + overflow = 1 if overflow

```
 [value1-high, value1-mid1, value1-mid2, value1-low] += [value2-high, value2-mid1, value2-mid2, value2-low]
```

Destroys: value1-high, value1-mid1, value1-mid2, value1-low, overflow, vF

#### Parameters
- value1-high
- value1-mid1
- value1-mid2
- value1-low
- value2-high
- value2-mid1
- value2-mid2
- value2-low
- overflow

### `math.add-16-to-i`

_library/lib/math.8o:243_

Add 16-bit value to index register

```
 i += [value-high, value-low]
```

Destroys: value-high, vF

#### Parameters
- value-high
- value-low


# Subtraction

## Macros

### `math.sub-8-16`

_library/lib/math.8o:263_

Subtract a 16-bit value from an 8-bit value, resulting in 8-bit result + carry = 0 if carry

```
 value1 -= [value2-high, value2-low]
```

Destroys: value1, vF

#### Parameters
- value1
- value2-high
- value2-low
- carry

### `math.sub-8-32`

_library/lib/math.8o:275_

Subtract a 32-bit value from an 8-bit value, resulting in 8-bit result + carry = 0 if carry

```
 value1 -= [value2-high, value2-mid1, value2-mid2, value2-low]
```

Destroys: value1, vF

#### Parameters
- value1
- value2-high
- value2-mid1
- value2-mid2
- value2-low
- carry

### `math.sub-16-8`

_library/lib/math.8o:289_

Subtract an 8-bit value from a 16-bit value, resulting in 16-bit result + vF = 0 if carry

```
 [value1-high, value1-low] -= value2
```

Destroys: value1-high, value1-low, vF

#### Parameters
- value1-high
- value1-low
- value2

### `math.sub-16-16`

_library/lib/math.8o:303_

Subtract two 16-bit values, resulting in 16-bit result + carry = 0 if carry

```
 [value1-high, value1-low] -= [value2-high, value2-low]
```

Destroys: value1-high, value1-low, carry, vF

#### Parameters
- value1-high
- value1-low
- value2-high
- value2-low
- carry

### `math.sub-16-32`

_library/lib/math.8o:321_

Subtract a 32-bit value from a 16-bit value, resulting in 16-bit result + carry = 0 if carry

```
 [value1-high, value1-low] -= [value2-high, value2-mid1, value2-mid2, value2-low]
```

Destroys: value1-high, value1-low, carry, vF

#### Parameters
- value1-high
- value1-low
- value2-high
- value2-mid1
- value2-mid2
- value2-low
- carry

### `math.sub-32-8`

_library/lib/math.8o:341_

Subtract an 8-bit value from a 32-bit value, resulting in 32-bit result + vF = 0 if carry

```
 [value1-high, value1-mid1, value1-mid2, value1-low] -= value2
```

Destroys: value1-high, value1-mid1, value1-mid2, value1-low, vF

#### Parameters
- value1-high
- value1-mid1
- value1-mid2
- value1-low
- value2

### `math.sub-32-16`

_library/lib/math.8o:363_

Subtract a 16-bit value from a 32-bit value, resulting in 32-bit result + carry = 0 if carry

```
 [value1-high, value1-mid1, value1-mid2, value1-low] -= [value2-high, value2-low]
```

Destroys: value1-high, value1-mid1, value1-mid2, value1-low, carry, vF

#### Parameters
- value1-high
- value1-mid1
- value1-mid2
- value1-low
- value2-high
- value2-low
- carry

### `math.sub-32-32`

_library/lib/math.8o:397_

Subtract a 32-bit value from a 32-bit value, resulting in 32-bit result + carry = 0 if carry

```
 [value1-high, value1-mid1, value1-mid2, value1-low] -= [value2-high, value2-mid1, value2-mid2, value2-low]
```

Destroys: value1-high, value1-mid1, value1-mid2, value1-low, carry, vF

#### Parameters
- value1-high
- value1-mid1
- value1-mid2
- value1-low
- value2-high
- value2-mid1
- value2-mid2
- value2-low
- carry


