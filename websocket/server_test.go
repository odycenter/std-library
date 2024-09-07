package websocket_test

import (
	"context"
	"fmt"
	"github.com/odycenter/std-library/websocket"
	"github.com/olahol/melody"
	"log"
	"sync"
	"testing"
)

type Session struct {
	m map[string]*melody.Session
	sync.RWMutex
}

func (s *Session) get(k string) (*melody.Session, bool) {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.m[k]
	return v, ok
}

func (s *Session) set(k string, session melody.Session) {
	s.Lock()
	s.m[k] = &session
	s.Unlock()
}

func (s *Session) del(k string) {
	s.Lock()
	delete(s.m, k)
	s.Unlock()
}

func TestWS(t *testing.T) {
	sessions := Session{
		m: make(map[string]*melody.Session),
	}
	opt := websocket.Option{}
	websocket.New(&opt)
	opt.WithOnConnect(func(session *melody.Session) {
		fmt.Println("connected")
		req := websocket.GetRequest(session)
		ID := req.URL.Query().Get("A")
		sessions.set(ID, *session)
	})
	opt.WithOnDisconnect(func(session *melody.Session) {
		req := websocket.GetRequest(session)
		ID := req.URL.Query().Get("A")
		sessions.del(ID)
		fmt.Println("disconnected")
	})
	opt.WithOnClose(func(session *melody.Session, i int, s string) error {
		fmt.Printf("closed:%d=>%s\n", i, s)
		return nil
	})
	opt.WithOnMessage(func(session *melody.Session, b []byte) {
		log.Printf("message:%s\n", `b`)
		err := websocket.Cli().Broadcast(b)
		if err != nil {
			log.Panicln(err)
		}
		b = append(b, []byte(" [BroadcastOthers]")...)
		err = websocket.Cli().BroadcastOthers(b, session)
		if err != nil {
			log.Panicln(err)
		}
		b = append(b, []byte(" [BroadcastFilter]")...)
		err = websocket.Cli().BroadcastFilter(b, func(session *melody.Session) bool {
			uri := websocket.GetURL(session).Query().Get("A")
			ok := uri == "1"
			if ok {
				log.Printf("[%s]send to [1]:%s\n", uri, b)
			}
			return ok
		})
		if err != nil {
			log.Panicln(err)
		}
		v, ok := sessions.get("-1")
		if ok {
			err := v.Write(append([]byte("confirmed receive msg:"), b...))
			if err != nil {
				return
			}
		}
		b = append(b, []byte(" [BroadcastMultiple]")...)
		s, ok := sessions.get("2")
		if ok {
			err = websocket.Cli().BroadcastMultiple(b, []*melody.Session{s})
		}
	})
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	err := websocket.Cli().ListenWithServer(ctx)
	if err != nil {
		log.Panic(err)
		return
	}
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := websocket.Option{}
		websocket.New(&opt)
	}
}
