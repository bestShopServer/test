package util

import (
	"io/ioutil"
	"github.com/json-iterator/go"
)

func LoadJson(filename string, v interface{}) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		Logger(ERROR_LEVEL, "LoadFile", "Load " + filename + " Err:" + err.Error())
		return
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		Logger(ERROR_LEVEL, "LoadFile", "JsonUnmarshal " + filename + " Err:" + err.Error())
		return
	}
}
