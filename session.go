package glive

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"github.com/nohns/go-liveview/internal/rng"
)

type Session struct {
	id     string
	ctx    context.Context
	mounts map[string]View
	views  map[string]View

	term      chan struct{}
	rerenders chan string // chan of view ids
}

func StartSession(ctx context.Context) (*Session, error) {
	id, err := rng.SessionID()
	if err != nil {
		return nil, err
	}
	ses := &Session{
		id:     id,
		ctx:    ctx,
		views:  make(map[string]View),
		mounts: make(map[string]View),
	}
	return ses, nil
}

func (s *Session) IsServing() bool {
	return s.term != nil
}

type hydration struct {
	ViewID string     `json:"vid"`
	Diff   [][]string `json:"h"`
}

func hydrateView(sock *websocket.Conn, v View) error {
	var b templ.DiffBuffer
	if err := v.Render().Render(context.TODO(), &b); err != nil {
		return fmt.Errorf("hydration render: %v", err)
	}
	b.Flush()

	if err := sock.WriteJSON(map[string]any{"type": "hydration", "data": hydration{ViewID: v.ID(), Diff: [][]string{b.Segs, b.Vals}}}); err != nil {
		return fmt.Errorf("write hydration json: %v", err)
	}
	return nil
}

type diff struct {
	ViewID string   `json:"vid"`
	Vals   []string `json:"v"`
}

func diffView(sock *websocket.Conn, v View) error {
	var b templ.DiffBuffer
	if err := v.Render().Render(context.TODO(), &b); err != nil {
		return fmt.Errorf("diff render: %v", err)
	}

	if err := sock.WriteJSON(map[string]any{"type": "diff", "data": diff{ViewID: v.ID(), Vals: b.Vals}}); err != nil {
		return fmt.Errorf("write diff json: %v", err)
	}
	return nil
}

func (s *Session) Serve(sock *websocket.Conn) error {
	defer sock.Close()
	s.term = make(chan struct{})
	s.rerenders = make(chan string)

	// Hydrate all live views on page
	for _, v := range s.views {
		if err := hydrateView(sock, v); err != nil {
			return err
		}
	}

	// Read incoming messages on chan
	msgch := make(chan []byte)
	go func() {
		for {
			_, p, err := sock.ReadMessage()
			if err != nil {
				// TODO: EOF will probably block forever
				s.term <- struct{}{} // terminate on error
				return
			}
			msgch <- p
		}
	}()

	for {
		select {
		case p := <-msgch:
			// handle incomming msg
			fmt.Printf("recv msg len %d\n", len(p))
		case id := <-s.rerenders:
			if err := diffView(sock, s.views[id]); err != nil {
				return err
			}
		case <-s.term:
			for _, v := range s.views {
				h, ok := v.(ViewLifecycleUnmount)
				if !ok {
					continue
				}
				h.OnUnmount()
			}
			return fmt.Errorf("session terminated")
		}
	}
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Context() context.Context {
	return s.ctx
}

func (s *Session) BodyAttr() templ.Attributes {
	return templ.Attributes{
		"data-glive-session": s.id,
	}
}

func (s *Session) Mount(name string, builder ViewBuilder) templ.Component {
	v, ok := s.mounts[name]
	if !ok {
		v = builder()
		s.mounts[name] = v
		s.views[v.ID()] = v

		h, ok := v.(ViewLifecycleMount)
		if ok {
			h.OnMount(s)
		}
	}
	return v.Render()
}

func (s *Session) Rerender(v View) error {
	s.rerenders <- v.ID()
	return nil
}

func (s *Session) Close() error {
	if s.term == nil {
		return nil
	}
	s.term <- struct{}{}
	return nil
}
