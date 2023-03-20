package blockexplorerclient

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"code.cryptopower.dev/group/instantswap/blockexplorer/global/errors"
	"code.cryptopower.dev/group/instantswap/blockexplorer/global/utils"
)

const (
	defaultHttpClientTimeout = 30

	//outputResponse = true
)

type HandleRequest func(r *http.Request)

// Client is the base for http calls
type Client struct {
	apiBase        string
	libName        string
	httpClient     *http.Client
	Debug          bool
	OutputResponse bool
	handleRequest  HandleRequest
}

// NewClient return a new HTTP client
func NewClient(apiBase, apiSecret string, enableOutput bool, handleRequest HandleRequest) (c *Client) {
	return &Client{
		apiBase:        apiBase,
		libName:        apiSecret,
		httpClient:     &http.Client{},
		Debug:          false,
		OutputResponse: enableOutput,
		handleRequest:  handleRequest,
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

// doTimeoutRequest do a HTTP request with timeout
func (c *Client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		if c.Debug {
			c.dumpRequest(req)
		}
		resp, err := c.httpClient.Do(req)
		if c.Debug {
			c.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		errTimeout := &errors.Error{Err: errors.New("timeout on reading data from " + c.libName + " API"), Kind: 8}
		return nil, errTimeout
	}
}

// Do do prepare and process HTTP request to API
func (c *Client) Do(method, path string, payload interface{}, authNeeded bool) (response []byte, err error) {
	var connectTimer *time.Timer
	//connectTimer := time.NewTimer(DEFAULT_HTTPCLIENT_TIMEOUT * time.Second)
	var rawUrl string
	if strings.HasPrefix(path, "http") {
		rawUrl = path
	} else {
		rawUrl = fmt.Sprintf("%s%s", c.apiBase, path)
	}
	var req *http.Request

	reqInfo := AuthInfo{
		exchange: c.libName,
		c:        c,
		method:   method,
		payload:  payload,
		url:      rawUrl,
	}
	reqResult, err := getRequestType(reqInfo)
	if err != nil {
		return nil, err
	}
	if c.handleRequest != nil {
		c.handleRequest(reqResult.request)
	}
	req = reqResult.request
	connectTimer = reqResult.connectTimer

	if req == nil {
		err = errors.New("blockexplorerclient error: request was nil")
		return nil, err
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)

	//test
	if c.OutputResponse {
		fmt.Println(fmt.Sprintf("reponse %s", response), err)
	}

	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		var errStr string
		responseStr := string(response)
		if responseStr != "" {
			if strings.Contains(strings.ToLower(responseStr), "<body>") {
				responseStr = utils.GetStringBefore(responseStr, "<body>")
			}
		}
		res := "'" + responseStr + "'"
		if c.OutputResponse {
			errStr = "blockexplorerclient:error:" + resp.Status + ":" + res
		} else {
			errStr = res
		}

		err = errors.New(errStr)
	}
	return response, err
}

func getRequestType(info AuthInfo) (result AuthInfo, err error) {
	result.connectTimer = time.NewTimer(defaultHttpClientTimeout * time.Second)
	req, err := http.NewRequest(info.method, info.url, strings.NewReader(info.payload.(string)))
	if err != nil {
		return result, err
	}
	if info.method == "POST" || info.method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}

	req.Header.Add("Accept", "application/json")
	result.request = req
	return
}
func sleepPrint(dur int) {
	send_delay := dur
	for j := 1; j <= send_delay; j++ {
		fmt.Printf(".")
		////
		time.Sleep(time.Second)
	}
	// fmt.Println("waiting.......")
	// time.Sleep(time.Second * dur)
}

type AuthInfo struct {
	exchange     string
	c            *Client
	request      *http.Request
	payload      interface{}
	method       string
	url          string
	resource     string //only used for coinswitch right now
	connectTimer *time.Timer
}
