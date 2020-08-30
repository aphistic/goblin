package goblin

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenFSFile(t *testing.T) {
	tf, err := ioutil.TempFile("", testTempPattern)
	require.NoError(t, err)
	defer tf.Close()

	f := newOpenFSFile(tf.Name(), tf)
	require.NotNil(t, f)
	assert.Equal(t, tf.Name(), f.fullPath)
	assert.Equal(t, tf, f.f)
}

func TestOpenFSFileStat(t *testing.T) {
	tf, err := ioutil.TempFile("", testTempPattern)
	require.NoError(t, err)
	defer tf.Close()

	f := newOpenFSFile(tf.Name(), tf)
	require.NotNil(t, f)

	fi, err := f.Stat()
	require.NoError(t, err)
	assert.Contains(t, tf.Name(), fi.Name())
}
