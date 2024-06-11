package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	glive "github.com/nohns/go-liveview"
	"github.com/nohns/go-liveview/cmd/view"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const _sessionIdCookie = "glive_session_id"

func main() {
	s := server{sessions: make(map[string]*glive.Session)}
	s.serve()
}

type server struct {
	sessions map[string]*glive.Session
}

func (srv *server) serve() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fp := filepath.Join(cwd, "assets")
	fmt.Printf("assets: %v\n", fp)
	// Serve assets which will contain ws setup
	fs := http.FileServer(http.Dir(fp))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))

	http.HandleFunc("/ws/{sid}", func(w http.ResponseWriter, r *http.Request) {
		sid := r.PathValue("sid")
		if sid == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ses, ok := srv.sessions[sid]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("error occured when establishing ws conn: %v\n", err)
			return
		}

		go func(sid string) {
			if err := ses.Serve(conn); err != nil {
				fmt.Printf("error occured when serving session: %v", err)
			}
			delete(srv.sessions, sid)
		}(ses.ID())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ses, err := glive.StartSession(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		srv.sessions[ses.ID()] = ses

		// Automatically cleanup session, if not serving within a given time frame
		time.AfterFunc(10*time.Second, func() {
			if ses.IsServing() {
				return
			}
			ses.Close()
			delete(srv.sessions, ses.ID())
		})

		view.HomePage(ses).Render(r.Context(), w)
	})

	// serve blocking
	fmt.Println("serving http on port 7070...")
	if err := http.ListenAndServe(":7070", nil); err != nil {
		panic(err)
	}
}

func (srv *server) resolveSession(w http.ResponseWriter, r *http.Request) (*glive.Session, bool) {
	var ses *glive.Session
	c, err := r.Cookie(_sessionIdCookie)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusBadRequest)
		return nil, false
	}

	// no cookie found or the session found in cookie does not exist on server
	if errors.Is(err, http.ErrNoCookie) || srv.sessions[c.Value] == nil {
		ses, err = glive.StartSession(context.TODO())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil, false
		}
		http.SetCookie(w, &http.Cookie{
			Name:     _sessionIdCookie,
			Value:    ses.ID(),
			HttpOnly: true,
		})
		srv.sessions[ses.ID()] = ses
		return ses, true
	}
	return srv.sessions[c.Value], true
}
