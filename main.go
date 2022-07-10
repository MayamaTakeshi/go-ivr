package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func KeyValueString2Map(s string, sep string, kv_sep string) map[string]string {
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
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func main() {
	ListenAndServe(":8084", handler)
}

func handler(c *Connection) {
	go func () {
		fmt.Println("new client:", c.RemoteAddr())

		c.SendCommand("connect")

		ev, err := c.ReadEvent()
		if err != nil {
			log.Println(err)
			return
		}

		var xml_url = ""

		fmt.Println("got ev:")
		fmt.Println(ev)

		goivr_config := ev.Get("Variable_goivr_config")
		fmt.Println("goivr_config:")
		fmt.Println(goivr_config)

		m := KeyValueString2Map(goivr_config, ";", "=")
		if _, ok := m["xml_url"]; ok {
			xml_url = m["xml_url"]
		} else {
			log.Println("Could not resolve xml_url")
		}

		if xmlBytes, err := getXML(xml_url); err != nil {
			log.Printf("Failed to get XML: %v", err)
			log.Println("Got error while getting XML: %s", err)
		} else {
			fmt.Println(string(xmlBytes))
		}

		c.SendCommand("linger 10")
		ev, err = c.ReadEvent()
		fmt.Println(ev)

		c.SendCommand("myevents")
		ev, err = c.ReadEvent()
		fmt.Println(ev)

		ev, err = c.Execute("hangup", "USER_BUSY", true)
		if err != nil {
			log.Println(err)
		}
		ev.PrettyPrint()
		for {
			ev, err = c.ReadEvent()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println("\nNew event")
			ev.PrettyPrint()
		}
	}()
}
