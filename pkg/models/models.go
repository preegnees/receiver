package models

import (
	pb "receiver/pkg/proto"
)

// ---> Conf

// CnfServer. Конифгурация для сервера
type CnfServer struct {
	Addr                 string
	ShutdownTimeout      int
	CertPem              string
	KeyPem               string
	Debug                bool
	StreamStoragePathLog string
	ServerPathLog        string
}

// CnfClient. Конфигурация для клиента
type CnfClient struct {
	Addr              string
	RequestTimeout    int
	KeepaliveInterval int
	Reconnect         bool
	ReconnectTimeout  int
	IdChannel         string
	Name              string
	AllowedNames      string
}

// IServ. Интерфейс для сервера
type IServ interface {
	Run(CnfServer) error
	Stop()
}

// // ICli. Интерфейс для клиента
// type ICli interface {
// 	Run(CnfClient) error
// 	Stop()
// }

// ---> MyMemStorage

type Username string
type Token string

type Peer struct {
	Username   Username
	Token      Token
	GrpcStream pb.StreamingService_StreamingServer
}

/*
IMemStorage. Интерфейс для взаимодействия с базой данных, которая хронит стримы.
SavePeer: сохранение пира при подключении;
DeletePeer: удаление пира при отключении;
*/
type IStreamStorage interface {
	SavePeer(peer Peer)
	DeletePeer(token Token)
}
