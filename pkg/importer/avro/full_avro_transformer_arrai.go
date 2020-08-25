// Package avro ...

package avro

const AvroTransformerScript string = `### ------------------------------------------------------------------------ ###
###  avro_util.arrai                                                         ###
### ------------------------------------------------------------------------ ###

# Utilities
let indent = "    ";

let getAnnotationS = \item \name
    let annotation = item(name)? .s:"";
    let annotationTxt = 
        cond annotation {
            "":"", 
            _: $'${cond name {'doc':'description', _:name}}="${annotation}"'
        };
    annotationTxt
;

# Pretty annotation array
let prettyAnnotations = \annotationArray
    let str = $'${let a = annotationArray where .@item != "" rank (:.@); a >>> \i \item item::, }';
    cond str {'':$'', _:$'[${str}]'}
;

# Transform Avro primitive types to Sysl ones
let transformPrimitiveType = \type
    cond type {
        'null': 'null',
        'boolean': 'bool',
        'int': 'int32',
        'long': 'int64',
        'float': 'float32',
        'double': 'float64',
        'bytes': "bytes",
        'string': 'string',
        _: 'NonePrimitiveType'              
    }
;

let transformLogicalTypes = \type
    cond type('logicalType').s {
        'decimal': $'decimal(${type('precision')}.${type('scale')})',
        'uuid': 'string',
        'date': 'date',
        'time-millis': 'uint32',
        'time-micros': 'uint64',
        'timestamp-millis': 'datetime',
        'timestamp-micros': 'datetime',
        'local-timestamp-millis': 'datetime',
        'local-timestamp-micros': 'datetime',
        'duration': 'bytes',
        _: 'NoneLogicalType'  
    }
;

# Get type name from types which looks like:
# ['null', 'string'], ['string', 'int'] or ['null', {'type':'record', ...}] etc.
let getTypeName = \t
    cond t {
        (:s): t,
        {'type': (s: type), 'items': (:s), ...}:
            cond type {
                'map': (s: $'${//str.title(t('items').s)}Set', type_S: $'set of ${t('items').s}'), # will be displayed in union definition
                'array': (s: $'${//str.title(t('items').s)}Sequence', type_S: $'sequence of ${t('items').s}'), # will be displayed in union definition
            },
        {'type': (s: type), 'items': {'type': (:s), ...}, ...}:
            cond type {
                'map': (s: $'${t('items')('name').s}Set', type_S: $'set of ${t('items')('name').s}'), # will be displayed in union definition
                'array': (s: $'${t('items')('name').s}Sequence', type_S: $'sequence of ${t('items')('name').s}'), # will be displayed in union definition
            },
        {'type': (:s), ...}: (s : t('name').s),
    }
;

let getDisplayedTypeName = \.
    cond . {(type_S:type_S, ...): type_S, _: .s}
;

let combineTypes = \types
    let fullTypeNames = types >> getTypeName(.);
    let types = fullTypeNames >> .s;
    cond //seq.contains(['null'])(types) {
        true:
            let typeNames = fullTypeNames where .@item.s != 'null' rank (:.@);
            cond types count {
                2: //seq.concat(typeNames >> getDisplayedTypeName(.)) + '?',
                _: //seq.concat(typeNames >> //str.title(.s)) + '?'
            }
        ,
        _: //seq.concat(types >> //str.title(.))
    }
;

let transformType = \type
    cond type {
        # (s: 'string') etc.
        (s: typeName):
            let primitive = transformPrimitiveType(typeName);
            cond primitive {
                "null": "string[~null]",
                'NonePrimitiveType': typeName,
                _: primitive
            },
        (a: typeArray):
            # type can be ['null', 'string'], ['string', 'int'] or ['null', {'type':'record', ...}] etc.
            # [(s: 'null'), (s: 'string')] to ['null', 'string']
            combineTypes(typeArray),
        {'logicalType': logicalTypeName, 'type': typeName, ...}: # Must be before {'type': typeName, ...}
            transformLogicalTypes(type), 
        {'type': typeName, ...}:
            # it is array, map
            cond typeName.s {
                'array': 'sequence of ' + cond type {{'items': {'name': name, ...}, ...}: name.s, _: type('items').s},
                'map': 'set of ' + cond type {{'items': {'name': name, ...}, ...}: 'String' + name.s + 'Item', _: type('items').s},
                _: type('name').s
            },
    }
;

# TODO default = {} can't work still
let getDefaultVal = \.
    let default = .('default')?:""; 
    cond default {
        "": "",
        ():'default="null"', 
        (b: true): $'default="true"',
        (b: false): $'default="false"',
        (a: {}): $'default="[]"',
        (a: [(:s), ...]): $'default="${default.a >> .s}"',
        _: $'default="${default}"'
    } 
;

let util = (
    : prettyAnnotations,
    : transformType,
    : getAnnotationS,
    : getTypeName,
    : combineTypes,
    : indent,
    : getDefaultVal
);let avro_util_arrai = 
util
;

### ------------------------------------------------------------------------ ###
###  to_sysl_alias.arrai                                                     ###
### ------------------------------------------------------------------------ ###

let util = avro_util_arrai;

let t = \.
cond . {
{'type': 'fixed', ...}:
    $'${util.indent}!alias ${.('name').s}${
        util.prettyAnnotations(['~fixed', util.getAnnotationS(., 'namespace')])
    }: bytes(${.('size')})
    ${
        let aliases = .('aliases')? .a:[];
        aliases >> \alias ($'!alias ${alias.s}: ${.('name').s}')::
    }',
{'type': (s: 'array'), ...}:
    $'${util.indent}!alias ${//str.title(.('items').s)}Sequence${
        util.prettyAnnotations([util.getDefaultVal(.)])
    }: sequence of ${.('items').s}'
};let to_sysl_alias_arrai = t
;

### ------------------------------------------------------------------------ ###
###  to_sysl_enum.arrai                                                      ###
### ------------------------------------------------------------------------ ###

# Transform Avro enum to Sysl enum
let util = avro_util_arrai;

let t = \enum

let annotations = [util.getAnnotationS(enum, 'namespace'), util.getAnnotationS(enum, 'default'), util.getAnnotationS(enum, 'doc')];

$'${util.indent}!enum ${enum('name').s}${util.prettyAnnotations(annotations)}:
${
    enum('symbols').a >>> \i \item  $'${util.indent}${util.indent}${item.s}: ${i}'::
}
${
let aliases = enum('aliases')? .a:[]; 
aliases >> ($'

${util.indent}!alias ${.s}:
${util.indent}${util.indent}${enum('name').s}')::\n}
';let to_sysl_enum_arrai = t
;

### ------------------------------------------------------------------------ ###
###  to_sysl_type.arrai                                                      ###
### ------------------------------------------------------------------------ ###

# Transform Avro record to Sysl record
let util = avro_util_arrai;

# Build type annotation like [union=["int32", "null"]]
let buildTypeAnnotation = \type
    cond type {
        (a: typeArray):
            let types = (typeArray >> util.getTypeName(.)) >> .s;
            cond //seq.contains(['null'])(types) {
                true:
                    cond //seq.has_prefix(['null'])(types) {
                        true: # types is ['null', 'string']
                            '',
                        _: # types is ['string', 'null']
                            $'union=${types}'
                    }
                ,
                _: ''
            }
        ,
        _: ''
    }
;

let buildLogicalTypeAnnotation = \.
    cond .('type') {
        {'type': typeName, 'logicalType': logicalTypeName, ...}:
            cond logicalTypeName.s {
                'uuid': '[~uuid]',
                'date': '',
                'time-millis': ['~time-millis'],
                'time-micros': ['~time-micros'],
                'timestamp-millis': ['~timestamp-millis'],
                'timestamp-micros': ['~timestamp-micros'],
                'local-timestamp-millis': ['~local-timestamp-millis'],
                'local-timestamp-micros': ['~local-timestamp-micros'],
                'duration': ['~fixed', 'fixed_size="12"', '~duration'],
                _: []
            },
        _: []
    }
;

let buildAnnotations = \.
    buildLogicalTypeAnnotation(.) + 
    [
        buildTypeAnnotation(.('type')),
        util.getAnnotationS(., 'doc'),
        util.getDefaultVal(.),
        util.getAnnotationS(., 'order'),
    ]
;

let t = \record
$'${
    cond record {
        {'isMap':(s: 'true'), ...}:
            let valueType = cond record {
                {"name": (s: name)}: name,
                {"items": (s: type), ...}: type
            };
            $'${'\n'}${util.indent}!type String${valueType}Item${
                    util.prettyAnnotations(['json_map_key="key"', util.getDefaultVal(record)])
                }:${'\n'}${util.indent}${util.indent}key <: string${
                '\n'}${util.indent}${util.indent}value <: ${valueType}${'\n\n'}',
    }
}${
cond record {
{'type': (s: 'record'), ...}:
$'${util.indent}!type ${record('name').s}${
    util.prettyAnnotations([util.getAnnotationS(record, 'namespace'), util.getAnnotationS(record, 'doc')])
}:
${
    let fields = record('fields')? .a:[]; 
    fields >> 
        ($'${util.indent}${util.indent}${.('name').s} <: ${util.transformType(.('type'))} ${
                util.prettyAnnotations(buildAnnotations(.))
            }${
                let aliases = .('aliases')? .a:[];
                cond aliases {
                    []: '',
                    _: $'${'\n'}${aliases >> \alias 
                                    $'${util.indent}${util.indent}${alias.s} <: ${util.transformType(.('type'))} ${
                                        util.prettyAnnotations(buildAnnotations(.) + [$'alias_of="${.('name').s}"'])}'::\i}'
                }
            }')::
}

${
let aliases = record('aliases')? .a:[];
aliases >> 
($'${util.indent}!alias ${.s}:
${util.indent}${util.indent}${record('name').s}')::\i
}'}}';let to_sysl_type_arrai = t
;

### ------------------------------------------------------------------------ ###
###  to_sysl_union.arrai                                                     ###
### ------------------------------------------------------------------------ ###

# Transform Avro union to Sysl union
let util = avro_util_arrai;

let combineTypes = \types
    let fullTypeNames = types >> util.getTypeName(.);
    let types = fullTypeNames >> .s;
    cond //seq.contains(['null'])(types) {
        true:
            let typeNames = fullTypeNames where .@item.s != 'null' rank (:.@);
            cond types count {
                2: //seq.concat(typeNames >> .s) + '?',
                _: //seq.concat(typeNames >> //str.title(.s)) + '?'
            }
        ,
        _: //seq.concat(types >> //str.title(.))
    }
;

let printUnion = \union
    let types = union >> util.getTypeName(.) >> cond . {(type: typeS, ...): typeS, _: .s};
    $'${util.indent}!union ${//seq.sub("?", "")(combineTypes(union))}:
    ${types where .@item != 'null' rank (:.@) >> $'${util.indent}${.}'::\n}${'\n'}'
;

let t = \union
    let types = union >> util.getTypeName(.) >> .s;
    cond //seq.contains(['null'])(types) {
        true:
            cond types count > 2 {
                true: printUnion(union),
            }
        ,
        _: printUnion(union)
    }
;let to_sysl_union_arrai = t
;

### ------------------------------------------------------------------------ ###
###  transformer.arrai                                                       ###
### ------------------------------------------------------------------------ ###

let util = avro_util_arrai;

# Extra items transformed to Sysl type !enum
let rec extraEnums = \schema
  cond schema {
    {'type': (s: 'enum'), 'name': name, ...}:
      {schema},
    {'type': (s: 'record'), 'name': name, "fields": (a : fields), ...}:
      //rel.union(fields => extraEnums(.@item("type"))),
    {'type': type, 'items': items, ...}: # Map or array
      cond items {
        (a: types): //rel.union(types => extraEnums(.@item)), # array
        (:s): extraEnums(s), # primitive
        _: extraEnums(items),
      },
    {'type': (a: types)}: # type is an array
      //rel.union(types => extraEnums(.@item)),  
    (a: types): # root is json array 
      //rel.union(types => extraEnums(.@item)),
  }
;

# Extra items transformed to Sysl type !type 
let rec extraTypes = \schema
  cond schema {
    {"type": (s: "record"), "name": (s: name), "fields": (a : fields), ...}:
      {schema} | //rel.union(fields => extraTypes(.@item("type"))),
    {'type': (s: type), 'items': items, ...}: # array or map
      cond items {
        (a: types): //rel.union(types => extraTypes(.@item)), # array
        (:s): cond type {'map': {schema | {'isMap': (s: 'true')}}}, # primitive
        _:
          cond type {
            'map': {schema('items') | {'isMap': (s: 'true')}} | extraTypes(items),
            _: extraTypes(items),
          } 
      },
    {'type': (a: types)}: # type is an array
      //rel.union(types => extraTypes(.@item)),  
    (a: types): # root is json array 
      //rel.union(types => extraTypes(.@item)),
  }
;

# Extra items transformed to Sysl type !alias
let rec extraAliases = \schema
  cond schema {
    {'type': (s: 'fixed'), ...}: # fixed type
      {schema},
    {"type": (s: "record"), "name": (s: name), "fields": (a : fields), ...}:
      //rel.union(fields => extraAliases(.@item("type"))),
    {'type': type, 'items': items, ...}: # array or map
      cond items {
        (a: types): //rel.union(types => extraAliases(.@item)), # array
        (:s): cond type.s {'array': {schema}}, # primitive
        _: extraAliases(items),
      },
    {'type': (a: types)}: # type is an array
      //rel.union(types => extraAliases(.@item)),  
    (a: types): # root is json array 
      //rel.union(types => extraAliases(.@item)),
};

# Extra item transformed to Sysl type !union
let rec extraUnions = \schema
  cond schema {
    (a: types): {types} | //rel.union(types => extraUnions(.@item)),
    {'type': (s: 'record'), "fields": (a : fields), ...}:
      //rel.union(fields => extraUnions(.@item("type"))),
    {'type': type, 'items': items, ...}: # array or map
      cond items {
        (a: types): //rel.union(types => extraUnions(.@item)), # array
        (:s): extraUnions(s), # primitive
        _: extraUnions(items),
      },
    {'type': (a: types)}: # type is an array
      //rel.union(types => extraUnions(.@item)),  
    (a: types): # root is json array 
      //rel.union(types => extraUnions(.@item)),    
  }
;

let avroTransform = \schema \appName \packageName
  # Load Avro spec whose format is protobuf 
  let schema = //encoding.json.decode(schema);

$'
##########################################
##                                      ##
##  AUTOGENERATED CODE -- DO NOT EDIT!  ##
##                                      ##
##########################################
${appName}${util.prettyAnnotations(['avro_spec="1.0"',$'${cond packageName {'':'', _:$'package="${packageName}"'}}'])}:
${util.indent}# Types

${extraTypes(schema) => to_sysl_type_arrai(.) orderby . ::}
${util.indent}# Aliases

${extraAliases(schema) => to_sysl_alias_arrai(.) orderby . ::\n}
${util.indent}# Unions

${
# It applies to_sysl_union(.) to every item and will remove duplidated items if to_sysl_union(.)
# produces, as it iterates set.
extraUnions(schema) => to_sysl_union_arrai(.) orderby . ::\n}
${util.indent}# Enums

${extraEnums(schema) => to_sysl_enum_arrai(.) orderby . ::\n}
';
avroTransform
`
