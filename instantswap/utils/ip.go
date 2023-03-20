package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetPublicIP() (ip string, err error) {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ipAddress := fmt.Sprintf("%s", response)
	if ipAddress == "" {
		err = errors.New("myexternalip.com returned a blank ip address... ")
		return "", err
	}
	ipAddress = strings.TrimSpace(ipAddress)

	return ipAddress, nil
}
