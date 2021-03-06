package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bunniesandbeatings/baduk/architecture"
	"github.com/bunniesandbeatings/baduk/contexts"
	"github.com/bunniesandbeatings/baduk/parser"
)

func TestParseVoidMethod(t *testing.T) {
	method := findMethod(t, "Void", parseTypesMethods(t))
	checkFuncShape(t, 0, 0, method)
}

func TestParseScalarMethod(t *testing.T) {
	method := findMethod(t, "Scalar", parseTypesMethods(t))
	checkFuncShape(t, 2, 2, method)
	checkParmTypes(t, method, "int", "string")
	checkReturnTypes(t, method, "int64", "uint32")
}

func TestParseSliceMethod(t *testing.T) {
	method := findMethod(t, "Slice", parseTypesMethods(t))
	checkFuncShape(t, 1, 1, method)
	checkParmTypes(t, method, "[]int")
	checkReturnTypes(t, method, "[]string")
}

func TestParseArrayMethod(t *testing.T) {
	method := findMethod(t, "Array", parseTypesMethods(t))
	checkFuncShape(t, 1, 1, method)
	checkParmTypes(t, method, "[2]int")
	checkReturnTypes(t, method, "[3]string")
}

func TestParsePointerMethod(t *testing.T) {
	method := findMethod(t, "Pointer", parseTypesMethods(t))
	checkFuncShape(t, 1, 1, method)
	checkParmTypes(t, method, "*int")
	checkReturnTypes(t, method, "*string")
}

func TestParseRichMethod(t *testing.T) {
	method := findMethod(t, "Rich", parseTypesMethods(t))
	checkFuncShape(t, 3, 3, method)
	checkParmTypes(t, method, "T", "[]T", "*T")
	checkReturnTypes(t, method, "*T", "T", "[]T")
}

func TestParseInterfaceMethod(t *testing.T) {
	method := findMethod(t, "Interface", parseTypesMethods(t))
	checkFuncShape(t, 1, 1, method)
	checkParmTypes(t, method, "I")
	checkReturnTypes(t, method, "I")
}

func TestParseMethodOfInterface(t *testing.T) {
	ifaces := parseTypesInterfaces(t)
	if len(ifaces) != 1 {
		t.Errorf("Unexpected number of interfaces: %d", len(ifaces))
	}
	methods := ifaces[0].Methods
	if len(methods) != 2 {
		t.Errorf("Unexpected number of methods: %d", len(methods))
	}

	method := findMethod(t, "MethodOfInterface", methods)
	checkFuncShape(t, 0, 0, method)

	method = findMethod(t, "AnotherMethodOfInterface", methods)
	checkFuncShape(t, 0, 0, method)
}

func TestParseImportedInterfaceMethod(t *testing.T) {
	method := findMethod(t, "ImportedInterface", parseTypesMethods(t))
	checkFuncShape(t, 1, 1, method)
	checkParmTypes(t, method, "io.Reader")
	checkReturnTypes(t, method, "io.Writer")
}

func TestParseDeeplyTypedMethod(t *testing.T) {
	method := findMethod(t, "DeeplyTyped", parseTypesMethods(t))
	checkFuncShape(t, 2, 1, method)
	checkParmTypes(t, method, "[]*int", "**[]**[2]**int")
	checkReturnTypes(t, method, "[]***string")
}

func TestParseEllipsisMethod(t *testing.T) {
	method := findMethod(t, "Ellipsis", parseTypesMethods(t))
	checkFuncShape(t, 1, 0, method)
	checkParmTypes(t, method, "...int")
}

func parseTypesPackage(t *testing.T) *architecture.Package {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	p := parser.NewParser(contexts.CreateBuildContext(contexts.CommandContext{
		GoPath: filepath.Join(pwd, "../fixtures/types"),
	}))

	p.ParseImportSpec("test_types")
	arch := p.GetArchitecture()
	return arch.FindDirectory("test_types").Package
}

func parseTypesMethods(t *testing.T) []*architecture.Method {
	return parseTypesPackage(t).Methods
}

func parseTypesInterfaces(t *testing.T) []*architecture.Interface {
	return parseTypesPackage(t).Interfaces
}

func findMethod(t *testing.T, name string, methods []*architecture.Method) *architecture.Method {
	for _, meth := range methods {
		if meth.Name == name {
			return meth
		}
	}
	t.Errorf("method %s not found", name)
	t.FailNow()
	return nil
}

func checkFuncShape(t *testing.T, numParms, numReturns int, method *architecture.Method) {
	if len(method.ParmTypes) != numParms {
		t.Errorf("Method %s parsed incorrectly: %d instead of %d parameters", method, len(method.ParmTypes), numParms)
	}
	if len(method.ReturnTypes) != numReturns {
		t.Errorf("Method %s parsed incorrectly: %d instead of %d return values", method, len(method.ReturnTypes), numReturns)
	}
}

func checkParmTypes(t *testing.T, method *architecture.Method, parmTypes ...architecture.Type) {
	if len(parmTypes) != len(method.ParmTypes) {
		t.Errorf("expected %d parameter(s), but found %d", len(parmTypes), len(method.ParmTypes))
	}
	for i, parmType := range parmTypes {
		if method.ParmTypes[i] != parmType {
			t.Errorf("Method %s parameter %d type parsed incorrectly as %s instead of %s", method, i, method.ParmTypes[i], parmType)
		}
	}
}

func checkReturnTypes(t *testing.T, method *architecture.Method, returnTypes ...architecture.Type) {
	if len(returnTypes) != len(method.ReturnTypes) {
		t.Errorf("expected %d return value(s), but found %d", len(returnTypes), len(method.ReturnTypes))
	}
	for i, returnType := range returnTypes {
		if method.ReturnTypes[i] != returnType {
			t.Errorf("Method %s return value %d type parsed incorrectly as %s instead of %s", method, i, method.ReturnTypes[i], returnType)
		}
	}
}
