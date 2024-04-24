package main

import (
	"log/slog"
	"net"
)

type Message struct {
	cmd  string
	peer *Peer
}

type Server struct {
	ListenAddr string
	ln         net.Listener
	msgCh      chan Message
	peers      map[*Peer]bool
	addPeerCh  chan *Peer
	delPeerCh  chan *Peer
}

func NewServer(listenAddr string) *Server {
	return &Server{
		ListenAddr: listenAddr,
		msgCh:      make(chan Message),
		peers:      make(map[*Peer]bool),
		addPeerCh:  make(chan *Peer),
		delPeerCh:  make(chan *Peer),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}

	s.ln = ln

	go s.loop()

	slog.Info("goecho server is running", "ListenAddr", s.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) loop() error {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw message error", "err", err)
				return err
			}
		case peer := <-s.addPeerCh:
			slog.Info("peer connected", "remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true
		case peer := <-s.delPeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
			defer peer.conn.Close()
		}
	}
}

func (s *Server) handleMessage(msg Message) error {
	peer := msg.peer
	_, err := peer.conn.Write([]byte(msg.cmd))

	if err != nil {
		slog.Info("error sending message back", "remoteAddr", peer.conn.RemoteAddr())
		return err
	}

	slog.Info("sending message back: ", "message", msg.cmd)

	return nil
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh, s.delPeerCh)
	s.addPeerCh <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}
