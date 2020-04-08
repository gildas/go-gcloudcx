package purecloud_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
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
