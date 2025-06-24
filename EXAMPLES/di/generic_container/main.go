// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using the generic container from the dependency injection package
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abitofhelp/servicelib/di"
	"go.uber.org/zap"
)

// Implement the Repository interface
type ProductRepository struct {
	id string
}

func (r *ProductRepository) GetID() string {
	return r.id
}

func (r *ProductRepository) GetProducts() []string {
	return []string{"Product 1", "Product 2", "Product 3"}
}

// Initialize repository function
func InitProductRepository(ctx context.Context, connectionString string, logger *zap.Logger) (*ProductRepository, error) {
	logger.Info("Initializing product repository", zap.String("connection", connectionString))
	return &ProductRepository{id: "product-repo-1"}, nil
}

// Implement the DomainService interface
type ProductDomainService struct {
	id   string
	repo *ProductRepository
}

func (s *ProductDomainService) GetID() string {
	return s.id
}

func (s *ProductDomainService) GetFeaturedProducts() []string {
	allProducts := s.repo.GetProducts()
	// In a real application, this would apply domain logic to filter featured products
	return allProducts[:2] // Return first two products as featured
}

// Initialize domain service function
func InitProductDomainService(repo *ProductRepository) (*ProductDomainService, error) {
	return &ProductDomainService{
		id:   "product-domain-service-1",
		repo: repo,
	}, nil
}

// Implement the ApplicationService interface
type ProductApplicationService struct {
	id            string
	domainService *ProductDomainService
	repo          *ProductRepository
}

func (s *ProductApplicationService) GetID() string {
	return s.id
}

func (s *ProductApplicationService) GetProductRecommendations() []string {
	// In a real application, this would combine domain services and repositories
	// to provide application-specific functionality
	return s.domainService.GetFeaturedProducts()
}

// Initialize application service function
func InitProductApplicationService(domainService *ProductDomainService, repo *ProductRepository) (*ProductApplicationService, error) {
	return &ProductApplicationService{
		id:            "product-app-service-1",
		domainService: domainService,
		repo:          repo,
	}, nil
}

// Simple configuration type
type SimpleConfig struct{}

func main() {
	// Create a context
	ctx := context.Background()

	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create a simple configuration
	cfg := SimpleConfig{}

	// Connection string for the repository
	connectionString := "postgres://user:password@localhost:5432/products"

	// Create a generic container
	container, err := di.NewGenericAppContainer[*ProductRepository, *ProductDomainService, *ProductApplicationService, SimpleConfig](
		ctx,
		logger,
		cfg,
		connectionString,
		InitProductRepository,
		InitProductDomainService,
		InitProductApplicationService,
	)
	if err != nil {
		log.Fatalf("Failed to create generic container: %v", err)
	}

	// Use the container
	fmt.Println("Generic container created successfully")

	// Get the repository from the container
	repository := container.GetRepository()
	fmt.Printf("Repository ID: %s\n", repository.GetID())
	fmt.Printf("Products: %v\n", repository.GetProducts())

	// Get the domain service from the container
	domainService := container.GetDomainService()
	fmt.Printf("Domain service ID: %s\n", domainService.GetID())
	fmt.Printf("Featured products: %v\n", domainService.GetFeaturedProducts())

	// Get the application service from the container
	appService := container.GetApplicationService()
	fmt.Printf("Application service ID: %s\n", appService.GetID())
	fmt.Printf("Product recommendations: %v\n", appService.GetProductRecommendations())
}
