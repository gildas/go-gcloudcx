package purecloud_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gildas/go-logger"
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

func CreateLogger(filename string) *logger.Logger {
	folder := filepath.Join(".", "log")
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}
	return logger.CreateWithStream("test", &logger.FileStream{Path: filepath.Join(folder, filename)})
}
