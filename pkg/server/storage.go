package server

import (
	"errors"
	"sync"
)

var streams = make(map[Token]Peer)
var mx sync.Mutex

func addPeer(peer Peer) error {

	mx.Lock()
	defer mx.Unlock()

	_, ok := streams[peer.Token]
	if ok {
		return errors.New("Conn Is exists")
	}

	streams[peer.Token] = peer
	return nil
}

func DeletePeer(token Token) {

	mx.Lock()
	defer mx.Unlock()

	peer, ok := streams[token]
	if !ok {
		return
	}

	peer.GrpcStream.Context().Done()
	delete(streams, token)
}