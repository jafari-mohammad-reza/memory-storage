package main

import (
	"fmt"
	"log/slog"
	"net"
)

type Server struct {
	Cfg         Config
	ln          net.Listener
	peers       map[*Peer]bool
	addPeerChan chan *Peer
	quitChan    chan struct{}
	msgChan     chan []byte
	respChan    chan []byte
	storage     StorageCore
}

func NewServer(cfg Config) *Server {
	return &Server{
		Cfg:         cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		msgChan:     make(chan []byte),
		respChan:    make(chan []byte),
		storage:     NewMemoryStorageCore(),
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
				go func() {
					s.respChan <- []byte(err.Error())
				}()
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
	peer := NewPeer(conn, s.msgChan, s.respChan)
	s.addPeerChan <- peer
	go func() {
		if err := peer.readLoop(); err != nil {
			slog.Error("peer read loop err", "err", err)
		}
	}()
	go func() {
		if err := peer.respLoop(); err != nil {
			slog.Error("peer resp loop err", "err", err)
		}
	}()
}

func (s *Server) handleRawMsg(msg []byte) error {
	cmd, err := parseCommand(string(msg))
	if err != nil {
		return err
	}
	storage := s.storage.GetStorage(0, true)
	err = cmd.Execute(storage) // later we add ability to select database numbere
	if err != nil {
		return err
	}
	return nil
}
