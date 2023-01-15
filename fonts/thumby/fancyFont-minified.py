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
		V=strLen;U=string;O=xMax;N=xStart;H=yMax;C=self;A=yPos;U:ptr8=ptr8(U);V:int=int(V);N:int=int(N);A:int=int(A);O:int=int(O);H:int=int(H);W:int=0;P:int=0;F:int=0;G:int=0;K:int=0;D:int=0;I:int=0;J:int=0;Q:int=0;E:int=N;X:int=E;Y:int=A;R:ptr8=ptr8(C.displayBuffer);Z=C.fontFile;b=C.characterBuffer;a:ptr8=ptr8(b);L:int=int(C.displayWidth);M:int=int(C.characterWidth);S:int=int(C.characterHeight);c:int=int(C.characterMarginWidth);d:int=int(C.numCharactersInFont);T:int=int(len(C.displayBuffer))
		if M==VARIABLE_WIDTH:e=C.characterIndices
		while W<V:
			F=U[W];W+=1
			if F==NEWLINE:A+=S+1;E=N;continue
			F=F-SPACE
			if not 0<=F<d:continue
			if M==VARIABLE_WIDTH:Z.seek(int(e[F]))
			else:Z.seek(F*M)
			Z.readinto(b)
			if M==VARIABLE_WIDTH:G=a[0];P=1
			else:G=M;P=0
			if E+G<=0 or E>=O or A+S<=0:E+=G+c;continue
			if A>=H:break
			J=O-E
			if J>G:J=G
			Q=255
			if A+S>H:Q>>=8-(H-A)
			Y=int(max(Y,E+J-1));X=int(max(X,min(A+S-1,H-1)));K=A&7;D=(A>>3)*L+E
			if color==0:
				for B in range(0,J):
					I=a[P+B]&Q
					if 0<=D+B<T:R[D+B]&=255^I<<K
					if 0<=D+B+L<T:R[D+B+L]&=255^I>>8-K
			else:
				for B in range(0,J):
					I=a[P+B]&Q
					if 0<=D+B<T:R[D+B]|=I<<K
					if 0<=D+B+L<T:R[D+B+L]|=I>>8-K
			E+=G+c
		return bytearray([Y,X])
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