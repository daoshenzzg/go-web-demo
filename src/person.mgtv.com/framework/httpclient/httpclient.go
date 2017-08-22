package httpclient

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"person.mgtv.com/framework/logs"
)

const (
	DEFAULT_IDLE_CONN_TIMEOUT = 30 * time.Second
	DEFAULT_CONN_TIMEOUT      = 1 * time.Second
	DEFAULT_RW_TIMEOUT        = 1 * time.Second
	DEFAULT_MAX_IDLE_CONN     = 100
)

var httpclient *http.Client

func init() {
	httpclient = NewClient(DEFAULT_CONN_TIMEOUT, DEFAULT_RW_TIMEOUT, DEFAULT_MAX_IDLE_CONN)
	logs.GetLogger("system").Infof("HttpClient[connTimeout=%v, rwTimeout=%v, maxIdleConn=%v] init success.",
		DEFAULT_CONN_TIMEOUT, DEFAULT_CONN_TIMEOUT, DEFAULT_MAX_IDLE_CONN)
}

// TimeoutDialer implements our own dialer in order to set conn and read and write idle timeouts.
func TimeoutDialer(connTimeout, rwTimeout time.Duration) func(net, addr string) (net.Conn, error) {
	return func(network, address string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, address, connTimeout)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func NewClient(connTimeout, rwTimeout time.Duration, maxIdleConn int) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial:                TimeoutDialer(connTimeout, rwTimeout),
			MaxIdleConnsPerHost: maxIdleConn,
			IdleConnTimeout:     DEFAULT_IDLE_CONN_TIMEOUT,
		},
	}
}

func Get(url string) ([]byte, error) {
	resp, err := httpclient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Post(url string, params url.Values) ([]byte, error) {
	resp, err := httpclient.Post(url, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
