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

		err = validate(c.xml, 0)
		if err != nil {
			log.Println(err)
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
		
	}()
}
