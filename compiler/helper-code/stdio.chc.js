module.exports = `

safe function clear_screen() {
  raw cls
}

function print_byte(byte value, byte xpos, byte ypos) {
  raw <<

  ld i, stdio_buffer
  bcd (i), value
  ld v0-v0, (i)
  sne v0, 0
  jp print_byte_first_zero
  call print_byte_draw
  jp print_byte_first_not_zero

print_byte_first_zero:
  add xpos, 5
  ld i, stdio_buffer + 1
  ld v0-v0, (i)
  se v0, 0
  call print_byte_draw
  jp print_byte_third_digit

print_byte_first_not_zero:
  add xpos, 5
  ld i, stdio_buffer + 1
  ld v0-v0, (i)
  call print_byte_draw

print_byte_third_digit:
  add xpos, 5
  ld i, stdio_buffer + 2
  ld v0-v0, (i)
  call print_byte_draw
  ret

print_byte_draw:
  getfont v0
  drw xpos, ypos, 5
  ret

  >>
}

function modulo(byte x, byte y) {
  return x - (y * (x/y));
}

function point(byte x, byte y) {
  raw <<

  ld i, stdio_buffer
  ld v0, 128
  ld (i), v0-v0
  drw x, y, 1

  >>
}

// // Bresenham's line algorithm adapted for CHIPcode
// // https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
// function drawLine(byte x0, byte y0, byte x1, byte y1) {
//   byte dx =  abs(x1-x0);
//   byte sx = x0 < x1 ? 1 : -1;
//   byte dy = -abs(y1-y0);
//   byte sy = y0 < y1 ? 1 : -1;
//   byte err = dx + dy;  /* error value e_xy */
//   while ( x0 == x1 && y0 == y1 ) {
//     plot(x0, y0);
//     byte e2 = err + err; // 2 * err
//     if ( e2 + 1 > dy ) { // e2 >= dy
//         err = err + dy; /* e_xy+e_x > 0 */
//         x0 = x0 + sx;
//     }
//     if ( e2 < dx + 1 ) { // e2 <= dx
//         err = err + dx; /* e_xy+e_y < 0 */
//         y0 = y0 + sy;
//     }
//   }
//   plot(x0, y0);
// }

raw <<

stdio_buffer:
  .db 0,0,0,0,0

>>

`;
