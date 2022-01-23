package main
 
import (
    . "github.com/0x19/goesl"
	"github.com/sbabiv/xml2map"
    "strings"
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
)

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
 
var (
    goeslMessage = "Hello from GoESL. Open source FreeSWITCH event socket wrapper written in Go!"
)
 
func main() {
    defer func() {
        if r := recover(); r != nil {
            Error("Recovered in: ", r)
        }
    }()
 
    if s, err := NewOutboundServer("127.0.0.1:8084"); err != nil {
        Error("Got error while starting FreeSWITCH outbound server: %s", err)
    } else {
        go handle(s)
        s.Start()
    }
 
}
 
// handle - Running under goroutine here to explain how to run tts outbound server
func handle(s *OutboundServer) {
    for {
 
        select {
 
        case conn := <-s.Conns:
            Notice("New incomming connection: %v", conn)
 
            if err := conn.Connect(); err != nil {
                Error("Got error while accepting connection: %s", err)
                break
            }

			var xml_url = "http://127.0.0.1/ivr"

			if xmlBytes, err := getXML(xml_url); err != nil {
			  log.Printf("Failed to get XML: %v", err)
              Error("Got error while getting XML: %s", err)
			} else {
				fmt.Println(xmlBytes)	
			} 
			
            //if hm, err := conn.ExecuteHangup(cUUID, "USER_BUSY", false); err != nil {
            if hm, err := conn.ExecuteHangup("", "USER_BUSY", false); err != nil {
                Error("Got error while executing hangup: %s", err)
                break
            } else {
                Debug("Hangup Message: %s", hm)
            }
 
            go func() {
                for {
                    msg, err := conn.ReadMessage()
 
                    if err != nil {
 
                        // If it contains EOF, we really dont care...
                        if !strings.Contains(err.Error(), "EOF") {
                            Error("Error while reading Freeswitch message: %s", err)
                        }
                        break
                    }
 
                    Debug("Got message: %s", msg)
                }
            }()
 
        default:
        }
    }
 
}
