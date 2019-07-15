package dao

import (
	xjson "github.com/buger/jsonparser"
	"go-web-demo/library/log"
	"io/ioutil"
	"net/http"
)

var (
	url = "http://fantuan.bz.mgtv.com/fantuan/soKeyword"
)

func (d *Dao) SearchKeyword() (keyword string, err error) {
	resp, err := d.httpClient.Get(url)
	if err != nil {
		log.Error("SearchKeyword error(%v)", err)
		return
	}
	// see: https://golang.org/pkg/net/http/
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("http request status:%s", resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("SearchKeyword read body error(%v)", err)
		return
	}

	keyword, err = xjson.GetString(body, "data", "keyword")
	if err != nil {
		log.Error("SearchKeyword read body error(%v)", err)
		return
	}

	log.Info("SearchKeyword keyword = %s", keyword)
	return
}
