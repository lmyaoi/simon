package httputil

import (
	"io/ioutil"
	"net/http"
	"simon/log"
	"time"
)

func Discard(res *http.Response, err error) {
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	defer ioutil.ReadAll(res.Body)
}

func Retry(client *http.Client, req *http.Request, retries int) (res *http.Response, err error) {
	for i := 0; i < 1+retries; i++ {
		res, err = client.Do(req)
		if err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}
