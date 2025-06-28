// Copyright (c) 2025 A Bit of Help, Inc.

package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestMockStringRepository tests that the MockStringRepository satisfies the Repository interface
func TestMockStringRepository(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mocks.NewMockStringRepository(ctrl)

	// Test GetByID
	t.Run("GetByID", func(t *testing.T) {
		// Setup expectations
		mockRepo.EXPECT().GetByID(gomock.Any(), "1").Return("test-entity", nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), "2").Return("", errors.New("entity not found"))

		// Test successful retrieval
		entity, err := mockRepo.GetByID(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, "test-entity", entity)

		// Test error case
		entity, err = mockRepo.GetByID(context.Background(), "2")
		assert.Error(t, err)
		assert.Equal(t, "", entity)
		assert.Equal(t, "entity not found", err.Error())
	})

	// Test GetAll
	t.Run("GetAll", func(t *testing.T) {
		// Setup expectations
		mockRepo.EXPECT().GetAll(gomock.Any()).Return([]string{"entity1", "entity2"}, nil)
		mockRepo.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("database error"))

		// Test successful retrieval
		entities, err := mockRepo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []string{"entity1", "entity2"}, entities)

		// Test error case
		entities, err = mockRepo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, entities)
		assert.Equal(t, "database error", err.Error())
	})

	// Test Save
	t.Run("Save", func(t *testing.T) {
		// Setup expectations
		mockRepo.EXPECT().Save(gomock.Any(), "entity1").Return(nil)
		mockRepo.EXPECT().Save(gomock.Any(), "entity2").Return(errors.New("save error"))

		// Test successful save
		err := mockRepo.Save(context.Background(), "entity1")
		assert.NoError(t, err)

		// Test error case
		err = mockRepo.Save(context.Background(), "entity2")
		assert.Error(t, err)
		assert.Equal(t, "save error", err.Error())
	})
}

// TestMockRepositoryFactory tests that the MockRepositoryFactory satisfies the RepositoryFactory interface
func TestMockRepositoryFactory(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository factory
	mockFactory := mocks.NewMockRepositoryFactory(ctrl)

	// Test GetRepository
	t.Run("GetRepository", func(t *testing.T) {
		// Create a mock repository to return
		mockRepo := mocks.NewMockStringRepository(ctrl)

		// Setup expectations
		mockFactory.EXPECT().GetRepository().Return(mockRepo)
		mockFactory.EXPECT().GetRepository().Return(nil)

		// Test successful retrieval
		repo := mockFactory.GetRepository()
		assert.NotNil(t, repo)
		assert.IsType(t, &mocks.MockStringRepository{}, repo)

		// Test nil case
		repo = mockFactory.GetRepository()
		assert.Nil(t, repo)
	})
}

// TestRepositoryInterface verifies that the Repository interface is properly defined
func TestRepositoryInterface(t *testing.T) {
	// This test doesn't actually execute any code, it just verifies that the types
	// satisfy the interfaces at compile time

	// Verify that MockStringRepository implements Repository[string]
	var _ Repository[string] = (*mocks.MockStringRepository)(nil)
}

// TestRepositoryFactoryInterface verifies that the RepositoryFactory interface is properly defined
func TestRepositoryFactoryInterface(t *testing.T) {
	// This test doesn't actually execute any code, it just verifies that the types
	// satisfy the interfaces at compile time

	// Verify that MockRepositoryFactory implements RepositoryFactory
	var _ RepositoryFactory = (*mocks.MockRepositoryFactory)(nil)
}
