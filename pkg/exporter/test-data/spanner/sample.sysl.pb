�
App�

App"
patterns:
"db2�
Types�B?
indexes4:2
0:.
	"name:Ix
!"key_parts:float64_,int64_(desc)�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl:�
q

string_100cB
patterns	:
"hexRd�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl'
k

string_max]B
patterns	:
"max�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl"
�
array�B
patterns:
"pk�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl%zF�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl%
�
	array_set��A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysljF�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl
`
bytes_1UB
length"1�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl$
O
optH`�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl
P
int64_F�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl
|
timestamp_commithB 
allow_commit_timestamp"true�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.syslE

O
date_F�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl	
T

timestamp_F�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl

X
string_1LR�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl
�
bytes_1024_hexoB
length"1024B
patterns	:
"hex�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl4
R
float64_F�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl
Q
string_F�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl
�
cust_id|B
patterns:
"fk�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl#J 


AppTypesCustomerid
P
bytes_F�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl
j
	bytes_max]B
patterns	:
"max�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl		 
O
bool_F�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl
array2�
Customer��A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl:�
b
id\B
patterns:
"pk�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl
p
ref_descdB
patterns:
"pk
"desc�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl#
n
ref_asccB
patterns:
"pk
"asc�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl!
id
ref_desc
ref_asc2�
Constraints��A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl.:�
�
cust_id�B
patterns:
"fk�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl&&#J&


AppConstraintsCustomerid
�
	types_int�B
patterns:
"fk�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl''&J'


AppConstraintsTypesint64_
�
types_string�B
patterns:
"fk�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl((*J(


AppConstraintsTypesstring_2�
NotTable��?
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl.0N
L
idF�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl//�A
3dev/sysl/pkg/exporter/test-data/spanner/sample.sysl