package main

import (
	"encoding/json"
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
	msgChan     chan Message
	storage     StorageCore
}

type Message struct {
	Msg    []byte
	Sender *Peer
}

func NewServer(cfg Config) *Server {
	return &Server{
		Cfg:         cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		msgChan:     make(chan Message),
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
		case msg := <-s.msgChan:
			err := s.handleMsg(msg)
			if err != nil {
				msg.Sender.Send([]byte(err.Error()))
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

func (s *Server) handleMsg(msg Message) error {
	cmd, err := parseCommand(string(msg.Msg))
	if err != nil {
		return err
	}
	storage := s.storage.GetStorage(0, true)
	resp, err := cmd.Execute(storage) // later we add ability to select database numbere
	if err != nil {
		return err
	}
	sendData, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	msg.Sender.Send(sendData)
	return nil
}
