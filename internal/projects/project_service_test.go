package projects

import (
	"context"
	"innotech/internal/storage/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProjectRepository struct{ mock.Mock }

func (m *MockProjectRepository) Create(ctx context.Context, p *postgres.Project) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockProjectRepository) GetAll(ctx context.Context) ([]postgres.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]postgres.Project), args.Error(1)
}
func (m *MockProjectRepository) GetByID(ctx context.Context, id int) (*postgres.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*postgres.Project), args.Error(1)
}
func (m *MockProjectRepository) Update(ctx context.Context, p *postgres.Project) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockProjectRepository) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func TestProjectService_CRUD(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	svc := NewService(mockRepo)

	ctx := context.Background()
	now := time.Now()
	project := &postgres.Project{ID: 1, Name: "Demo", DateCreated: now, DateUpdated: now}
	projects := []postgres.Project{*project}

	mockRepo.On("Create", ctx, project).Return(nil)
	mockRepo.On("GetAll", ctx).Return(projects, nil)
	mockRepo.On("GetByID", ctx, 1).Return(project, nil)
	mockRepo.On("Update", ctx, project).Return(nil)
	mockRepo.On("Delete", ctx, 1).Return(nil)

	assert.NoError(t, svc.Create(ctx, project))
	items, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, items, 1)

	got, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Demo", got.Name)

	assert.NoError(t, svc.Update(ctx, project))
	assert.NoError(t, svc.Delete(ctx, 1))
}
