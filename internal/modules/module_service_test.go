package modules

import (
	"context"
	"innotech/internal/storage/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockModuleRepository struct{ mock.Mock }

func (m *MockModuleRepository) Create(ctx context.Context, p *postgres.Module) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockModuleRepository) GetAll(ctx context.Context) ([]postgres.Module, error) {
	args := m.Called(ctx)
	return args.Get(0).([]postgres.Module), args.Error(1)
}
func (m *MockModuleRepository) GetByID(ctx context.Context, id int) (*postgres.Module, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*postgres.Module), args.Error(1)
}
func (m *MockModuleRepository) Update(ctx context.Context, p *postgres.Module) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockModuleRepository) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func TestModuleService_CRUD(t *testing.T) {
	mockRepo := new(MockModuleRepository)
	svc := NewService(mockRepo)

	ctx := context.Background()
	now := time.Now()
	module := &postgres.Module{ID: 1, ProjectID: 1, Name: "Core", DateCreated: now, DateUpdated: now}
	modules := []postgres.Module{*module}

	mockRepo.On("Create", ctx, module).Return(nil)
	mockRepo.On("GetAll", ctx).Return(modules, nil)
	mockRepo.On("GetByID", ctx, 1).Return(module, nil)
	mockRepo.On("Update", ctx, module).Return(nil)
	mockRepo.On("Delete", ctx, 1).Return(nil)

	assert.NoError(t, svc.Create(ctx, module))
	list, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	got, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Core", got.Name)

	assert.NoError(t, svc.Update(ctx, module))
	assert.NoError(t, svc.Delete(ctx, 1))
}
