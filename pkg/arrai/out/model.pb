�
A�

A*�
Fetch�
Fetch:

SourceRead:�B�
�ok <: A.FetchResponse [dataflow={"A.FetchResponse.ax": "Source.Foo.x", "A.FetchResponse.ay": "Source.Foo.y"}, description="1:1 transform"]�

model.sysl2x
FetchResponseg�

model.sysl!J
#
ax�

model.sysl
#
ay�

model.sysl�

model.sysl�
B�

B*�
Fetch�
Fetch:

SourceRead:�B�
�ok <: B.FetchResponse [dataflow={"B.FetchResponse.bx": "Source.Foo.x", "B.FetchResponse.by": "Source.Foo.y"}, description="1:1 transform"]�

model.sysl"&2x
FetchResponseg�

model.sysl&*J
#
by�

model.sysl((
#
bx�

model.sysl''�

model.sysl!!�
C�

C*�
Fetch�
Fetch:

CFetchA:

CFetchB:�B�
�ok <: C.FetchResponse [dataflow={"C.FetchResponse.cx": ["A.FetchResponse.ax", "B.FetchResponse.bx"], "C.FetchResponse.cy": ["A.FetchResponse.ay", "B.FetchResponse.by"]}, description="1:1 transform"]�

model.sysl38*�
FetchA�
FetchA:

AFetch:�B�
�ok <: C.FetchResponse [dataflow={"C.FetchResponse.cx": "A.FetchResponse.ax", "C.FetchResponse.cy": "A.FetchResponse.ay"}, description="1:1 transform"]�

model.sysl+/*�
FetchB�
FetchB:

BFetch:�B�
�ok <: C.FetchResponse [dataflow={"C.FetchResponse.cx": "B.FetchResponse.bx", "C.FetchResponse.cy": "B.FetchResponse.by"}, description="1:1 transform"]�

model.sysl/32x
FetchResponseg�

model.sysl8<J
#
cx�

model.sysl99
#
cy�

model.sysl::�

model.sysl**�
D�

D*�
Fetch�
Fetch:

AFetch:

CFetch:�B�
�ok <: D.FetchResponse [dataflow={"D.FetchResponse.dx": ["A.FetchResponse.ax", "C.FetchResponse.cx"], "D.FetchResponse.dy": ["A.FetchResponse.ay", "C.FetchResponse.cy"]}, description="1:1 transform"]�

model.sysl=B2x
FetchResponseg�

model.syslBFJ
#
dx�

model.syslCC
#
dy�

model.syslDD�

model.sysl<<�
Client�

Client*�
Do�
Do:

DFetch:�B�
�ok <: Client.Screen [dataflow={"Client.Screen.xx": "D.FetchResponse.dx", "Client.Screen.yy": "D.FetchResponse.dy"}, description="1:1 transform"]�

model.syslGK2q
Screeng�

model.syslKOJ
#
xx�

model.syslLL
#
yy�

model.syslMM�

model.syslFF�
all{

all*Y
allR
all:

Source:

Client:
A:
B:
C:
D�

model.syslPW�

model.syslOO�
Source�

Source"
patterns:
"db"+
description"A database.
 Stores data.
*?
Read7
Read:B
ok <: Source.Foo�

model.sysl2�
Foo�B0
description!"A Foo.
 Represents foo things.
�

model.sysl:�
A
x<B
description"The x value.�

model.sysl

T
yOB0
description!"A Foo.
 Represents foo things.
�

model.sysl2�
Bar�B
description"A bar table.�

model.sysl:�
W
aRB
patterns:
"pkB
description"A bar table.�

model.sysl
F
bAB 
description"An optional int`�

model.sysl
]
xXB
description"A foreign key�

model.syslJ


SourceBarFoox
a�

model.sysl
