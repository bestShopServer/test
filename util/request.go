package util

import (
	"github.com/valyala/fasthttp"
	"time"
)

// Func: Send Key Value Request
func SendKVReq(arg *fasthttp.Args, method string, requestURI string, cookies map[string]interface{}) ([]byte, int, error) {
	req := &fasthttp.Request{}
	switch method {
	case "GET":
		req.Header.SetMethod(method)
		requestURI = requestURI + "?" + arg.String()
	case "POST":
		req.Header.SetMethod(method)
		arg.WriteTo(req.BodyWriter())
	}
	if cookies != nil {
		for key, v := range cookies {
			req.Header.SetCookie(key, v.(string))
		}
	}
	req.SetRequestURI(requestURI)

	resp := &fasthttp.Response{}
	err := fasthttp.DoTimeout(req, resp, time.Second*30)

	return resp.Body(), resp.StatusCode(), err
}

// Func: Send Json Request
func SendJsonReq(method string, url, bodyjson string) ([]byte, int, error) {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	switch method {
	case "GET":
		req.Header.SetMethod(method)
	case "POST":
		req.Header.SetMethod(method)
	}

	req.Header.SetContentType("application/json")
	req.SetBodyString(bodyjson)

	req.SetRequestURI(url)

	err := fasthttp.DoTimeout(req, resp, time.Second*30)
	return resp.Body(), resp.StatusCode(), err
}
