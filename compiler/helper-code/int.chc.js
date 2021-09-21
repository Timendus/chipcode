module.exports = `

raw <<

store-int:
  call get-frame-pointer
  add i, v10
  ld v0, v11
  ld v1, v12
  ld v2, v13
  ld v3, v14
  ld (i), v0-v3
  ret

load-int:
  call get-frame-pointer
  add i, v10
  ld v0-v3, (i)
  ld v11, v0
  ld v12, v1
  ld v13, v2
  ld v14, v3
  ret

add-int:
  call get-frame-pointer
  add i, v10
  ld v0-v3, (i)
  ld v4, v0
  ld v5, v1
  ld v6, v2
  ld v7, v3
  call get-frame-pointer
  add i, v11
  ld v0-v3, (i)
  ld v14, v7
  add v14, v3
  ld v13, v15
  add v13, v6
  add v13, v2
  ld v12, v15
  add v12, v5
  add v12, v1
  ld v11, v15
  add v11, v4
  add v11, v0
  ret

>>

`;
