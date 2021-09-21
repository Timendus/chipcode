module.exports = `

raw <<

get-frame-pointer:
  ld i, frame-pointer
  ld v0-v1, (i)
  ld v2, $A0
  or v0, v2
  ld i, get-frame-pointer-smc
  ld (i), v0-v1
get-frame-pointer-smc:
  .db 0, 0
  ret

push-stack-frame:
  ld i, frame-pointer
  ld v0-v1, (i)
  ld v3, v0
  ld v4, v1
  call get-frame-pointer
  add i, v10
  ld v0, v3
  ld v1, v4
  ld (i), v0-v1
  ld i, frame-pointer
  add v10, 2
  add v1, v10
  add v0, v15
  ld (i), v0-v1
  ret

pop-stack-frame:
  ld i, frame-pointer
  ld v0-v1, (i)
  ld v15, 2
  sub v1, v15
  sub v0, v15
  ld v2, $A0
  or v0, v2
  ld i, pop-stack-frame-smc
  ld (i), v0-v1
pop-stack-frame-smc:
  .db 0, 0
  ld v0-v1, (i)
  ld i, frame-pointer
  ld (i), v0-v1
  ret

frame-pointer:
  .ref stack

stack:

>>

`;
