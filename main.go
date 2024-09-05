package main

import (
	"flag"
	"fmt"
	"github.com/AlessandroFinocchi/sdcc_common/pb"
	u "github.com/AlessandroFinocchi/sdcc_common/utils"
	"log"
	"net"
	"os"
	m "sdcc_registry/model"
	s "sdcc_registry/services"
	ur "sdcc_registry/utils"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	// value represent the default port number if "-port" flag is not provided
	port = flag.Int("port", 50051, "Server port")
)

func main() {
	flag.Parse()

	fmt.Println("Current process PID: ", os.Getgid())

	nodeListW := m.NewNodeListWrapper()
	nodeListMutex := &sync.Mutex{}
	timeoutDuration, err := u.ReadConfigUInt64("config.ini", "heartbeat", "timeout_duration")
	if err != nil {
		log.Fatalf("failed to read timeout duration: %v", err)
	}

	connectorService := s.NewConnector(nodeListMutex, nodeListW)
	heartbeatService := s.NewHeartbeat(nodeListMutex, nodeListW)

	// Get the port
	serverAddress := fmt.Sprintf(":%d", *port)

	// Create TCP listener on the port; gRPC server will bind to it
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tlsCredentials, err := ur.LoadServerTLSCredentials()
	if err != nil {
		log.Fatalf("cannot load TLS credentials: %w")
	}
	serverOptions := []grpc.ServerOption{grpc.Creds(tlsCredentials)}

	// Create new gRPC server instance by calling gRPC Go APIs
	registry := grpc.NewServer(serverOptions...)

	pb.RegisterConnectorServer(registry, connectorService)
	pb.RegisterHeartbeatServer(registry, heartbeatService)
	fmt.Printf("server listening at %v", lis.Addr())

	// Start listening to the incoming messages on the port
	go func() {
		err := registry.Serve(lis)
		if err != nil {
			log.Fatalf("failed to call Serve API")
		}
	}()

	ticker := time.NewTicker(time.Duration(timeoutDuration) * time.Second)
	for range ticker.C {
		heartbeatService.OnTimeout()
	}

	return
}
