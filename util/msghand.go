package util

import (
	"encoding/json"
	"github.com/pebbe/zmq4"
	//zmq "github.com/zeromq/goczmq"
	"sync"
)

type TaskComm struct {
	*TaskObj
	//ZmqSock  *zmq.Sock
	SendData []byte
	RecvData []byte
	Error    error
}

//var Sock *zmq.Sock
var MQSock *zmq4.Socket
var wg sync.WaitGroup
var mutex sync.Mutex

//var TaskWorker chan *TaskComm
//
//func Connect(addr string, timeout int) (*zmq.Sock, error) {
//	s := zmq.NewSock(zmq.Req)
//
//	err := s.Connect(addr)
//	if err != nil {
//		return nil, err
//	}
//	s.SetRcvtimeo(timeout)
//
//	return s, nil
//}

func InitZeroMQClient(addr string, timeout int) {
	var err error
	//Sock = zmq.NewSock(zmq.Req)
	//
	//err := Sock.Connect(addr)
	//if err != nil {
	//	Error("ERROR:%v", err.Error())
	//}
	//Sock.SetRcvtimeo(timeout)

	//context, err := gozmq.NewContext()
	//if err != nil {
	//	Error("ERROR:%v", err.Error())
	//}
	//MQSock, err = context.NewSocket(gozmq.REQ)
	//if err != nil {
	//	Error("ERROR:%v", err.Error())
	//}
	//err = MQSock.Connect(addr)
	//if err != nil {
	//	Error("ERROR:%v", err.Error())
	//}

	MQSock, err = zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		Error("ERROR:%v", err.Error())
	}
	if MQSock != nil {
		err = MQSock.Connect(addr)
		if err != nil {
			Error("ERROR:%v", err.Error())
		}
		Debug("init socket connect ...")
	} else {
		Error("MQSock ERROR: is null")
		return
	}

}

//
//func TaskNew(s *zmq.Sock, data []byte) ([]byte, error) {
//	Debug("TaskNew...")
//	defer s.Destroy()
//
//	err := s.SendFrame(data, zmq.FlagNone)
//	if err != nil {
//		return nil, err
//	}
//	Debug("SendFrame...")
//
//	msg, _, err := s.RecvFrame()
//	if err != nil {
//		return nil, err
//	}
//	Debug("RecvFrame...")
//
//	err = s.SendFrame(data, zmq.FlagNone)
//	if err != nil {
//		return nil, err
//	}
//	Debug("SendFrame2...")
//
//	msg, _, err = s.RecvFrame()
//	if err != nil {
//		return nil, err
//	}
//	Debug("RecvFrame2...")
//
//	err = s.SendFrame(data, zmq.FlagNone)
//	if err != nil {
//		return nil, err
//	}
//	Debug("SendFrame3...")
//
//	msg, _, err = s.RecvFrame()
//	if err != nil {
//		return nil, err
//	}
//	Debug("RecvFrame3...")
//
//	return msg, nil
//}

func TaskNew2(data []byte) (msg []byte, err error) {
	Debug("TaskNew...")
	mutex.Lock()
	l, err := MQSock.SendBytes(data, 0)
	if err != nil {
		Error("ERROR:%v", err.Error())
		return nil, err
	}
	Debug("SendFrame len %v ...", l)

	msg, err = MQSock.RecvBytes(0)
	if err != nil {
		Error("ERROR:%v", err.Error())
		return nil, err
	}
	Debug("RecvFrame ...")
	mutex.Unlock()

	return msg, err
}

//func HandDBMsg(obj *TaskComm) {
//	Debug("HandDBMsg ...")
//	mutex := sync.Mutex{}
//	mutex.Lock()
//	s, err := Connect(obj.Addr, obj.Timeout)
//	if err != nil {
//		Error("ERROR:%v", err.Error())
//	}
//
//	obj.RecvData, err = TaskNew(s, obj.SendData)
//	if err != nil {
//		Error("ERROR:%v", err.Error())
//	}
//	mutex.Unlock()
//	//wg.Done()
//}

func (this *TaskObj) TaskJsonNew(data interface{}) (jsonResp map[string]interface{}, err error) {
	Debug("data[%v]", data)

	jraw, err := json.Marshal(data)
	if err != nil {
		Error("ERROR:%v", err.Error())
		return jsonResp, err
	}
	//var comm TaskComm
	//comm.TaskObj = this
	//comm.SendData = jraw

	//wg.Add(1)
	//go HandDBMsg(&comm)
	//HandDBMsg(&comm)
	res, err := TaskNew2(jraw)
	if err != nil {
		Error("ERROR:%v", err.Error())
		return jsonResp, err
	}

	//wg.Wait()

	err = json.Unmarshal(res, &jsonResp)
	if err != nil {
		Error("ERROR[%v]", err.Error())
		return jsonResp, err
	}
	Debug("jsonResp...")

	return jsonResp, err
}
