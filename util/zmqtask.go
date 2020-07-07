package util

import (
	"encoding/json"
	zmq "github.com/zeromq/goczmq"
)

func connect(addr string, timeout int) (*zmq.Sock, error) {
	s := zmq.NewSock(zmq.Req)

	err := s.Connect(addr)

	if err != nil {
		Error("ERROR[%v]", err.Error())
		return nil, err
	}

	s.SetRcvtimeo(timeout)

	return s, nil
}

type ZmqResult struct {
	Data []byte
	Err  error
}

// type task struct {
// 	input  []byte
// 	output chan *ZmqResult
// }

// type zmq_worker_pool struct {
// 	sock     *zmq.Sock
// 	id       int
// 	timeout  int
// 	receiver chan *task
// 	done     chan int
// }

// func (this *zmq_worker_pool) gotTask() {
// 	for {
// 		task, ok := <-this.receiver

// 		if !ok {
// 			this.done <- 1
// 			break
// 		}

// 		for this.sock.Pollin() {
// 			this.sock.RecvFrame()
// 		}

// 		result := new(ZmqResult)

// 		err := this.sock.SendFrame(task.input, zmq.FlagNone)

// 		if err != nil {
// 			result.Err = err
// 			task.output <- result
// 		} else {
// 			msg, _, err := this.sock.RecvFrame()
// 			if err != nil {
// 				result.Err = err
// 				task.output <- result
// 			} else {
// 				result.Data = msg
// 				task.output <- result
// 			}
// 		}
// 	}
// }

// type ZmqPool struct {
// 	worker []*zmq_worker_pool
// 	i      int
// 	c      int
// 	s      bool
// 	rwmx   *sync.RWMutex
// }

// func NewZmqPool(worker int, timeout int, addr string) (*ZmqPool, error) {
// 	pool := new(ZmqPool)
// 	pool.s = false
// 	pool.worker = make([]*zmq_worker_pool, worker)
// 	pool.c = worker
// 	pool.i = 0
// 	pool.rwmx = new(sync.RWMutex)

// 	for i := 0; i < worker; i++ {
// 		w := new(zmq_worker_pool)
// 		s, err := connect(addr, timeout)

// 		if err != nil {
// 			return nil, err
// 		}

// 		w.receiver = make(chan *task)
// 		w.done = make(chan int)
// 		w.sock = s
// 		go w.gotTask()
// 		w.timeout = timeout
// 		w.id = i + 1
// 		pool.worker[i] = w
// 	}

// 	return pool, nil
// }

type TaskObj struct {
	Addr     string
	Timeout  int
	Host     string
	MaxQueue int
}

func (this *TaskObj) Task(data []byte) (*ZmqResult, error) {
	s, err := connect(this.Addr, this.Timeout)
	if err != nil {
		Error("ERROR[%v]", err.Error())
	}
	defer s.Destroy()

	err = s.SendFrame(data, zmq.FlagNone)
	if err != nil {
		Error("ERROR[%v]", err.Error())
		return nil, err
	}

	msg, _, err := s.RecvFrame()

	if err != nil {
		Error("ERROR[%v]", err.Error())
		return nil, err
	}

	return &ZmqResult{
		msg,
		nil,
	}, nil
}

//func SendMessageQueue(sock *zmq.Sock) {
//	for {
//		identity := <-global.SendQueue
//		err := sock.SendFrame(identity, zmq.FlagNone)
//		if err != nil {
//			Error("ERROR:%v", err.Error())
//		}
//	}
//}
//
//func QueueTask(sock *zmq.Sock) {
//	Debug("TaskNew...")
//	mx := sync.Mutex{}
//	for {
//		mx.Lock()
//		tmp := <-global.SendQueue
//		err := sock.SendFrame(tmp.SendMsg, zmq.FlagNone)
//		if err != nil {
//			Error("ERROR:%v", err.Error())
//		}
//
//		msg, i, err := sock.RecvFrame()
//		if err != nil {
//			Error("ERROR:%v", err.Error())
//		}
//		Debug("i:%v", i)
//		tmp.RecvMsg = msg
//		mx.Unlock()
//	}
//
//}

func (this *TaskObj) TaskJson(data interface{}) (map[string]interface{}, error) {
	jraw, err := json.Marshal(data)
	if err != nil {
		Error("ERROR[%v]", err.Error())
		return nil, err
	}

	//global.SendQueue <- NewIdentityPackage(jraw)
	result, err := this.Task(jraw)
	//result, err := TaskNew(jraw)
	if err != nil {
		Error("ERROR[%v]", err.Error())
		return nil, err
	}

	var jsonResp map[string]interface{}

	err = json.Unmarshal(result.Data, &jsonResp)
	if err != nil {
		Error("ERROR[%v]", err.Error())
		return nil, err
	}

	return jsonResp, nil
}

// func (this *ZmqPool) Shutdown() {
// 	this.rwmx.Lock()
// 	this.s = true
// 	this.rwmx.Unlock()

// 	for i := 0; i < this.c; i++ {
// 		w := this.worker[i]
// 		close(w.receiver)
// 		<-w.done
// 		close(w.done)
// 		w.sock.Destroy()
// 	}
// }
