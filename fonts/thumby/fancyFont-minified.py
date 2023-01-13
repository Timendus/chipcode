_A=None
from thumbyGraphics import display
from os import stat
VARIABLE_WIDTH=const(0)
NEWLINE=const(10)
SPACE=const(32)
class FancyFont:
	@micropython.native
	def setFont(self,fontPath,width:int=_A,height:int=_A,space:int=1):
		D=height;C=width;B=fontPath;A=self;A.fontFile=open(B,'rb');A.characterBuffer=bytearray(9)
		if C==_A and D==_A:A.characterWidth=VARIABLE_WIDTH;A.characterMarginWidth=0;A.fontFile.readinto(A.characterBuffer);A.characterHeight=A.characterBuffer[0];A.numCharactersInFont=A.characterBuffer[1];A._collectCharacterIndices()
		else:A.characterWidth=C;A.characterHeight=D;A.characterMarginWidth=space;A.numCharactersInFont=stat(B)[6]//A.characterWidth
	@micropython.native
	def drawTextWrapped(self,string,xPos:int,yPos:int,color:int=1,xMax:int=display.width,yMax:int=display.height):A=string;B=self._wrapText(A,len(A),xPos,yPos,xMax,yMax);return self._drawText(B,len(B),xPos,yPos,color,xMax,yMax)
	@micropython.native
	def drawText(self,string,xPos:int,yPos:int,color:int=1,xMax:int=display.width,yMax:int=display.height):A=string;return self._drawText(A,len(A),xPos,yPos,color,xMax,yMax)
	@micropython.viper
	def _collectCharacterIndices(self):
		A=self;A.characterIndices=[];B:int=2
		for C in range(int(A.numCharactersInFont)):A.characterIndices.append(B);A.fontFile.seek(B);A.fontFile.readinto(A.characterBuffer);B+=int(A.characterBuffer[0])+1
	@micropython.viper
	def _drawText(self,string:ptr8,strLen:int,xStart:int,yPos:int,color:int,xMax:int,yMax:int):
		T=strLen;S=string;O=yMax;N=xMax;M=xStart;E=self;B=yPos;S:ptr8=ptr8(S);T:int=int(T);M:int=int(M);B:int=int(B);N:int=int(N);O:int=int(O);U:int=0;P:int=0;F:int=0;G:int=0;J:int=0;C:int=0;H:int=0;I:int=0;e:int=0;D:int=M;V:int=D;W:int=B;Q:ptr8=ptr8(display.display.buffer);X=E.fontFile;a=E.characterBuffer;Y:ptr8=ptr8(a);K:int=int(display.width);L:int=int(E.characterWidth);Z:int=int(E.characterHeight);b:int=int(E.characterMarginWidth);c:int=int(E.numCharactersInFont);R:int=int(len(display.display.buffer))
		if L==VARIABLE_WIDTH:d=E.characterIndices
		while U<T:
			F=S[U];U+=1
			if F==NEWLINE:B+=Z+1;D=M;continue
			F=F-SPACE
			if not 0<=F<c:continue
			if L==VARIABLE_WIDTH:X.seek(int(d[F]))
			else:X.seek(F*L)
			X.readinto(a)
			if L==VARIABLE_WIDTH:G=Y[0];P=1
			else:G=L;P=0
			if D+G<=0 or D>=N or B+Z<=0:D+=G+b;continue
			if B>=O:break
			I=N-D
			if I>G:I=G
			W=int(max(W,D+I-1));V=int(max(V,min(B+Z-1,O-1)));J=B&7;C=(B>>3)*K+D
			if color==0:
				for A in range(0,I):
					H=Y[P+A]
					if 0<=C+A<R:Q[C+A]&=255^H<<J
					if 0<=C+A+K<R:Q[C+A+K]&=255^H>>8-J
			else:
				for A in range(0,I):
					H=Y[P+A]
					if 0<=C+A<R:Q[C+A]|=H<<J
					if 0<=C+A+K<R:Q[C+A+K]|=H>>8-J
			D+=G+b
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
fancyFont=FancyFont()