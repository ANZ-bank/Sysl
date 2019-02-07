package main

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.WarnLevel)
}

func TestGenerateCode(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", ".", "tests/test.gen.sysl", "tests/test.gen.g", "javaFile")
	root := output[0].output
	assert.Equal(t, 1, len(output), "unexpected length of output")
	assert.Equal(t, 1, len(root), "unexpected length of root")
	n1 := root[0].(Node)
	assert.Equal(t, 4, len(n1), "unexpected length of javaFile")
	package1 := n1[0].(Node)
	comment1 := n1[1].(Node)
	import1 := n1[2].(Node)
	definition1 := n1[3].(string)
	assert.Equal(t, 1, len(package1), "unexpected length of package")
	assert.Equal(t, 2, len(comment1), "unexpected length of comment")
	assert.Equal(t, 2, len(import1), "unexpected length of import")
	assert.Equal(t, "some_value", definition1, "unexpected value of definition")

	package2 := package1[0].(Node)
	assert.Equal(t, 3, len(package2), "unexpected length of package2")
	assert.Equal(t, "com.example.gen", package2[1].(string), "unexpected length of package2")

	comment := []string{"comment1", "comment2"}
	for i := range comment {
		comment_0 := comment1[i].(Node)
		assert.Equal(t, 1, len(comment_0), "unexpected length of comment2")
		comment_0_0 := comment_0[0].(string)
		assert.Equal(t, comment[i], comment_0_0, "unexpected length of comment_i")
	}

	imports := []string{"import1", "import2"}
	for i := range imports {
		import_0 := import1[i].(Node)
		assert.Equal(t, 1, len(import_0), "unexpected length of import2")
		import_0_0 := import_0[0].(Node)
		assert.Equal(t, 3, len(import_0_0), "unexpected length of import2")
		assert.Equal(t, imports[i], import_0_0[1].(string), "unexpected length of import_i")
	}
}

func TestGenerateCodeNoComment(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", ".", "tests/test.gen_no_comment.sysl", "tests/test.gen.g", "javaFile")
	assert.Equal(t, 1, len(output), "unexpected length of output")
	root := output[0].output
	assert.Equal(t, 1, len(root), "unexpected length of root")
	n1 := root[0].(Node)
	assert.Equal(t, 3, len(n1), "unexpected length of javaFile")
	package1 := n1[0].(Node)
	import1 := n1[1].(Node)
	definition1 := n1[2].(string)
	assert.Equal(t, 1, len(package1), "unexpected length of package")
	assert.Equal(t, 2, len(import1), "unexpected length of comment")
	assert.Equal(t, "some_value", definition1, "unexpected value of definition")

	package2 := package1[0].(Node)
	assert.Equal(t, 3, len(package2), "unexpected length of package2")
	assert.Equal(t, "com.example.gen", package2[1].(string), "unexpected length of package2")

	imports := []string{"import1", "import2"}
	for i := range imports {
		import_0 := import1[i].(Node)
		assert.Equal(t, 1, len(import_0), "unexpected length of import2")
		import_0_0 := import_0[0].(Node)
		assert.Equal(t, 3, len(import_0_0), "unexpected length of import2")
		assert.Equal(t, imports[i], import_0_0[1].(string), "unexpected length of import_i")
	}
}

func TestGenerateCodeNoPackage(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", ".", "tests/test.gen_no_package.sysl", "tests/test.gen.g", "javaFile")
	root := output[0].output
	assert.Nil(t, root, "unexpected root")
}

func TestGenerateCodeMultipleAnnotations(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", ".", "tests/test.gen_multiple_annotations.sysl", "tests/test.gen.g", "javaFile")
	root := output[0].output
	assert.Nil(t, root, "unexpected root")
}

func TestGenerateCodePerType(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", ".", "tests/multiple_file.gen.sysl", "tests/test.gen.g", "javaFile")
	assert.Equal(t, 1, len(output), "unexpected length of output")
	assert.Equal(t, "Request.java", output[0].filename, "unexpected length of output")

	root := output[0].output
	assert.Equal(t, 1, len(root), "unexpected length of javaFile")

	requestRoot := root[0].(Node)
	assert.Equal(t, 4, len(requestRoot), "unexpected length of requestRoot")

	package1 := requestRoot[0].(Node)
	comment1 := requestRoot[1].(Node)
	import1 := requestRoot[2].(Node)
	definition1 := requestRoot[3].(string)
	assert.Equal(t, 1, len(package1), "unexpected length of package")
	assert.Equal(t, 2, len(comment1), "unexpected length of comment")
	assert.Equal(t, 2, len(import1), "unexpected length of import")
	assert.Equal(t, "Request", definition1, "unexpected value of definition")
}

func TestSerialize(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", ".", "tests/test.gen.sysl", "tests/test.gen.g", "javaFile")
	out := new(bytes.Buffer)
	Serialize(out, " ", output[0].output)
	golden := "package com.example.gen \n comment1 comment2 import import1 \n import import2 \n some_value "
	assert.Equal(t, golden, out.String(), "unexpected value of out string")
}
