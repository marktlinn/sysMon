package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

type Server struct {
	SubscriberMsgBuffer int
	Mux                 http.ServeMux
	subscriberMutex     sync.Mutex
	Subscribers         map[*Subscriber]struct{}
}

type Subscriber struct {
	msgs chan []byte
}

func NewServer() *Server {
	s := &Server{
		SubscriberMsgBuffer: 10,
		Subscribers:         make(map[*Subscriber]struct{}),
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get CWD setting to relative './' : %s\n", err)
		cwd = "./"
	}

	htmxDir := filepath.Join(cwd, "htmx")

	s.Mux.Handle("/", http.FileServer(http.Dir(htmxDir)))

	s.Mux.HandleFunc("/ws", s.subscriberHandler)
	return s
}

func (s *Server) addSubscriber(subscriber *Subscriber) {
	s.subscriberMutex.Lock()
	defer s.subscriberMutex.Unlock()
	s.Subscribers[subscriber] = struct{}{}
	log.Printf("subscriber %v added to Subscribers: %v", subscriber, s.Subscribers)
}

func (s *Server) removeSubscriber(subscriber *Subscriber) {
	s.subscriberMutex.Lock()
	defer s.subscriberMutex.Unlock()

	_, ok := s.Subscribers[subscriber]
	if ok {
		delete(s.Subscribers, subscriber)
		log.Printf("subscriber %v removed from Subscribers: %v", subscriber, s.Subscribers)
	} else {
		log.Printf("subscriber: %v, not found in Subscribers: %v\n", subscriber, s.Subscribers)
	}
}

func (s *Server) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var c *websocket.Conn
	subscriber := &Subscriber{
		msgs: make(chan []byte, s.SubscriberMsgBuffer),
	}

	s.addSubscriber(subscriber)

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return fmt.Errorf("failed to accept websocket connection: %w", err)
	}

	defer func() {
		if closeErr := c.CloseNow(); closeErr != nil {
			log.Printf("failed to close websocket: %s\n", err)
		}
	}()

	ctx = c.CloseRead(ctx)

	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			err := c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return fmt.Errorf("failed to sent message through websocket: %w", err)
			}
		case <-ctx.Done():
			s.removeSubscriber(subscriber)
			return fmt.Errorf("context done: %w", ctx.Err())
		}
	}
}

func (s *Server) subscriberHandler(w http.ResponseWriter, r *http.Request) {
	err := s.subscribe(r.Context(), w, r)
	if err != nil {
		log.Printf("failed to subscribe to server: %s\n", err)
		return
	}
}

func (s *Server) Broadcast(msg []byte) {
	s.subscriberMutex.Lock()
	for sub := range s.Subscribers {
		sub.msgs <- msg
	}
	s.subscriberMutex.Unlock()
}
