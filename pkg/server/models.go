package server

import (
	"context"

	protobuf "receiver/pkg/server/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//go:generate mockgen -source=models.go -destination=mock/fileSaverMock.go 
type IFileSaver interface {
	Write(*protobuf.Message) error
}

//go:generate mockgen -source=models.go -destination=mock/repositoryMock.go
type IRepository interface {
	CheckToken(Username, Token) (bool, error)
}

//go:generate mockgen -source=models.go -destination=mock/messageSenderMock.go
type IMessageSender interface {
	ReportState()
	ReportError(error)
}

type Username string
type Token string

type Peer struct {
	Username   Username
	Token      Token
	GrpcStream protobuf.StreamingService_StreamingServer
	IP         string
}

type CnfServer struct {
	Addr                 string
	ShutdownTimeout      int
	CertPem              string
	KeyPem               string
	Debug                bool
	StreamStoragePathLog string
	ServerPathLog        string
	FileSaver            IFileSaver
	Repository           IRepository
	MessageSender        IMessageSender
}

type rpcServer struct {
	protobuf.StreamingServiceServer
	CnfServer

	stopCh        chan struct{}
	restartCh     chan struct{}
	doneCh        chan struct{}
	errCh         chan error
	runningCtx    context.Context
	runningCancel context.CancelFunc
	running       bool
	grpcServer    *grpc.Server
	log           *logrus.Logger
}

type connPeer struct {
	current       Peer
	log           *logrus.Logger
	fileSaver     IFileSaver
	MessageSender IMessageSender
}

func New(cnf CnfServer) rpcServer {

	return rpcServer{
		stopCh:    make(chan struct{}),
		restartCh: make(chan struct{}),
		doneCh:    make(chan struct{}),
		errCh:     make(chan error),
		CnfServer: cnf,
		running:   false,
		log:       logrus.New(),
	}
}
