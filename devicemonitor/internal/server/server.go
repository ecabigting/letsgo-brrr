package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/ecabigting/letsgo-brrr/devicemonitor/internal/hardware"
)

type Server struct {
	Port string
	// to keep track of Subscribers
	// we create a map of subscribers
	// as anonymous struct
	Subscribers             map[*Subscriber]struct{} // anonymous struct
	Mux                     http.ServeMux
	SubscriberMutex         sync.Mutex
	SubscriberMessageBuffer int
}

type Subscriber struct {
	// channel where we
	// can send messages
	// for the Subscriber
	msgs chan []byte
}

func NewServer(sub int, port string) *Server {
	s := &Server{
		Port:                    port,
		SubscriberMessageBuffer: sub,
		Subscribers:             make(map[*Subscriber]struct{}),
	}

	s.Mux.Handle("/", http.FileServer(http.Dir("./htmx")))
	s.Mux.HandleFunc("/ws", s.SubscribeHandler)
	return s
}

func (s *Server) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.Subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Server) AddSubscriber(subscriber *Subscriber) {
	// lets lock the subscriber via mutex
	// before we do some changes
	s.SubscriberMutex.Lock()
	s.Subscribers[subscriber] = struct{}{}
	// unlock the subscriber so
	// other resources can update it
	s.SubscriberMutex.Unlock()
	fmt.Println("Added subscriber", subscriber)
}

// Subscribe by creating connections
// that we can keep track and able to send
// messages to that connection
func (s *Server) Subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// creating the connection
	// via websockets using
	// the Subscriber struct
	var c *websocket.Conn
	subscriber := &Subscriber{
		msgs: make(chan []byte, s.SubscriberMessageBuffer),
	}

	// lets add
	// the new connection
	// to our server
	s.AddSubscriber(subscriber)

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}

	// close the connection
	// after the execution
	defer c.CloseNow()

	// close using
	// the implemention of the
	// io closer for this
	// context
	ctx = c.CloseRead(ctx)

	// lets loop thru
	// the websockets connected
	// whenever a new message is
	// available
	for {
		select {
		// lets write the new message
		// on this new channel(msgs)
		// in this new subscriber
		// we just created
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			err := c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// broadcast all the messages
// to all connected subscribers
// via channels
// note the usage of Mutex
// to protect the thread
func (s *Server) Broadcast(msg []byte) {
	s.SubscriberMutex.Lock()
	for subs := range s.Subscribers {
		subs.msgs <- msg
	}
	s.SubscriberMutex.Unlock()
}

func (s *Server) GetHardwareData() (string, string, string) {
	sectionSys, err := hardware.GetSystemSection()
	if err != nil {
		fmt.Println(err)
	}

	sectionDisc, err := hardware.GetDiskSection()
	if err != nil {
		fmt.Println(err)
	}

	sectionCPU, err := hardware.GetCpuSection()
	if err != nil {
		fmt.Println(err)
	}
	return sectionSys, sectionCPU, sectionDisc
}
