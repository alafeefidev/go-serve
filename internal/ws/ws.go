package ws

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type WS struct {
	conn         *websocket.Conn
	cancel       context.CancelFunc
	message      chan []byte
	connected    chan struct{}
	disconnected chan struct{}
	mu           sync.Mutex
}

func NewWS() *WS {
	return &WS{
		message:      make(chan []byte, 256),
		connected:    make(chan struct{}, 1),
		disconnected: make(chan struct{}, 1),
	}
}

func (w *WS) Run() {
	for msg := range w.message {
		w.mu.Lock()
		conn := w.conn
		w.mu.Unlock()

		if conn == nil {
			slog.Warn("no client connected, skipping message")
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := conn.Write(ctx, websocket.MessageText, msg)
		cancel()
		if err != nil {
			slog.Error("error writing to websocket", "error", err)
			// w.Disconnect()
		}

	}
}

func (w *WS) Send(msg []byte) {
	w.message <- msg
}

func (w *WS) Connected() <-chan struct{} {
	return w.connected
}

func (w *WS) Disconnected() <-chan struct{} {
	return w.disconnected
}

func (w *WS) Connect(conn *websocket.Conn, cancel context.CancelFunc) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.conn != nil {
		slog.Warn("new client connected, dropping previous connection")
		w.cancel()
	}
	w.conn = conn
	w.cancel = cancel

	select {
	case <-w.disconnected:
	default:
	}
	select {
	case w.connected <- struct{}{}:
	default:
	}
}

func (w *WS) Disconnect() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.conn != nil {
		w.cancel()
		w.conn = nil
		w.cancel = nil
	}

	select {
	case <-w.connected:
	default:
	}
	select {
	case w.disconnected <- struct{}{}:
	default:
	}
}

func (w *WS) Sender(msg []byte) {
	for {

	}
}

func StartServerWS(ws *WS) (string, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", fmt.Errorf("error creating tcp listener: %w", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			slog.Error("error accepting websocket connection", "error", err)
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		ws.Connect(conn, cancel)
		defer ws.Disconnect()

		for {
			defer ws.Disconnect()
			if _, _, err = conn.Read(ctx); err != nil {
				break
			}
		}
	})

	go func() {
		if err := http.Serve(listener, mux); err != nil {
			slog.Warn("websocekt server stopped", "error", err)
		}
	}()

	addr := fmt.Sprintf("ws://%s/ws", listener.Addr().String())
	addr = strings.ReplaceAll(addr, "[::]", "localhost")
	return addr, nil
}
