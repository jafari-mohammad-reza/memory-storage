package main

import (
	"log/slog"
	"net"
)

type Peer struct {
	conn    net.Conn
	msgChan chan Message
}

func NewPeer(conn net.Conn, msgChan chan Message) *Peer {
	return &Peer{
		conn,
		msgChan,
	}
}

func (p *Peer) readLoop() error {
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			slog.Error("peer read error", "err", err)
			return err
		}
		msgBuff := make([]byte, n)
		copy(msgBuff, buf[:n])
		msg := Message{
			Msg:    msgBuff,
			Sender: p,
		}
		p.msgChan <- msg
	}
}

func (p *Peer) Send(msg []byte) error {
	_, err := p.conn.Write(msg)
	if err != nil {
		slog.Error("peer failed to send message", "err", err)
		return err
	}
	return nil
}
