package Common

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func HttpSetExtraHeader(header string)  {
	
}

func HttpClient(timeout time.Duration) (client *http.Client, err error) {
	client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(timeout))
				return conn, nil
			},
			ResponseHeaderTimeout: timeout,
		},
	}
	return
}

func HttpGet(timeout time.Duration, url string) (body []byte, err error) {
	client, err := HttpClient(timeout)
	if nil != err {
		return nil, err
	}
	resp, err := client.Get(url)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	return

}

func HttpPostForm(timeout time.Duration, url string, data url.Values) (body []byte, err error) {
	client, err := HttpClient(timeout)
	if nil != err {
		return nil, err
	}
	resp, err := client.PostForm(url, data)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	return

}
func HttpGetRealTargetURL(timeout time.Duration, url string) (realTargetURL string, urlList []map[int]string, err error) {

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(timeout))
				return conn, nil
			},
			ResponseHeaderTimeout: timeout,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			log.Println("CheckRedirect:%#v",req.Response.StatusCode, req.URL.String())
			stat := make(map[int]string)
			stat[req.Response.StatusCode] = req.URL.String()
			urlList = append(urlList, stat)
			return nil
		},
	}
	stat := make(map[int]string)
	stat[302] = url
	urlList = append(urlList, stat)

	resp, err := client.Get(url)

	if nil == err {
		_lastURL := resp.Request.URL.String()

		if nil != err {
			return "", []map[int]string{}, err
		}
		if _lastURL == url {
			return url, []map[int]string{map[int]string{resp.StatusCode:url}}, nil
		}
		return _lastURL, urlList, nil
	} else {
		return "", urlList, err
	}

	return
}
