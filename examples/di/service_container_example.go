// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using the service container from the dependency injection package
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abitofhelp/servicelib/config"
	"github.com/abitofhelp/servicelib/di"
	"go.uber.org/zap"
)

// Implement the Repository interface
type MyRepository struct {
	id string
}

func (r *MyRepository) GetID() string {
	return r.id
}

func (r *MyRepository) GetData() string {
	return "Sample data from repository"
}

// Implement the DomainService interface
type MyDomainService struct {
	id   string
	repo *MyRepository
}

func NewMyDomainService(repo *MyRepository) (*MyDomainService, error) {
	return &MyDomainService{
		id:   "domain-service-1",
		repo: repo,
	}, nil
}

func (s *MyDomainService) GetID() string {
	return s.id
}

func (s *MyDomainService) ProcessData() string {
	data := s.repo.GetData()
	return fmt.Sprintf("Processed: %s", data)
}

// Implement the ApplicationService interface
type MyApplicationService struct {
	id            string
	domainService *MyDomainService
	repo          *MyRepository
}

func NewMyApplicationService(domainService *MyDomainService, repo *MyRepository) (*MyApplicationService, error) {
	return &MyApplicationService{
		id:            "app-service-1",
		domainService: domainService,
		repo:          repo,
	}, nil
}

func (s *MyApplicationService) GetID() string {
	return s.id
}

func (s *MyApplicationService) Execute() string {
	return s.domainService.ProcessData()
}

// Implement the AppConfig interface
type MyAppConfig struct {
	AppName    string
	AppVersion string
	AppEnv     string
}

func (c *MyAppConfig) GetName() string {
	return c.AppName
}

func (c *MyAppConfig) GetVersion() string {
	return c.AppVersion
}

func (c *MyAppConfig) GetEnvironment() string {
	return c.AppEnv
}

// Implement the DatabaseConfig interface
type MyDatabaseConfig struct {
	DBType             string
	DBConnectionString string
	DBName             string
	Collections        map[string]string
}

func (c *MyDatabaseConfig) GetType() string {
	return c.DBType
}

func (c *MyDatabaseConfig) GetConnectionString() string {
	return c.DBConnectionString
}

func (c *MyDatabaseConfig) GetDatabaseName() string {
	return c.DBName
}

func (c *MyDatabaseConfig) GetCollectionName(entityType string) string {
	if name, ok := c.Collections[entityType]; ok {
		return name
	}
	return entityType + "s" // Default to pluralized entity type
}

// Implement the Config interface
type MyConfig struct {
	App      *MyAppConfig
	Database *MyDatabaseConfig
}

func (c *MyConfig) GetApp() config.AppConfig {
	return c.App
}

func (c *MyConfig) GetDatabase() config.DatabaseConfig {
	return c.Database
}

func main() {
	// Create a context
	ctx := context.Background()

	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create a configuration
	cfg := &MyConfig{
		App: &MyAppConfig{
			AppName:    "service-example",
			AppVersion: "1.0.0",
			AppEnv:     "development",
		},
		Database: &MyDatabaseConfig{
			DBType:             "postgres",
			DBConnectionString: "postgres://user:password@localhost:5432/mydb",
			DBName:             "mydb",
			Collections: map[string]string{
				"user": "users",
				"item": "items",
			},
		},
	}

	// Create a repository
	repo := &MyRepository{id: "repo-1"}

	// Create a service container
	container, err := di.NewServiceContainer[*MyRepository, *MyDomainService, *MyApplicationService, *MyConfig](
		ctx,
		logger,
		cfg,
		repo,
		NewMyDomainService,
		NewMyApplicationService,
	)
	if err != nil {
		log.Fatalf("Failed to create service container: %v", err)
	}

	// Use the container
	fmt.Println("Service container created successfully")

	// Get the repository from the container
	repository := container.GetRepository()
	fmt.Printf("Repository ID: %s\n", repository.GetID())
	fmt.Printf("Repository data: %s\n", repository.GetData())

	// Get the domain service from the container
	domainService := container.GetDomainService()
	fmt.Printf("Domain service ID: %s\n", domainService.GetID())
	fmt.Printf("Domain service result: %s\n", domainService.ProcessData())

	// Get the application service from the container
	appService := container.GetApplicationService()
	fmt.Printf("Application service ID: %s\n", appService.GetID())
	fmt.Printf("Application service result: %s\n", appService.Execute())
}
