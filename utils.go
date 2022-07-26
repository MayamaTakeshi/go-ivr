package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func keyValueString2Map(m map[string]string, s string, sep string, kv_sep string) *GoIvrHalt {
	tokens := strings.Split(s, sep)
	for _, token := range tokens {
		tks := strings.Split(token, kv_sep)
		if len(tks) != 2 {
			return &GoIvrHalt{Error, "keyValueString2Map invalid length"}
		}
		m[tks[0]] = tks[1]
	}

	return nil
}

func getXML(url string) ([]byte, *GoIvrHalt) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, &GoIvrHalt{Error, fmt.Sprintf("GET %s error: %v", url, err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, &GoIvrHalt{Error, fmt.Sprintf("GET %s status error: %v", url, resp.StatusCode)}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, &GoIvrHalt{Error, fmt.Sprintf("GET %s read body error: %v", url, err)}
	}

	return data, nil
}
