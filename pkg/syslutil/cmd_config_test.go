package syslutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCMDFlags(t *testing.T) {
	t.Parallel()
	var flags []string
	var err error

	flags, err = ReadCMDFlags("tests/config.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(flags))
	assert.Equal(t, "--grammar=go.gen.g", flags[0])
	assert.Equal(t, "--transform=go.gen.sysl", flags[1])
	assert.Equal(t, "--app-name=Test", flags[2])
	assert.Equal(t, "model.sysl", flags[3])

	flags, err = ReadCMDFlags("tests/config1.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(flags))
	assert.Equal(t, "--grammar=go.gen.g\"", flags[0])
	assert.Equal(t, "--transform=\"go.gen.sysl", flags[1])
	assert.Equal(t, "--app-name=Test POC", flags[2])
	assert.Equal(t, "model.sysl", flags[3])

	flags, err = ReadCMDFlags("tests/config2.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(flags))
	assert.Equal(t, "--grammar=go.gen.g\"", flags[0])
	assert.Equal(t, "--transform=\"go.gen.sysl", flags[1])
	assert.Equal(t, "--app-name=Test \" POC", flags[2])
	assert.Equal(t, "model.sysl", flags[3])
}

func TestPopulateCMDFlagsFromFile(t *testing.T) {
	t.Parallel()
	var cmdArgs, flags []string
	var err error
	///////////
	cmdArgs = []string{"sysl", "codegen", "@tests/config1.txt"}
	flags, err = PopulateCMDFlagsFromFile(cmdArgs)
	assert.Equal(t, nil, err)
	assert.Equal(t, 6, len(flags))
	assert.Equal(t, "sysl", flags[0])
	assert.Equal(t, "codegen", flags[1])
	assert.Equal(t, "--grammar=go.gen.g\"", flags[2])
	assert.Equal(t, "--transform=\"go.gen.sysl", flags[3])
	assert.Equal(t, "--app-name=Test POC", flags[4])
	assert.Equal(t, "model.sysl", flags[5])
	///////////
	cmdArgs = []string{"sysl", "codegen", "@tests/config.txt"}
	flags, err = PopulateCMDFlagsFromFile(cmdArgs)
	assert.Equal(t, nil, err)
	assert.Equal(t, 6, len(flags))
	assert.Equal(t, "sysl", flags[0])
	assert.Equal(t, "codegen", flags[1])
	assert.Equal(t, "--grammar=go.gen.g", flags[2])
	assert.Equal(t, "--transform=go.gen.sysl", flags[3])
	assert.Equal(t, "--app-name=Test", flags[4])
	assert.Equal(t, "model.sysl", flags[5])
}
