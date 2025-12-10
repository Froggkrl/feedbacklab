package messageattachments

import (
	"context"
	"errors"
	"innotech/internal/storage/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock repository implementing new interface
type mockRepoMA struct {
	mock.Mock
}

func (m *mockRepoMA) Create(ctx context.Context, att *postgres.MessageAttachment) error {
	return m.Called(ctx, att).Error(0)
}

func (m *mockRepoMA) GetByID(ctx context.Context, id int) (*postgres.MessageAttachment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postgres.MessageAttachment), args.Error(1)
}

func (m *mockRepoMA) GetByChatID(ctx context.Context, chatID int) ([]postgres.MessageAttachment, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgres.MessageAttachment), args.Error(1)
}

func (m *mockRepoMA) Update(ctx context.Context, att *postgres.MessageAttachment) error {
	return m.Called(ctx, att).Error(0)
}

func (m *mockRepoMA) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func TestCreate_ValidationFails_WhenMissingFields(t *testing.T) {
	repo := new(mockRepoMA)
	svc := NewService(repo)
	ctx := context.Background()

	err := svc.Create(ctx, &postgres.MessageAttachment{ChatID: 1, FilePath: ""})
	assert.Error(t, err)

	err = svc.Create(ctx, &postgres.MessageAttachment{ChatID: 0, FilePath: "p"})
	assert.Error(t, err)
}

func TestCreate_HappyPath_CallsRepo(t *testing.T) {
	repo := new(mockRepoMA)
	svc := NewService(repo)
	ctx := context.Background()

	a := &postgres.MessageAttachment{ChatID: 2, FilePath: "p"}
	repo.On("Create", ctx, a).Return(nil).Once()

	err := svc.Create(ctx, a)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError_ReturnsError(t *testing.T) {
	repo := new(mockRepoMA)
	svc := NewService(repo)
	ctx := context.Background()

	a := &postgres.MessageAttachment{ChatID: 2, FilePath: "p"}
	repo.On("Create", ctx, a).Return(errors.New("db")).Once()

	err := svc.Create(ctx, a)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetUpdateDelete_PassesThrough(t *testing.T) {
	repo := new(mockRepoMA)
	svc := NewService(repo)
	ctx := context.Background()

	exp := &postgres.MessageAttachment{ID: 11}

	repo.On("GetByID", ctx, 11).Return(exp, nil).Once()
	repo.On("GetByChatID", ctx, 3).Return([]postgres.MessageAttachment{{ID: 1}}, nil).Once()
	repo.On("Update", ctx, exp).Return(nil).Once()
	repo.On("Delete", ctx, 11).Return(nil).Once()

	got, err := svc.GetByID(ctx, 11)
	assert.NoError(t, err)
	assert.Equal(t, exp, got)

	list, err := svc.GetByChatID(ctx, 3)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	assert.NoError(t, svc.Update(ctx, exp))
	assert.NoError(t, svc.Delete(ctx, 11))
	repo.AssertExpectations(t)
}
