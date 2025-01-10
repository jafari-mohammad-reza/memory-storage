package main

import (
	"log/slog"
	"net"
)

type Peer struct {
	conn     net.Conn
	msgChan  chan []byte
	respChan chan []byte
}

func NewPeer(conn net.Conn, msgChan chan []byte, respChan chan []byte) *Peer {
	return &Peer{
		conn,
		msgChan,
		respChan,
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
		p.msgChan <- msgBuff
	}
}

func (p *Peer) respLoop() error {
	for resp := range p.respChan {
		_, err := p.conn.Write(resp)
		if err != nil {
			slog.Error("writing response error", "err", err)
			return err
		}
	}
	return nil
}
