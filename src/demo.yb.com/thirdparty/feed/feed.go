package feed

import (
	"encoding/json"
	"fmt"

	"person.mgtv.com/framework/config"
	"person.mgtv.com/framework/httpclient"
	"person.mgtv.com/framework/logs"
)

// 判断是否关注
func IsFollowed(uid, artistId string) (bool, error) {
	url := config.Section("thirdparty.feed").Key("is_followed_url").String()
	url = fmt.Sprintf(url, uid, artistId)

	body, err := httpclient.Get(url)
	if err != nil {
		return false, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logs.GetLogger("system").Errorf("url=%s, response=%s", url, string(body))
		return false, err
	}

	logs.GetLogger("system").Debugf("Response data is '%s'", data)

	return true, nil
}
