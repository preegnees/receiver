package common

import (
	"fmt"
	
	m "receiver/pkg/models"
)

func PeertoString(peer m.Peer) string {
	return fmt.Sprintf(
		"Username=%s, token=%s, GrpcStream=%v",
		peer.Username, peer.Token, peer.GrpcStream,
	)
}