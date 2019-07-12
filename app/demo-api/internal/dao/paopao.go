package dao

import (
	xjson "github.com/buger/jsonparser"
	"go-web-demo/library/log"
	"io/ioutil"
	"net/http"
)

var (
	url = "https://paopao.iqiyi.com/apis/e/paopao/search/searchkeyword.action?agenttype=116&agentversion=9.7.5&appid=44&authcookie=e4tVLHncJZXL3wm3ZfFwQv0K6AjPMMtxGDLoV9ui1lPgam2Eph6KgpA0duNFHm17V6K0Hfb&business_type=1&device_id=12447C8C-9FDE-4383-B9A2-2DD5B7FD918B&iOSVersion=12.3&m_device_id=12447C8C-9FDE-4383-B9A2-2DD5B7FD918B&openudid=12447C8C-9FDE-4383-B9A2-2DD5B7FD918B&playPlatform=12&qypid=02032001010000000000&sign=d95d97d309a0d958538f6df7b2a50105&sourceid=44&timestamp=1562830151326&uid=1481574257"
)

func (d *Dao) SearchKeyword() (keyword string, err error) {
	res, err := d.httpClient.Get(url)
	if err != nil || http.StatusOK != res.StatusCode {
		log.Error("SearchKeyword error(%v)", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
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
