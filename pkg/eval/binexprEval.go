package eval

import (
	"fmt"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/pkg/errors"
)

type evalValueFunc func(*sysl.Value, *sysl.Value) *sysl.Value
type evalExprFunc func(
	ee *exprEval,
	assign Scope,
	list *sysl.Value,
	scopeVar string,
	rhs *sysl.Expr,
) *sysl.Value

// Strategy interface to evaluate binary expr
type Strategy interface {
	eval(*exprEval, Scope, *sysl.Expr_BinExpr) *sysl.Value
}

// DefaultBinExprStrategy is to evaluate lhs and rhs expr's first and then pass it to fn
type DefaultBinExprStrategy struct{}

// NegateBinExprStrategy is to evaluate the negative of the DefaultBinExprStrategy
type NegateBinExprStrategy struct{}

// LHSOverRHSStrategy binds rhs expression over each element of the LHS.
// Assumes lhs is a collection
type LHSOverRHSStrategy struct{}

//nolint:gochecknoglobals
var (
	functionEvalStrategy = map[sysl.Expr_BinExpr_Op]Strategy{
		sysl.Expr_BinExpr_EQ:      DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_ADD:     DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_SUB:     DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_MUL:     DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_MOD:     DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_DIV:     DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_IN:      DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_NOT_IN:  DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_BITOR:   DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_GT:      DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_LT:      DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_GE:      DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_LE:      DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_NE:      NegateBinExprStrategy{},
		sysl.Expr_BinExpr_AND:     DefaultBinExprStrategy{},
		sysl.Expr_BinExpr_FLATTEN: LHSOverRHSStrategy{},
		sysl.Expr_BinExpr_WHERE:   LHSOverRHSStrategy{},
	}
)

//nolint:gochecknoglobals
var (
	// key = op, lhs & rhs types
	valueFunctions = map[string]evalValueFunc{
		makeKey(sysl.Expr_BinExpr_ADD, ValueInt, ValueInt):        addInt64,
		makeKey(sysl.Expr_BinExpr_ADD, ValueString, ValueString):  addString,
		makeKey(sysl.Expr_BinExpr_AND, ValueBool, ValueBool):      andBool,
		makeKey(sysl.Expr_BinExpr_BITOR, ValueList, ValueList):    concatListList,
		makeKey(sysl.Expr_BinExpr_BITOR, ValueSet, ValueSet):      setUnion,
		makeKey(sysl.Expr_BinExpr_BITOR, ValueList, ValueSet):     concatListSet,
		makeKey(sysl.Expr_BinExpr_DIV, ValueInt, ValueInt):        divInt64,
		makeKey(sysl.Expr_BinExpr_EQ, ValueBool, ValueBool):       cmpBool,
		makeKey(sysl.Expr_BinExpr_EQ, ValueInt, ValueInt):         cmpInt,
		makeKey(sysl.Expr_BinExpr_EQ, ValueInt, ValueNull):        cmpNullFalse,
		makeKey(sysl.Expr_BinExpr_EQ, ValueNull, ValueNull):       cmpNullTrue,
		makeKey(sysl.Expr_BinExpr_EQ, ValueNull, ValueInt):        cmpNullFalse,
		makeKey(sysl.Expr_BinExpr_EQ, ValueNull, ValueString):     cmpNullFalse,
		makeKey(sysl.Expr_BinExpr_EQ, ValueString, ValueNull):     cmpNullFalse,
		makeKey(sysl.Expr_BinExpr_EQ, ValueString, ValueString):   cmpString,
		makeKey(sysl.Expr_BinExpr_EQ, ValueList, ValueNull):       cmpListNull,
		makeKey(sysl.Expr_BinExpr_GE, ValueInt, ValueInt):         geInt64,
		makeKey(sysl.Expr_BinExpr_GT, ValueInt, ValueInt):         gtInt64,
		makeKey(sysl.Expr_BinExpr_IN, ValueString, ValueList):     stringInList,
		makeKey(sysl.Expr_BinExpr_IN, ValueString, ValueNull):     stringInNull,
		makeKey(sysl.Expr_BinExpr_IN, ValueString, ValueSet):      stringInSet,
		makeKey(sysl.Expr_BinExpr_IN, ValueString, ValueMap):      stringInMapKey,
		makeKey(sysl.Expr_BinExpr_NOT_IN, ValueString, ValueList): stringNotInList,
		makeKey(sysl.Expr_BinExpr_NOT_IN, ValueString, ValueNull): stringNotInNull,
		makeKey(sysl.Expr_BinExpr_NOT_IN, ValueString, ValueSet):  stringNotInSet,
		makeKey(sysl.Expr_BinExpr_NOT_IN, ValueString, ValueMap):  stringNotInMapKey,
		makeKey(sysl.Expr_BinExpr_LE, ValueInt, ValueInt):         leInt64,
		makeKey(sysl.Expr_BinExpr_LT, ValueInt, ValueInt):         ltInt64,
		makeKey(sysl.Expr_BinExpr_MOD, ValueInt, ValueInt):        modInt64,
		makeKey(sysl.Expr_BinExpr_MUL, ValueInt, ValueInt):        mulInt64,
		makeKey(sysl.Expr_BinExpr_SUB, ValueInt, ValueInt):        subInt64,
	}

	// key = op, outer container_type, inner container type
	exprFunctions = map[string]evalExprFunc{
		makeKey(sysl.Expr_BinExpr_FLATTEN, ValueList, ValueNoArg): flattenListList, // empty list
		makeKey(sysl.Expr_BinExpr_FLATTEN, ValueList, ValueList):  flattenListList,
		makeKey(sysl.Expr_BinExpr_FLATTEN, ValueList, ValueSet):   flattenListSet,
		makeKey(sysl.Expr_BinExpr_FLATTEN, ValueSet, ValueList):   flattenSetList,
		makeKey(sysl.Expr_BinExpr_FLATTEN, ValueList, ValueMap):   flattenListMap,
		makeKey(sysl.Expr_BinExpr_FLATTEN, ValueSet, ValueMap):    flattenSetMap,
		makeKey(sysl.Expr_BinExpr_FLATTEN, ValueSet, ValueNoArg):  flattenSetSet, // empty list
		makeKey(sysl.Expr_BinExpr_FLATTEN, ValueSet, ValueSet):    flattenSetSet,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueList, ValueNoArg):   whereList,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueList, ValueList):    whereList,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueList, ValueString):  whereList,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueList, ValueMap):     whereList,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueMap, ValueNoArg):    whereMap,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueSet, ValueNoArg):    whereSet,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueSet, ValueMap):      whereSet,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueSet, ValueInt):      whereSet,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueSet, ValueFloat):    whereSet,
		makeKey(sysl.Expr_BinExpr_WHERE, ValueSet, ValueString):   whereSet,
	}
)

func makeKey(op sysl.Expr_BinExpr_Op, lhs, rhs valueType) string {
	return fmt.Sprintf("%s_%s_%s", op, lhs.String(), rhs.String())
}

func (op DefaultBinExprStrategy) eval(ee *exprEval, assign Scope, binexpr *sysl.Expr_BinExpr) *sysl.Value {
	lhsValue := Eval(ee, assign, binexpr.Lhs)
	rhsValue := Eval(ee, assign, binexpr.Rhs)
	key := makeKey(binexpr.Op, getValueType(lhsValue), getValueType(rhsValue))

	if f, has := valueFunctions[key]; has {
		return f(lhsValue, rhsValue)
	}

	panic(errors.Errorf("Unsupported operation:DefaultBinExprStrategy: %s", key))
}

func (op NegateBinExprStrategy) eval(ee *exprEval, assign Scope, binexpr *sysl.Expr_BinExpr) *sysl.Value {
	if binexpr.Op != sysl.Expr_BinExpr_NE {
		panic(errors.Errorf("Attempting to call Negate Strategy on a non == operator"))
	}

	negated := *binexpr
	negated.Op = sysl.Expr_BinExpr_EQ
	return unaryNeg(DefaultBinExprStrategy{}.eval(ee, assign, &negated))
}

func (op LHSOverRHSStrategy) eval(ee *exprEval, assign Scope, binexpr *sysl.Expr_BinExpr) *sysl.Value {
	lhsValue := Eval(ee, assign, binexpr.Lhs)
	vType := getValueType(lhsValue)
	itemType := getContainedType(lhsValue)
	key := makeKey(binexpr.Op, vType, itemType)
	scopeVarValue, hasScopeVar := assign[binexpr.Scopevar]
	if f, has := exprFunctions[key]; has {
		result := f(ee, assign, lhsValue, binexpr.Scopevar, binexpr.Rhs)
		delete(assign, binexpr.Scopevar)
		if hasScopeVar {
			assign[binexpr.Scopevar] = scopeVarValue
		}
		return result
	}
	panic(errors.Errorf("Unsupported operation: Cannot Execute %s", key))
}

func evalBinExpr(ee *exprEval, assign Scope, binexpr *sysl.Expr_BinExpr) *sysl.Value {
	if strategy, has := functionEvalStrategy[binexpr.Op]; has {
		return strategy.eval(ee, assign, binexpr)
	}
	panic(errors.Errorf("Unsupported operation: %s", binexpr.Op))
}
