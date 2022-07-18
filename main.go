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

		err = c.SendExecute("hangup", "USER_BUSY", true)
		if err != nil {
			log.Println(err)
		}

		ev, err := c.ReadEvent()
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
