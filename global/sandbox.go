package global

import (
	"DetectiveMasterServer/util"
	zmq "github.com/zeromq/goczmq"
)

var Task *util.TaskObj

type DBMessage struct {
	SendMsg []byte
	RecvMsg []byte
}

var SendQueue chan *DBMessage
var DBChannel *zmq.Channeler

const (
	ERR_DB_OK = iota
	ERR_DB_JSON
	ERR_DB_NOTFOUND
	ERR_DB_INTERNAL_SERVER
	ERR_DB_NOTFOUND_DATA
	ERR_DB_PARAMS_MISSING
	ERR_DB_PARAMS
	ERR_DB_EXISTS
	ERR_DB_NOT_EXISTS
)

func NewDBRequest(method string, params interface{}) map[string]interface{} {
	request := make(map[string]interface{})
	request["method"] = method
	request["params"] = params
	return request
}

func UnwrapPackage(data map[string]interface{}) (int, interface{}) {
	var m string
	var e int

	msg, ok := data["msg"]

	if ok {
		m = msg.(string)
	}

	err, ok := data["err"]

	if ok {
		e = int(err.(float64))
	}

	params, _ := data["params"]

	if e != 0 {
		util.Logger(util.ERROR_LEVEL, "Sandbox", m)
	}

	return e, params
}

func UnwrapArrayPackage(data map[string]interface{}) (int, []interface{}) {
	code, obj := UnwrapPackage(data)

	var wrap []interface{}

	if obj != nil {
		wrap, _ = obj.([]interface{})
	}

	return code, wrap
}

func UnwrapObjectPackage(data map[string]interface{}) (int, map[string]interface{}) {
	code, obj := UnwrapPackage(data)

	var wrap map[string]interface{}

	if obj != nil {
		wrap, _ = obj.(map[string]interface{})
	}

	return code, wrap
}
