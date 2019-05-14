package navpool

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PoolClient struct {
	host       string
	network    string
	httpClient *http.Client
}

func NewClient(host string, network string) (c *PoolClient, err error) {
	if len(host) == 0 {
		err = errors.New("bad call missing argument host")
		return
	}

	c = &PoolClient{host: host, network: network, httpClient: &http.Client{}}
	return
}

func (c *PoolClient) call(url string, method string, data interface{}) (response []byte, err error) {
	var body *bytes.Buffer
	if data != nil {
		body = bytes.NewBufferString(data.(string))
	} else {
		body = bytes.NewBufferString("")
	}

	req, err := http.NewRequest(method, c.host+url, body)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Network", string(c.network))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respContent, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respContent))

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("HTTP error: %s - %s", resp.Status, string(respContent)))
		return
	}

	return respContent, nil
}
