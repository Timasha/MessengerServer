package logger

import (
	"encoding/json"
	"messengerServer/internal/api_objects/logger/responses"
	"time"

	"github.com/valyala/fasthttp"
)

type FileLogger struct {
	Ip  string
	TLS bool
}

//TODO: finish this function, add defered logging
func (f FileLogger) Log(logLevel LogLevel, logInfo string) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	if f.TLS {
		req.SetRequestURI("https://" + f.Ip + "/log/file")
	} else {
		req.SetRequestURI("http://" + f.Ip + "/log/file")
	}
	req.Header.Set("Content-Type", "plain/text")
	req.SetBodyString("[" + time.Now().Format("2006-01-02T15:04:05-0700") + "] " + "[" + string(logLevel) + "] : " + logInfo)
	var resp *fasthttp.Response
	reqErr := fasthttp.Do(req, resp)
	if reqErr != nil {
		go func() {
			time.Sleep(time.Minute * 10)
			f.Log(logLevel, logInfo)
		}()
		return
	}
	var respBody responses.LogFileResponse
	unmarshErr := json.Unmarshal(resp.Body(), &respBody)
	if unmarshErr != nil {

	}
}
