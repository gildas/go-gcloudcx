package gcloudcx_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func LoadObject(filename string, object interface{}) (err error) {
	payload, err := LoadFile(filename)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(payload, &object); err != nil {
		return err
	}
	return nil
}

func LoadFile(filename string) (payload []byte, err error) {
	if payload, err = ioutil.ReadFile(filepath.Join(".", "testdata", filename)); err != nil {
		return nil, err
	}
	return payload, nil
}
