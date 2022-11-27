package server

import (
	"context"
	"errors"

	protobuf "receiver/pkg/server/proto"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var ErrGetMetadata = errors.New("Error Get Metadata")
var ErrMetadataNotFound = errors.New("Error Metadata Not Found")
var ErrTokenIsInvalid = errors.New("Error Token Is Invalid")

func (s *rpcServer) AuthPeer(stream protobuf.StreamingService_StreamingServer) (Peer, error) {

	md, err := getMetadata(stream.Context())
	if err != nil {
		return Peer{}, err
	}

	ip := getIpFromPeer(stream.Context())

	username := Username(md["username"][0])
	token := Token(md["token"][0]) 
	g := stream

	ok, err := s.Repository.CheckToken(username, token)
	if err != nil {
		s.MessageSender.ReportError(err)
		return Peer{}, err
	}
	if !ok {
		return Peer{}, ErrTokenIsInvalid
	}

	return Peer{
		Username:   username,
		Token:      token,
		GrpcStream: g,
		IP:         ip,
	}, nil
}

func getMetadata(ctx context.Context) (metadata.MD, error) {
	
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrGetMetadata
	}

	if err := validateMetadata(md); err != nil {
		return nil, err
	}
	return md, nil
}

func validateMetadata(md metadata.MD) error {
	
	if len(md["username"]) != 1 || len(md["token"]) != 1 {
		return ErrMetadataNotFound
	}
	return nil
}

func getIpFromPeer(ctx context.Context) string {
	
	peer, ok := peer.FromContext(ctx)
	var ip string
	if ok {
		ip = peer.Addr.String()
	}
	return ip
}