package util

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
)

func HttpRequest(url string) (int, []byte, error) {
	var client http.Client
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 500, nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s/v%s", common.APPNAME, common.APPVERSION))
	res, getErr := client.Do(req)
	if getErr != nil {
		return res.StatusCode, []byte(res.Status), getErr
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return res.StatusCode, []byte(res.Status), readErr
	}
	return res.StatusCode, body, nil
}
