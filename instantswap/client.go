package instantswap

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/vibros68/instantswap/instantswap/utils"
)

const (
	defaultHttpClientTimeout = 30
)

type CustomReqFunc func(r *http.Request, body string) error

type Client struct {
	exchange      string
	httpClient    *http.Client
	conf          *ExchangeConfig
	handleRequest CustomReqFunc
}

// NewClient return a new HTTP client
func NewClient(exchange string, conf *ExchangeConfig, handleRequests ...CustomReqFunc) (c *Client) {
	client := &Client{
		exchange:   exchange,
		conf:       conf,
		httpClient: &http.Client{},
	}
	if len(handleRequests) >= 1 {
		client.handleRequest = handleRequests[0]
	}
	return client
}

func (c *Client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		if c.conf.Debug {
			c.dumpRequest(req)
		}
		resp, err := c.httpClient.Do(req)
		if c.conf.Debug {
			c.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, fmt.Errorf("timeout on reading data from [%s] api", c.exchange)
	}
}

func (c Client) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (c Client) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

// Do do prepare and process HTTP request to API
func (c *Client) Do(apibase, method, resource string, payload string, authNeeded bool) (response []byte, err error) {
	var connectTimer = time.NewTimer(defaultHttpClientTimeout * time.Second)
	var rawurl string
	if strings.HasPrefix(resource, "http") {
		rawurl = resource
	} else {
		rawurl = fmt.Sprintf("%s%s", apibase, resource)
	}
	var req *http.Request
	req, err = http.NewRequest(method, rawurl, strings.NewReader(payload))
	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Set("Accept", "application/json")
	}
	if c.handleRequest != nil {
		err = c.handleRequest(req, payload)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)

	//test
	if c.conf.Debug {
		fmt.Printf("\n|*** URL %s RESPONSE ***|\n", req.URL)
		if err != nil {
			fmt.Printf("%s err: %s", response, err.Error())
		} else {
			fmt.Printf("%s", response)
		}
		fmt.Printf("\n|*** END RESPONSE ***|\n")
	}

	if err != nil {
		return response, err
	}
	if resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusTooManyRequests {
			return response, TooManyRequestsError
		}
		var errStr string
		responseStr := string(response)
		if responseStr != "" {
			if strings.Contains(strings.ToLower(responseStr), "<body>") {
				responseStr = utils.GetStringBefore(responseStr, "<body>")
			}
		}
		res := "'" + responseStr + "'"
		errStr = "\nexchangeclient:error:" + resp.Status + ":" + res
		err = fmt.Errorf(errStr)
	}
	return response, err
}
