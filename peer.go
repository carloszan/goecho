package main

import (
	"io"
	"log/slog"
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

func (p *Peer) readLoop() error {
	slog.Info("new peer")
	for {
		buf := make([]byte, 1024)
		count, err := p.conn.Read(buf)

		if err == io.EOF {
			p.delCh <- p
			break
		}

		if err != nil {
			slog.Info("reading message error", "remoteAddr", p.conn.RemoteAddr())
			return err
		}

		message := string(buf[:count])

		p.msgCh <- Message{
			cmd:  message,
			peer: p,
		}
	}

	return nil
}
