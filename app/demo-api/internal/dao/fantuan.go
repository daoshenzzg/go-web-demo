package dao

import (
	"context"
	"go-web-demo/library/log"
)

var (
	_OK = 200
)

func (d *Dao) SearchKeyword(ctx context.Context) (keyword string, err error) {
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Keyword string `json:"keyword"`
		}
	}

	if err = d.fantuanClient.Get(ctx, d.searchKeywordURL, nil, &res); err != nil {
		log.Error("d.fantuanClient.Get(%s) error(%v)", d.searchKeywordURL, err)
		return
	}
	if _OK != res.Code {
		log.Error("d.fantuanClient.Get(%s) code error(%v)", d.searchKeywordURL, res)
		return
	}

	return res.Data.Keyword, nil
}
