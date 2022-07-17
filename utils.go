package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func keyValueString2Map(s string, sep string, kv_sep string) map[string]string {
	m := make(map[string]string)
	tokens := strings.Split(s, sep)
	for _, token := range tokens {
		tks := strings.Split(token, kv_sep)
		if len(tks) != 2 {
			log.Println("Invalid length")
		}
		m[tks[0]] = tks[1]
	}

	return m
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
