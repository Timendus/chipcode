_A=None
from os import stat
from sys import path
VARIABLE_WIDTH=const(0)
NEWLINE=const(10)
SPACE=const(32)
class FancyFont:
	@micropython.native
	def __init__(self,displayBuffer,displayWidth=72,displayHeight=40):A=self;A.displayBuffer=displayBuffer;A.displayWidth=int(displayWidth);A.displayHeight=int(displayHeight)
	@micropython.native
	def setFont(self,fontPath,width:int=_A,height:int=_A,space:int=1):
		D=height;C=width;B=fontPath;A=self;B=A._findFile(B);A.fontFile=open(B,'rb');A.characterBuffer=bytearray(9)
		if C==_A and D==_A:A.characterWidth=VARIABLE_WIDTH;A.characterMarginWidth=0;A.fontFile.readinto(A.characterBuffer);A.characterHeight=A.characterBuffer[0];A.numCharactersInFont=A.characterBuffer[1];A._collectCharacterIndices()
		else:A.characterWidth=C;A.characterHeight=D;A.characterMarginWidth=space;A.numCharactersInFont=stat(B)[6]//A.characterWidth
	@micropython.native
	def drawText(self,string,xPos:int,yPos:int,color:int=1,xMax:int=_A,yMax:int=_A):B=string;A=self;return A._drawText(B,len(B),xPos,yPos,color,xMax or A.displayWidth,yMax or A.displayHeight)
	@micropython.native
	def drawTextWrapped(self,string,xPos:int,yPos:int,color:int=1,xMax:int=_A,yMax:int=_A):B=string;A=self;C=A._wrapText(B,len(B),xPos,yPos,xMax or A.displayWidth,yMax or A.displayHeight);return A._drawText(C,len(C),xPos,yPos,color,xMax or A.displayWidth,yMax or A.displayHeight)
	def _findFile(D,filePath):
		B=filePath
		try:stat(B);return B
		except OSError:pass
		for A in path:
			try:C=(A+'/'if A and not A.endswith('/')else A)+B;stat(C);return C
			except OSError:pass
		raise OSError('Font file not found')
	@micropython.viper
	def _collectCharacterIndices(self):
		A=self;A.characterIndices=[];B:int=2
		for C in range(int(A.numCharactersInFont)):A.characterIndices.append(B);A.fontFile.seek(B);A.fontFile.readinto(A.characterBuffer);B+=int(A.characterBuffer[0])+1
	@micropython.viper
	def _drawText(self,string:ptr8,strLen:int,xStart:int,yPos:int,color:int,xMax:int,yMax:int):
		T=strLen;S=string;O=yMax;N=xMax;M=xStart;C=yPos;B=self;S:ptr8=ptr8(S);T:int=int(T);M:int=int(M);C:int=int(C);N:int=int(N);O:int=int(O);U:int=0;P:int=0;F:int=0;G:int=0;J:int=0;D:int=0;H:int=0;I:int=0;e:int=0;E:int=M;V:int=E;W:int=C;Q:ptr8=ptr8(B.displayBuffer);X=B.fontFile;a=B.characterBuffer;Y:ptr8=ptr8(a);K:int=int(B.displayWidth);L:int=int(B.characterWidth);Z:int=int(B.characterHeight);b:int=int(B.characterMarginWidth);c:int=int(B.numCharactersInFont);R:int=int(len(B.displayBuffer))
		if L==VARIABLE_WIDTH:d=B.characterIndices
		while U<T:
			F=S[U];U+=1
			if F==NEWLINE:C+=Z+1;E=M;continue
			F=F-SPACE
			if not 0<=F<c:continue
			if L==VARIABLE_WIDTH:X.seek(int(d[F]))
			else:X.seek(F*L)
			X.readinto(a)
			if L==VARIABLE_WIDTH:G=Y[0];P=1
			else:G=L;P=0
			if E+G<=0 or E>=N or C+Z<=0:E+=G+b;continue
			if C>=O:break
			I=N-E
			if I>G:I=G
			W=int(max(W,E+I-1));V=int(max(V,min(C+Z-1,O-1)));J=C&7;D=(C>>3)*K+E
			if color==0:
				for A in range(0,I):
					H=Y[P+A]
					if 0<=D+A<R:Q[D+A]&=255^H<<J
					if 0<=D+A+K<R:Q[D+A+K]&=255^H>>8-J
			else:
				for A in range(0,I):
					H=Y[P+A]
					if 0<=D+A<R:Q[D+A]|=H<<J
					if 0<=D+A+K<R:Q[D+A+K]|=H>>8-J
			E+=G+b
		return bytearray([W,V])
	@micropython.viper
	def _wrapText(self,string:ptr8,strLen:int,xStart:int,yPos:int,xMax:int,yMax:int):
		L=yMax;K=xMax;H=string;F=yPos;E=xStart;D=strLen;C=self;H:ptr8=ptr8(H);D:int=int(D);E:int=int(E);F:int=int(F);K:int=int(K);L:int=int(L);A:int=0;B:int=0;I:int=0;J:int=E;M=C.fontFile;O=C.characterBuffer;R:ptr8=ptr8(O);G:int=int(C.characterWidth);P:int=int(C.characterHeight);S:int=int(C.characterMarginWidth);T:int=int(C.numCharactersInFont)
		if G==VARIABLE_WIDTH:U=C.characterIndices
		Q=bytearray(D);N:ptr8=ptr8(Q)
		while A<D:
			B=H[A];N[A]=B;A+=1
			if B==NEWLINE:F+=P+1;J=E;continue
			B=B-SPACE
			if not 0<=B<T:continue
			if G==VARIABLE_WIDTH:M.seek(U[B])
			else:M.seek(B*G)
			M.readinto(O);I=G
			if G==VARIABLE_WIDTH:I=R[0]
			if F>=L:break
			if J+I>=K:
				A-=1;V=A
				while A>0 and H[A]!=SPACE:A-=1
				if A==0 or N[A]==NEWLINE:A=V
				else:
					if 0>A>=D:raise ValueError('Attempted to write outside string bounds: this should never happen!')
					N[A]=NEWLINE
				A+=1;F+=P+1;J=E;continue
			J+=I+S
		return Q