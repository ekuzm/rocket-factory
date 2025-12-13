package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	serverAddress = "127.0.0.1:50051"
)

type InventoryStorage struct {
	parts map[string]*inventoryV1.Part
	mtx   sync.RWMutex
}

func (i *InventoryStorage) GetPart(uuid string) (*inventoryV1.Part, error) {
	i.mtx.RLock()
	defer i.mtx.RUnlock()

	part, ok := i.parts[uuid]
	if !ok {
		return nil, errors.New("part is not found")
	}

	return part, nil
}

func initTestStorage() (*InventoryStorage, error) {
	storage := InventoryStorage{
		parts: make(map[string]*inventoryV1.Part),
		mtx:   sync.RWMutex{},
	}

	file, err := os.OpenFile("parts.json", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			log.Printf("failed to close file: %v", cerr)
			return
		}
	}()
	json.NewDecoder(file).Decode(&storage.parts)

	return &storage, nil
}

type InventoryService struct {
	storage *InventoryStorage
	inventoryV1.UnimplementedInventoryV1ServiceServer
}

func NewInventaryService(storage *InventoryStorage) *InventoryService {
	return &InventoryService{
		storage: storage,
	}
}

func (i *InventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := i.storage.GetPart(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &inventoryV1.GetPartResponse{Part: part}, nil
}

// func (i *InventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
// 	return nil, nil
// }

func main() {
	storage, err := initTestStorage()
	if err != nil {
		log.Printf("failed to init test storage: %v", err)
		return
	}

	service := NewInventaryService(storage)

	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Printf("failed to listening gRPC server: %v", err)
		return
	}
	defer func() {
		cerr := lis.Close()
		if cerr != nil {
			log.Printf("failed to close the listener: %v", cerr)
			return
		}
	}()

	server := grpc.NewServer()
	if err != nil {
		log.Printf("failed to create gRPC server: %v")
		return
	}

	inventoryV1.RegisterInventoryV1ServiceServer(server, service)
	reflection.Register(server)

	go func() {
		log.Printf("gRPC server start to listening at %s", serverAddress)

		if err = server.Serve(lis); err != nil {
			log.Printf("failed to serve gRPC server: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Printf("Shutting down gRPC server...")
	server.GracefulStop()
	log.Println("Server stopped")
}
