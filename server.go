package main

import (
	"fmt"
	"log/slog"
	"net"
)

type Peer struct {
	conn    net.Conn
	msgChan chan []byte
}

func NewPeer(conn net.Conn, msgChan chan []byte) *Peer {
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
		p.msgChan <- msgBuff
	}
}

type Server struct {
	Cfg         Config
	ln          net.Listener
	peers       map[*Peer]bool
	addPeerChan chan *Peer
	quitChan    chan struct{}
	msgChan     chan []byte
}

func NewServer(cfg Config) *Server {
	return &Server{
		Cfg:         cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		msgChan:     make(chan []byte),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Cfg.ServerListenAddr))
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()
	slog.Info("server running", "listenAddr", s.Cfg.ServerListenAddr)
	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgChan:
			err := s.handleRawMsg(rawMsg)
			if err != nil {
				slog.Error("handle message err", "err", err.Error())
			}
		case peer := <-s.addPeerChan:
			s.peers[peer] = true
		case <-s.quitChan:
			return

		}
	}
}
func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgChan)
	s.addPeerChan <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read loop err", "err", err)
	}
}

func (s *Server) handleRawMsg(msg []byte) error {
	return nil
}
