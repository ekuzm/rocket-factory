package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

const (
	serverAddress = "127.0.0.1:50051"
	dataFileName  = "parts.json"
)

type InventoryStorage struct {
	parts map[string]*inventoryV1.Part
	mtx   sync.RWMutex
}

func initTestStorage() (*InventoryStorage, error) {
	storage := InventoryStorage{
		parts: make(map[string]*inventoryV1.Part),
		mtx:   sync.RWMutex{},
	}

	file, err := os.OpenFile(dataFileName, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			log.Printf("Failed to close file: %v", cerr)
		}
	}()
	if err := json.NewDecoder(file).Decode(&storage.parts); err != nil {
		return nil, fmt.Errorf("failed to decode test data: %w", err)
	}

	return &storage, nil
}

func (i *InventoryStorage) GetPart(uuid string) *inventoryV1.Part {
	i.mtx.RLock()
	defer i.mtx.RUnlock()

	return i.parts[uuid]
}

func (i *InventoryStorage) ListParts() []*inventoryV1.Part {
	var parts []*inventoryV1.Part

	i.mtx.RLock()
	defer i.mtx.RUnlock()

	for _, part := range i.parts {
		parts = append(parts, part)
	}

	return parts
}

type PartPredicate func(*inventoryV1.Part) bool

func filterParts(inParts []*inventoryV1.Part, preds []PartPredicate) []*inventoryV1.Part {
	var outParts []*inventoryV1.Part

	for _, inPart := range inParts {
		if matchesAll(inPart, preds) {
			outParts = append(outParts, inPart)
		}
	}

	return outParts
}

func matchesAll(part *inventoryV1.Part, preds []PartPredicate) bool {
	for _, pred := range preds {
		if !pred(part) {
			return false
		}
	}

	return true
}

func byUUIDs(uuids []string) PartPredicate {
	return func(part *inventoryV1.Part) bool {
		for _, uuid := range uuids {
			if strings.EqualFold(strings.ToLower(part.Uuid), strings.ToLower(uuid)) {
				return true
			}
		}

		return false
	}
}

func byNames(names []string) PartPredicate {
	return func(part *inventoryV1.Part) bool {
		for _, name := range names {
			if strings.EqualFold(strings.ToLower(part.Name), strings.ToLower(name)) {
				return true
			}
		}

		return false
	}
}

func byCategories(categories []inventoryV1.Category) PartPredicate {
	return func(part *inventoryV1.Part) bool {
		for _, category := range categories {
			if part.Category == category {
				return true
			}
		}

		return false
	}
}

func byManufacturerCountries(countries []string) PartPredicate {
	return func(part *inventoryV1.Part) bool {
		for _, country := range countries {
			if strings.EqualFold(strings.ToLower(part.Manufacturer.Country), strings.ToLower(country)) {
				return true
			}
		}

		return false
	}
}

func byTags(tags []string) PartPredicate {
	return func(part *inventoryV1.Part) bool {
		var counter int

		for _, partTag := range part.Tags {
			for _, tag := range tags {
				if strings.EqualFold(strings.ToLower(partTag), strings.ToLower(tag)) {
					counter++
				}
			}
		}

		return counter == len(tags)
	}
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

func (i *InventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part := i.storage.GetPart(req.Uuid)
	if part == nil {
		return nil, status.Errorf(codes.NotFound, "the part with this uuid wasn't found")
	}

	return &inventoryV1.GetPartResponse{Part: part}, nil
}

func (i *InventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts := i.storage.ListParts()
	if parts == nil {
		return nil, status.Errorf(codes.NotFound, "parts weren't found")
	}

	var partPredicates []PartPredicate

	if len(req.Filter.Uuids) > 0 {
		partPredicates = append(partPredicates, byUUIDs(req.Filter.Uuids))
	}

	if len(req.Filter.Names) > 0 {
		partPredicates = append(partPredicates, byNames(req.Filter.Names))
	}

	if len(req.Filter.Categories) > 0 {
		partPredicates = append(partPredicates, byCategories(req.Filter.Categories))
	}

	if len(req.Filter.ManufacturerCountries) > 0 {
		partPredicates = append(partPredicates, byManufacturerCountries(req.Filter.ManufacturerCountries))
	}

	if len(req.Filter.Tags) > 0 {
		partPredicates = append(partPredicates, byTags(req.Filter.Tags))
	}

	parts = filterParts(parts, partPredicates)

	if len(parts) != len(req.Filter.Uuids) {
		return nil, status.Errorf(codes.NotFound, "some parts not found")
	}

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}

func main() {
	storage, err := initTestStorage()
	if err != nil {
		log.Fatalf("Failed to init test storage: %v", err)
	}

	service := NewInventaryService(storage)

	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Printf("failed to listening gRPC server: %v\n", err)
		return
	}
	defer func() {
		cerr := lis.Close()
		if cerr != nil {
			log.Printf("failed to close the listener: %v\n", cerr)
		}
	}()

	server := grpc.NewServer()

	inventoryV1.RegisterInventoryV1ServiceServer(server, service)
	reflection.Register(server)

	go func() {
		log.Printf("gRPC server start to listening at %s\n", serverAddress)

		if err = server.Serve(lis); err != nil {
			log.Printf("failed to serve gRPC server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Printf("Shutting down gRPC server...\n")
	server.GracefulStop()
	log.Printf("Stopped to server\n")
}
