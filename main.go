package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting server")
	ListenAndServe(":8084", handler)
	fmt.Println("Finished server")
}

func handler(c *Connection) {
	go func() {
		fmt.Println("new client:", c.RemoteAddr())

		err := c.initialize()
		if err != nil {
			log.Println(err)
			c.Close()
			return
		}

		err2 := validateXML(c.xml, 0)
		if err2 != nil {
			log.Println(err2)
			c.Close()
			return
		}

		err = c.initializeStack()
		if err != nil {
			log.Println(err)
			c.Close()
			return
		}

		err = c.process()
		if err != nil {
			log.Println(err)
			c.Close()
			return
		}
		
		fmt.Println("Finished")
	}()
}
