package formats

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteDOT(t *testing.T) {
	net := buildNetwork()
	net.Name = "TestNN"

	b := bytes.NewBufferString("")
	err := WriteDOT(b, net)
	require.NoError(t, err, "failed to DOT encode")
	t.Log(b)
	assert.NotEmpty(t, b)
}

func TestWriteDOT_Write_Error(t *testing.T) {
	net := buildNetwork()
	net.Name = "TestNN"

	errWriter := ErrorWriter(1)
	err := WriteDOT(&errWriter, net)
	assert.EqualError(t, err, alwaysErrorText)
}
