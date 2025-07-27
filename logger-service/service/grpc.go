package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

// LogServer represents the gRPC server for logging service
type LogServer struct {
	logs.UnimplementedLogServiceServer             // Embedding the unimplemented server to satisfy the interface
	Models                             data.Models // Dependency injection for data models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry() // Extract the log entry from the request

	// Create a log entry from the request data
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// Insert the log entry into the data store
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		// Return a failure response if insertion fails
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// Return a success response if insertion succeeds
	res := &logs.LogResponse{Result: "logged!"}
	return res, nil

}

// gRPCListen starts the gRPC server and listens for incoming requests
func (app *Config) gRPCListen() {
	// Listen on the specified port for gRPC requests
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the log service with the gRPC server
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	// Log the start of the gRPC server
	log.Printf("gRPC Server started on port %s", gRpcPort)

	// Serve gRPC requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
