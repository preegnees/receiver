package server

import (
	"io"

	protobuf "receiver/pkg/server/proto"
)

func (c *connPeer) Recv() {

	for {
		m, err := c.current.GrpcStream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}
		if err := c.fileSaver.Write(m); err != nil {
			c.MessageSender.ReportError(err)
		}
	}
}

func (s *rpcServer) Streaming(stream protobuf.StreamingService_StreamingServer) error {

	peer, err := s.AuthPeer(stream)
	if err != nil {
		return err
	}

	addPeer(peer)
	defer DeletePeer(peer.Token)

	newCli := connPeer{
		current:   peer,
		log:       s.log,
		fileSaver: s.FileSaver,
	}

	
	go newCli.Recv()
	

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-s.runningCtx.Done():
			return nil
		}
	}
}
