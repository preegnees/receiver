package server

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestAuth(t *testing.T) {

	mt := metadata.Pairs("username", "radmir", "token", "token1")
	ctx := metadata.NewIncomingContext(context.TODO(), mt)

	newMt, err := getMetadata(ctx)
	if err != nil {
		panic(err)
	}
	username := Username(newMt["username"][0])
	if username != "radmir" {
		panic("username != radmir")
	}

}