package purecloud_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Load(filename string, object interface{}) (err error) {
	var payload []byte

	if payload, err = ioutil.ReadFile(filepath.Join(".", "testdata", filename)); err != nil {
		return err
	}
	if err = json.Unmarshal(payload, &object); err != nil {
		return err
	}
	return nil
}

func RequireEqualJSON(t *testing.T, filename string, payload []byte) {
	expected, err := ioutil.ReadFile(filepath.Join(".", "testdata", filename))
	require.Nil(t, err, "Failed to load %s", filename)
	require.JSONEq(t, string(expected), string(payload))
}