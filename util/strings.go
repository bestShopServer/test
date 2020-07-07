package util

import (
	"fmt"
	"math/rand"
	"time"
)

// Func: Generate Random String by Length
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Func: Generate Random Int By Length
func GetRandomInt(l int) string {
	str := "123456789"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Func: Generate Random Int By Length
func GenOrderNo() string {
	result := fmt.Sprintf("%v%v", time.Now().Unix(), GetRandomInt(6))
	return result
}
