package httputil

import (
	"io/ioutil"
	"net/http"
	"vsync/log"
)

func Discard(res *http.Response, err error) {
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	defer ioutil.ReadAll(res.Body)
}