package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func MakeRequest(url string) (response []byte, err error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP call was unsuccessful, response code: %d", resp.StatusCode))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return []byte(body), nil
}

