package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	protobuf "receiver/pkg/server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (s *rpcServer) Run() error {

	go s.start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	select {
	case <-sigCh:
		s.stop()
		return nil
	case err := <-s.errCh:
		s.stop()
		return err
	}
}

func (s *rpcServer) Stop() {
	s.gracefulStop()
}

func (s *rpcServer) start() {

	for {
		s.runningCtx, s.runningCancel = context.WithCancel(context.Background())

		creds, err := credentials.NewServerTLSFromFile(s.CertPem, s.KeyPem)
		if err != nil {
			s.errCh <- fmt.Errorf("$Ошибка при добавлении TLS, err:=%v", err)
		}

		s.grpcServer = grpc.NewServer(
			grpc.Creds(creds),
		)

		protobuf.RegisterStreamingServiceServer(s.grpcServer, s)

		go func() {
			listener, err := net.Listen("tcp", s.Addr)
			if err != nil {
				s.runningCancel()
				s.errCh <- fmt.Errorf("$Ошибка при прослушивании сети, err:=%v", err)
			}
			s.running = true
			s.grpcServer.Serve(listener)
		}()

		select {
		case <-s.stopCh:
			s.gracefulStop()
			s.running = false
			s.doneCh <- struct{}{}
			return
		}
	}
}

func (s *rpcServer) gracefulStop() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.ShutdownTimeout)*time.Second)
	defer cancel()

	s.runningCancel()

	closed := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		if closed != nil {
			select {
			case closed <- struct{}{}:
			default:
			}
		}
	}()

	select {
	case <-closed:
		close(closed)
	case <-ctx.Done():
		s.grpcServer.Stop()
		close(closed)
		closed = nil
	}
}

func (s *rpcServer) stop() {

	if s.running {
		s.stopCh <- struct{}{}
		<-s.doneCh
	}
}