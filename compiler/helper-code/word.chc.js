module.exports = `

raw <<

store-word:
  call get-frame-pointer
  add i, v10
  ld v0, v11
  ld v1, v12
  ld (i), v0-v1
  ret

load-word:
  call get-frame-pointer
  add i, v10
  ld v0-v1, (i)
  ld v11, v0
  ld v12, v1
  ret

add-word:
  call get-frame-pointer
  add i, v10
  ld v0-v1, (i)
  ld v3, v0
  ld v4, v1
  call get-frame-pointer
  add i, v11
  ld v0-v1, (i)
  ld v12, v4
  add v12, v1
  ld v11, v15
  add v11, v3
  add v11, v0
  ret

>>

`;
