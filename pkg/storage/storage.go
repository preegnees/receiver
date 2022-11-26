package storage

import (
	"sync"

	m "receiver/pkg/models"
	

	l "github.com/sirupsen/logrus"
)


type storage struct {
	streams map[m.Token]m.Peer
	mx sync.Mutex
	log *l.Logger
}

func NewStorage(logger *l.Logger) *storage {

	streams := make(map[m.Token]m.Peer)
	return &storage{
		streams: streams,
		log:     logger,
	}
}

func (s *storage) SavePeer(peer m.Peer) {

	s.mx.Lock()
	defer s.mx.Unlock()

	s.streams[peer.Token] = peer
}

func (s *storage) DeletePeer(token m.Token) {

	s.mx.Lock()
	defer s.mx.Unlock()

	peer, ok := s.streams[token]
	if !ok {
		return 
	}

	peer.GrpcStream.Context().Done()
	delete(s.streams, token)
	return
}