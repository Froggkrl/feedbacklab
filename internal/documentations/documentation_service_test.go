package documentations

import (
	"context"
	"innotech/internal/storage/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDocumentationRepository struct{ mock.Mock }

func (m *MockDocumentationRepository) Create(ctx context.Context, p *postgres.Documentation) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockDocumentationRepository) GetAll(ctx context.Context) ([]postgres.Documentation, error) {
	args := m.Called(ctx)
	return args.Get(0).([]postgres.Documentation), args.Error(1)
}
func (m *MockDocumentationRepository) GetByID(ctx context.Context, id int) (*postgres.Documentation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*postgres.Documentation), args.Error(1)
}
func (m *MockDocumentationRepository) Update(ctx context.Context, p *postgres.Documentation) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockDocumentationRepository) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func TestDocumentationService_CRUD(t *testing.T) {
	mockRepo := new(MockDocumentationRepository)
	svc := NewService(mockRepo)

	ctx := context.Background()
	now := time.Now()
	version := "1.0"
	uploadedBy := "user1"
	doc := &postgres.Documentation{ID: 1, ProjectID: 1, FilePath: "/api/spec.pdf", Version: &version, UploadedBy: &uploadedBy, DateCreated: now, DateUpdated: now}
	docs := []postgres.Documentation{*doc}

	mockRepo.On("Create", ctx, doc).Return(nil)
	mockRepo.On("GetAll", ctx).Return(docs, nil)
	mockRepo.On("GetByID", ctx, 1).Return(doc, nil)
	mockRepo.On("Update", ctx, doc).Return(nil)
	mockRepo.On("Delete", ctx, 1).Return(nil)

	assert.NoError(t, svc.Create(ctx, doc))
	list, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	got, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "/api/spec.pdf", got.FilePath)

	assert.NoError(t, svc.Update(ctx, doc))
	assert.NoError(t, svc.Delete(ctx, 1))
}
