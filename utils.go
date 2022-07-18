package main

import (
	"fmt"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func keyValueString2Map(m map[string]string, s string, sep string, kv_sep string) error {
	tokens := strings.Split(s, sep)
	for _, token := range tokens {
		tks := strings.Split(token, kv_sep)
		if len(tks) != 2 {
			return errors.New("keyValueString2Map invalid length")
		}
		m[tks[0]] = tks[1]
	}

	return nil
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("read body: %v", err)
	}

	return data, nil
}
