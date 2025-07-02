package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	
	// Import your real generated proto code
	pb "github.com/lucasnoah/rescue-titan-proto/gen/go/services"
)

// RealDatabaseClient creates a real connection to your database service
func createRealDatabaseClient() (pb.DatabaseServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database service: %w", err)
	}
	
	client := pb.NewDatabaseServiceClient(conn)
	return client, conn, nil
}

// handleDatabaseServiceListProjectsREAL - REAL implementation that calls your actual service
func (s *MCPServer) handleDatabaseServiceListProjectsREAL(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	log.Printf("🔥 CALLING REAL DATABASE SERVICE FOR PROJECTS!")
	
	// Create real client
	client, conn, err := createRealDatabaseClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}
	defer conn.Close()
	
	// Create the real request
	req := &pb.ListProjectsRequest{
		PageSize: 10, // Get first 10 projects
	}
	
	// Add any filters from params
	if filter, ok := params["filter"].(string); ok && filter != "" {
		req.Filter = filter
	}
	
	log.Printf("📞 Making REAL gRPC call to DatabaseService.ListProjects...")
	
	// Make the REAL call to your service
	resp, err := client.ListProjects(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("REAL database service call failed: %w", err)
	}
	
	log.Printf("✅ GOT REAL DATA! Found %d projects, total count: %d", len(resp.Projects), resp.TotalCount)
	
	// Convert to JSON-friendly format
	result := map[string]interface{}{
		"projects":     resp.Projects,
		"total_count":  resp.TotalCount,
		"next_token":   resp.NextPageToken,
		"source":       "REAL_DATABASE_SERVICE",
		"proof":        "This data came from your actual running database service!",
	}
	
	return result, nil
}

// handleDatabaseServiceListUsersREAL - REAL implementation that calls your actual service
func (s *MCPServer) handleDatabaseServiceListUsersREAL(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	log.Printf("🔥 CALLING REAL DATABASE SERVICE FOR USERS!")
	
	// Create real client
	client, conn, err := createRealDatabaseClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}
	defer conn.Close()
	
	// Create the real request
	req := &pb.ListUsersRequest{
		PageSize: 10, // Get first 10 users
	}
	
	log.Printf("📞 Making REAL gRPC call to DatabaseService.ListUsers...")
	
	// Make the REAL call to your service
	resp, err := client.ListUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("REAL database service call failed: %w", err)
	}
	
	log.Printf("✅ GOT REAL DATA! Found %d users", len(resp.Users))
	
	// Convert to JSON-friendly format
	result := map[string]interface{}{
		"users":        resp.Users,
		"next_token":   resp.NextPageToken,
		"source":       "REAL_DATABASE_SERVICE",
		"proof":        "This data came from your actual running database service!",
	}
	
	return result, nil
}