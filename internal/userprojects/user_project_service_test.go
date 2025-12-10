package userprojects

import (
	"context"
	"innotech/internal/storage/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserProjectRepository struct{ mock.Mock }

func (m *MockUserProjectRepository) Create(ctx context.Context, p *postgres.UserProject) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockUserProjectRepository) GetAll(ctx context.Context) ([]postgres.UserProject, error) {
	args := m.Called(ctx)
	return args.Get(0).([]postgres.UserProject), args.Error(1)
}
func (m *MockUserProjectRepository) Get(ctx context.Context, userID string, projectID int) (*postgres.UserProject, error) {
	args := m.Called(ctx, userID, projectID)
	return args.Get(0).(*postgres.UserProject), args.Error(1)
}
func (m *MockUserProjectRepository) Update(ctx context.Context, p *postgres.UserProject) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockUserProjectRepository) Delete(ctx context.Context, userID string, projectID int) error {
	return m.Called(ctx, userID, projectID).Error(0)
}

func TestUserProjectService_CRUD(t *testing.T) {
	mockRepo := new(MockUserProjectRepository)
	svc := NewService(mockRepo)

	ctx := context.Background()
	now := time.Now()
	relation := &postgres.UserProject{UserID: "uuid-1", ProjectID: 2, Role: "developer", Permissions: []string{"read"}, DateCreated: now, DateUpdated: now}
	relations := []postgres.UserProject{*relation}

	mockRepo.On("Create", ctx, relation).Return(nil)
	mockRepo.On("GetAll", ctx).Return(relations, nil)
	mockRepo.On("Get", ctx, "uuid-1", 2).Return(relation, nil)
	mockRepo.On("Update", ctx, relation).Return(nil)
	mockRepo.On("Delete", ctx, "uuid-1", 2).Return(nil)

	assert.NoError(t, svc.Create(ctx, relation))
	list, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	got, err := svc.Get(ctx, "uuid-1", 2)
	assert.NoError(t, err)
	assert.Equal(t, "uuid-1", got.UserID)

	assert.NoError(t, svc.Update(ctx, relation))
	assert.NoError(t, svc.Delete(ctx, "uuid-1", 2))
}
