package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Log Level
const (
	INFO_LEVEL = iota
	ERROR_LEVEL
	FATAL_LEVEL
)

//启动服务
func InitLog() {
	//初始化日志打印文件
	fmt.Println("init logging to file ...")
	//f, _ := os.Create("log/debug.log")
	f, _ := os.OpenFile("log/debug.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	gin.DefaultWriter = io.MultiWriter(f)
}

//查看日志文件大小
func CheckLogFileSize() {
	var size int64
	err := filepath.Walk("log/debug.log", func(path string, info os.FileInfo, err error) error {
		size = info.Size()
		return nil
	})
	if err != nil {
		Error("ERROR:%v", err.Error())
	}
	Info("文件大小:%v", size)
	if size > 50000000 {
		Info("备份日志...")
		BackLogFile()
	}
}

//备份日志
func BackLogFile() {
	filename := "log/debug.log"
	mdate := time.Now().Format("20060102150405")
	err := os.Rename(filename, filename+"-"+mdate)
	if err != nil {
		Error("Move File ERROR[%v]", err.Error())
	}
	f, _ := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	gin.DefaultWriter = io.MultiWriter(f)
}

// Func: Print Different Level of Logs
func Logger(level int, category string, msg interface{}) {
	ct := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	filename := strings.Split(file, "/")
	switch level {
	case INFO_LEVEL:
		//fmt.Println(fmt.Sprintf("%s|INFO|%v|%v|%v|%v", ct, category, filename[len(filename)-1], line, msg))
		fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("%s|INFO|%v|%v|%v|%v", ct, category, filename[len(filename)-1], line, msg))
	case ERROR_LEVEL:
		//fmt.Println(fmt.Sprintf("%s|ERROR|%v|%v|%v|%v", ct, category, filename[len(filename)-1], line, msg))
		fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("%s|ERROR|%v|%v|%v|%v", ct, category, filename[len(filename)-1], line, msg))
	case FATAL_LEVEL:
		//fmt.Println(fmt.Sprintf("%s|FATAL|%v|%v|%v|%v", ct, category, filename[len(filename)-1], line, msg))
		fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("%s|FATAL|%v|%v|%v|%v", ct, category, filename[len(filename)-1], line, msg))
	default:
		//fmt.Println(fmt.Sprintf("%s|[DEBUG|%v|%v|%v|%v", ct, category, filename[len(filename)-1], line, msg))
		fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("%s|[DEBUG|%v|%v|%v|%v", ct, category, filename[len(filename)-1], line, msg))
	}
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

func Debug(format string, v ...interface{}) {
	ct := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	filename := strings.Split(file, "/")
	msg := fmt.Sprintf(formatLog(format, v...))
	//fmt.Println(fmt.Sprintf("%s|DEBUG|%v|%v|%v", ct, filename[len(filename)-1], line, msg))
	fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("%s|DEBUG|%v|%v|%v", ct, filename[len(filename)-1], line, msg))
}

func Info(format string, v ...interface{}) {
	ct := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	filename := strings.Split(file, "/")
	msg := fmt.Sprintf(formatLog(format, v...))
	//fmt.Println(fmt.Sprintf("%s|INFO|%v|%v|%v", ct, filename[len(filename)-1], line, msg))
	fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("%v|%v|%v|%v|%v", ct, "INFO", filename[len(filename)-1], line, msg))
}

func Error(format string, v ...interface{}) {
	ct := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	filename := strings.Split(file, "/")
	msg := fmt.Sprintf(formatLog(format, v...))
	//fmt.Println(fmt.Sprintf("%s|ERROR|%v|%v|%v", ct, filename[len(filename)-1], line, msg))
	fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("%s|ERROR|%v|%v|%v", ct, filename[len(filename)-1], line, msg))
}

func Fatal(format string, v ...interface{}) {
	ct := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	filename := strings.Split(file, "/")
	msg := fmt.Sprintf(formatLog(format, v...))
	//fmt.Println(fmt.Sprintf("%s|FATAL|%v|%v|%v", ct, filename[len(filename)-1], line, msg))
	fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("%s|FATAL|%v|%v|%v", ct, filename[len(filename)-1], line, msg))
}
