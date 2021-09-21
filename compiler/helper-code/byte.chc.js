module.exports = `

raw <<

store-byte:
  call get-frame-pointer
  add i, v10
  ld v0, v11
  ld (i), v0-v0
  ret

load-byte:
  call get-frame-pointer
  add i, v10
  ld v0-v0, (i)
  ld v11, v0
  ret

add-byte:
  call get-frame-pointer
  add i, v10
  ld v0-v0, (i)
  ld v3, v0
  call get-frame-pointer
  add i, v11
  ld v0-v0, (i)
  ld v11, v0
  add v11, v3
  ret

>>

`;
