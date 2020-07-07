package debug

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
)

func loadJsonToBytes(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Err: ReadFile ", err)
		return data
	} else {
		return data
	}
}

func loadJson(filename string, v interface{}) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Err: ReadFile ", err)
		return
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		log.Println("Err: JsonUnmarshal ", err)
		return
	}
}
