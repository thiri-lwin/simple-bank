package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	db "github.com/thiri-lwin/thiri-bank/db/sqlc"
	"github.com/thiri-lwin/thiri-bank/gapi"
	"github.com/thiri-lwin/thiri-bank/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:5432/thiri_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	store := db.NewStore(conn)
	//server := api.NewServer(store)
	//server.Start("0.0.0.0:8080")

	// server := api.NewMuxServer(store)
	// server.StartMuxServer("0.0.0.0:8080")
	go runGatewayServer(store)
	runGrpcServer(store)
}

func runGrpcServer(store *db.Store) {
	server := gapi.NewServer(store)
	grpcServer := grpc.NewServer()
	pb.RegisterThiriBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server: ", err)
	}
}

func runGatewayServer(store *db.Store) {
	server := gapi.NewServer(store)

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := pb.RegisterThiriBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", "0.0.0.0:8440")
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start http gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start http gateway server: ", err)
	}
}
