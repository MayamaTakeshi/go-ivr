// Adapted from code by Alexandre Fiori

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

const bufferSize = 1024 << 6 // For the socket reader
const eventsBuffer = 16      // For the events channel (memory eater!)

const (
	ContentTypePLAIN = "text/plain"
	ContentTypeJSON  = "application/json"
	ContentTypeXML   = "application/xml"
)

// Connection is the event socket connection handler.
type Connection struct {
	conn       net.Conn
	reader     *bufio.Reader
	textreader *textproto.Reader
	err        chan error
	evt        chan *Event

	goivrConfig   map[string]string
	sessionParams map[string]string
	xmlVars       map[string]interface{}

	xml *etree.Element

	stack Stack
}

// newConnection allocates a new Connection and initialize its buffers.
func newConnection(conn net.Conn) *Connection {
	c := Connection{
		conn:   conn,
		reader: bufio.NewReaderSize(conn, bufferSize),
		err:    make(chan error, 1),
		evt:    make(chan *Event, eventsBuffer),

		goivrConfig:   make(map[string]string),
		sessionParams: make(map[string]string),
		xmlVars:       make(map[string]interface{}),
	}
	c.textreader = textproto.NewReader(c.reader)
	return &c
}

// HandleFunc is the function called on new incoming connections.
type HandleFunc func(*Connection)

// ListenAndServe listens for incoming connections from FreeSWITCH and calls
// HandleFunc in a new goroutine for each client.
//
// Example:
//
//	func main() {
//		eventsocket.ListenAndServe(":9090", handler)
//	}
//
//	func handler(c *eventsocket.Connection) {
//		ev, err := c.Command("connect") // must always start with this
//		ev.PrettyPrint()             // print event to the console
//		...
//		c.Command("myevents")
//		for {
//			ev, err = c.ReadEvent()
//			...
//		}
//	}
//
func ListenAndServe(addr string, fn HandleFunc) error {
	srv, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := srv.Accept()
		if err != nil {
			return err
		}
		c := newConnection(conn)
		go c.readLoop()
		go fn(c)
	}
}

// readLoop calls readOne until a fatal error occurs, then close the socket.
func (c *Connection) readLoop() {
	for c.readOne() {
	}
	c.Close()
}

// readOne reads a single event and send over the appropriate channel.
// It separates incoming events from api and command responses.
func (c *Connection) readOne() bool {
	hdr, err := c.textreader.ReadMIMEHeader()
	if err != nil {
		c.err <- err
		return false
	}
	resp := new(Event)
	resp.Header = make(EventHeader)
	if v := hdr.Get("Content-Length"); v != "" {
		length, err := strconv.Atoi(v)
		if err != nil {
			c.err <- errors.New("failed to parse Content-Length header")
			return false
		}
		b := make([]byte, length)
		if _, err := io.ReadFull(c.reader, b); err != nil {
			c.err <- errors.New("failed to read data from ESL socket")
			return false
		}
		resp.Body = string(b)
	}
	switch hdr.Get("Content-Type") {
	case "command/reply":
		reply := hdr.Get("Reply-Text")

		if reply[0] == '%' {
			copyHeaders(&hdr, resp, true)
		} else {
			copyHeaders(&hdr, resp, false)
		}
		c.evt <- resp
	case "api/response":
		copyHeaders(&hdr, resp, false)
		fmt.Println("adding api/response to c.evt")
		c.evt <- resp
	case "text/event-plain":
		reader := bufio.NewReader(bytes.NewReader([]byte(resp.Body)))
		resp.Body = ""
		textreader := textproto.NewReader(reader)
		hdr, err = textreader.ReadMIMEHeader()
		if err != nil {
			c.err <- err
			return false
		}
		if v := hdr.Get("Content-Length"); v != "" {
			length, err := strconv.Atoi(v)
			if err != nil {
				c.err <- errors.New("failed to parse Content-Length header for body")
				return false
			}
			b := make([]byte, length)
			if _, err = io.ReadFull(reader, b); err != nil {
				c.err <- errors.New("failed to read body data from ESL socket")
				return false
			}
			resp.Body = string(b)
		}
		copyHeaders(&hdr, resp, true)
		c.evt <- resp
	case "text/disconnect-notice":
		copyHeaders(&hdr, resp, false)
		c.evt <- resp
	default:
		log.Fatal("Unsupported event:", hdr)
	}
	return true
}

// RemoteAddr returns the remote addr of the connection.
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// Close terminates the connection.
func (c *Connection) Close() {
	c.conn.Close()
}

// ReadEvent reads and returns events from the server (plain format only)
func (c *Connection) ReadEvent() (*Event, error) {
	var (
		ev  *Event
		err error
	)
	select {
	case err = <-c.err:
		return nil, err
	case ev = <-c.evt:
		return ev, nil
	}
}

// copyHeaders copies all keys and values from the MIMEHeader to Event.Header,
// normalizing header keys to their capitalized version and values by
// unescaping them when decode is set to true.
//
// It's used after parsing plain text event headers, but not JSON.
func copyHeaders(src *textproto.MIMEHeader, dst *Event, decode bool) {
	var err error
	for k, v := range *src {
		k = capitalize(k)
		if decode {
			dst.Header[k], err = url.QueryUnescape(v[0])
			if err != nil {
				dst.Header[k] = v[0]
			}
		} else {
			dst.Header[k] = v[0]
		}
	}
}

// capitalize capitalizes strings in a very particular manner.
// Headers such as Job-UUID become Job-Uuid and so on. Headers starting with
// Variable_ only replace ^v with V, and headers staring with _ are ignored.
func capitalize(s string) string {
	if s[0] == '_' {
		return s
	}
	ns := bytes.ToLower([]byte(s))
	if len(s) > 9 && s[1:9] == "ariable_" {
		ns[0] = 'V'
		return string(ns)
	}
	toUpper := true
	for n, c := range ns {
		if toUpper {
			if 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
			}
			ns[n] = c
			toUpper = false
		} else if c == '-' || c == '_' {
			toUpper = true
		}
	}
	return string(ns)
}

//----------------------------------------------------
//--------- Event
//----------------------------------------------------

// EventHeader represents events as a pair of key:value.
type EventHeader map[string]string

// Event represents a FreeSWITCH event.
type Event struct {
	Header EventHeader // Event headers, key:val
	Body   string      // Raw body, available in some events
}

func (r *Event) String() string {
	if r.Body == "" {
		return fmt.Sprintf("%s", r.Header)
	} else {
		return fmt.Sprintf("%s body=%s", r.Header, r.Body)
	}
}

// Get returns an Event value, or "" if the key doesn't exist.
func (r *Event) Get(key string) string {
	val, ok := r.Header[key]
	if !ok {
		return ""
	}
	return val
}

// PrettyPrint prints Event headers and body to the standard output.
func (r *Event) PrettyPrint() {
	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %#v\n", k, r.Header[k])
	}
	if r.Body != "" {
		fmt.Printf("BODY: %#v\n", r.Body)
	}
}

// Pretty returns Event headers and body with pretty format.
func (r *Event) Pretty() string {
	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var items []string
	for _, k := range keys {
		items = append(items, fmt.Sprintf("%s: %#v\n", k, r.Header[k]))
	}
	if r.Body != "" {
		items = append(items, "\n\n")
		items = append(items, fmt.Sprintf("BODY: %#v\n", r.Body))
	}
	return strings.Join(items, "")
}

func (c *Connection) SendCommand(command string) error {
	// Sanity check to avoid breaking the parser
	//if strings.IndexAny(command, "\r\n") > 0 {
	//	return nil, errInvalidCommand
	//}
	_, err := fmt.Fprintf(c.conn, "%s\r\n\r\n", command)
	return err
}

// MSG is the container used by SendMsg to store messages sent to FreeSWITCc.
// It's supposed to be populated with directives supported by the sendmsg
// command only, like "call-command: execute".
//
// See https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg for details.
type MSG map[string]string

// SendMsg sends messages to FreeSWITCH and returns a response Event.
//
// Examples:
//
//	SendMsg(MSG{
//		"call-command": "hangup",
//		"hangup-cause": "we're done!",
//	}, "", "")
//
//	SendMsg(MSG{
//		"call-command":     "execute",
//		"execute-app-name": "playback",
//		"execute-app-arg":  "/tmp/test.wav",
//	}, "", "")
//
// Keys with empty values are ignored; uuid and appData are optional.
// If appData is set, a "content-length" header is expected (lower case!).
//
// See https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg for details.
func (c *Connection) SendSendMsg(m MSG, appData string) error {
	b := bytes.NewBufferString("sendmsg")
	b.WriteString("\n")
	for k, v := range m {
		if v != "" {
			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}
	b.WriteString("\n")
	if m["content-length"] != "" && appData != "" {
		b.WriteString(appData)
	}
	if _, err := b.WriteTo(c.conn); err != nil {
		return err
	}

	return nil
}

// Execute is a shortcut to SendMsg with call-command: execute without UUID,
// suitable for use on outbound event socket connections (acting as server).
//
// Example:
//
//	Execute("playback", "/tmp/test.wav", false)
//
// See https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg for details.
func (c *Connection) SendExecute(appName, appArg string, lock bool) error {
	var evlock string
	if lock {
		// Could be strconv.FormatBool(lock), but we don't want to
		// send event-lock when it's set to false.
		evlock = "true"
	}
	return c.SendSendMsg(MSG{
		"call-command":     "execute",
		"execute-app-name": appName,
		"execute-app-arg":  appArg,
		"event-lock":       evlock,
	}, "")
}

func (c *Connection) IsTerminated() bool {
	for {
		select {
		case err := <-c.err:
			log.Println(err)
			return true
		case ev := <-c.evt:
			if ev.Header["Content-Type"] == "text/disconnect-notice" {
				return true
			}
			// discard any other messages
		default:
			return false
		}
	}
}

func (c *Connection) SendCommandAndWaitOK(command string) error {
	err := c.SendCommand(command)
	if err != nil {
		return err
	}

	for {
		select {
		case err := <-c.err:
			log.Println(err)
			return err
		case ev := <-c.evt:
			switch ev.Header["Content-Type"] {
			case "text/disconnect-notice":
				return errors.New("connection terminated")
			case "command/reply":
				reply := ev.Header["Reply-Text"]
				if reply[:2] == "-E" {
					return errors.New(reply[5:])
				} else {
					return nil
				}
				// ignore other messages (just consume them)
			}
		}
	}
}

func (c *Connection) initialize() error {
	c.SendCommand("connect")
	ev, err := c.ReadEvent()

	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("got ev: ", ev.Header["Event-Name"])

	goivr_config := ev.Get("Variable_goivr_config")
	fmt.Println("goivr_config:")
	fmt.Println(goivr_config)

	if goivr_config == "" {
		err := errors.New("variable_goivr_config absent")
		return err
	}

	err = keyValueString2Map(c.goivrConfig, goivr_config, ";", "=")
	if err != nil {
		return err
	}

	err = c.SendCommandAndWaitOK("linger 10")
	if err != nil {
		return err
	}

	err = c.SendCommandAndWaitOK("myevents")
	if err != nil {
		return err
	}

	if xml_url, ok := c.goivrConfig["xml_url"]; ok {
		doc := etree.NewDocument()

		if xmlBytes, err := getXML(xml_url); err != nil {
			err := fmt.Errorf("got error while getting XML: %w", err)

			return err
		} else {
			err := doc.ReadFromBytes(xmlBytes)
			if err != nil {
				err := errors.New("failed to parse xml")
				return err
			}
			fmt.Println(doc)
			root := doc.Root()
			for _, ch := range root.ChildElements() {
				fmt.Println(ch)
			}

			c.xml = root
		}
	} else {
		err := errors.New("could not resolve xml_url")
		return err
	}

	return nil
}

func (c *Connection) initializeStack() error {
	elems := c.xml.FindElements("//Section[name='main']")
	if len(elems) == 0 {
		temp := make([]*etree.Element, 0)
		copy(c.xml.ChildElements(), temp)
		c.stack.Push(temp)
	} else if len(elems) == 1 {
		c.stack.Push(elems[0].ChildElements())
	} else {
		return errors.New("more than one Section with @name=main")
	}
	return nil
}

func (c *Connection) process() error {
	var err error

	for !c.stack.IsEmpty() {
		var head *etree.Element
		elems := c.stack.Top()
		for len(elems) > 0 {
			head, elems = elems[0], elems[1:]
			err = c.processElement(head)
			if err != nil {
				break
			}

			if err != nil || c.IsTerminated() {
				break
			}
		}

		if err != nil || c.IsTerminated() {
			break
		}
	}
	return nil
}

func (c *Connection) processElement(elem *etree.Element) error {
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
}
