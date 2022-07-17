package main

import (
	"fmt"
	"log"

	"github.com/beevik/etree"
)

func main() {
	fmt.Println("Starting server")
	ListenAndServe(":8084", handler)
	fmt.Println("Finished server")
}

func handler(c *Connection) {
	go func() {
		fmt.Println("new client:", c.RemoteAddr())

		c.SendCommand("connect")

		ev, err := c.ReadEvent()
		if err != nil {
			log.Println(err)
			return
		}

		var xml_url = ""

		fmt.Println("got ev: ", ev.Header["Event-Name"])

		goivr_config := ev.Get("Variable_goivr_config")
		fmt.Println("goivr_config:")
		fmt.Println(goivr_config)

		m := keyValueString2Map(goivr_config, ";", "=")
		if _, ok := m["xml_url"]; ok {
			xml_url = m["xml_url"]
		} else {
			log.Println("Could not resolve xml_url")
		}

		doc := etree.NewDocument()

		if xmlBytes, err := getXML(xml_url); err != nil {
			log.Printf("Failed to get XML: %v", err)
			log.Printf("Got error while getting XML: %s", err)
		} else {
			fmt.Println(string(xmlBytes))
			err := doc.ReadFromBytes(xmlBytes)
			if err != nil {
				log.Println("Failed to parse xml")
				return
			}
			fmt.Println(doc)
			root := doc.Root()
			for _, ch := range root.ChildElements() {
				fmt.Println(ch)
			}
		}

		c.SendCommand("linger 10")
		ev, err = c.ReadEvent()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("got ev: ", ev.Header["Event-Name"])

		c.SendCommand("myevents")
		ev, err = c.ReadEvent()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("got ev: ", ev.Header["Event-Name"])

		ev, err = c.Execute("hangup", "USER_BUSY", true)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("got ev: ", ev.Header["Event-Name"])

		for {
			ev, err = c.ReadEvent()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println("\nNew event")
			fmt.Println("got ev: ", ev.Header["Event-Name"])
		}
	}()
}
