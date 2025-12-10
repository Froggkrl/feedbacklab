package ticketattachments

import (
	"context"
	"errors"
	"innotech/internal/storage/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepoTA struct{ mock.Mock }

func (m *mockRepoTA) Create(ctx context.Context, att *postgres.TicketAttachment) error {
	return m.Called(ctx, att).Error(0)
}
func (m *mockRepoTA) GetByID(ctx context.Context, id int) (*postgres.TicketAttachment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postgres.TicketAttachment), args.Error(1)
}
func (m *mockRepoTA) GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketAttachment, error) {
	args := m.Called(ctx, ticketID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgres.TicketAttachment), args.Error(1)
}
func (m *mockRepoTA) Update(ctx context.Context, att *postgres.TicketAttachment) error {
	return m.Called(ctx, att).Error(0)
}
func (m *mockRepoTA) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func TestCreate_ValidationFails_WhenEmptyFilePath(t *testing.T) {
	repo := new(mockRepoTA)
	svc := NewService(repo)

	err := svc.Create(context.Background(), &postgres.TicketAttachment{FilePath: ""})
	assert.Error(t, err)
}

func TestCreate_HappyPath_CallsRepo(t *testing.T) {
	repo := new(mockRepoTA)
	svc := NewService(repo)

	att := &postgres.TicketAttachment{TicketID: 1, FilePath: "p"}
	repo.On("Create", mock.Anything, att).Return(nil).Once()

	err := svc.Create(context.Background(), att)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError_ReturnsError(t *testing.T) {
	repo := new(mockRepoTA)
	svc := NewService(repo)

	att := &postgres.TicketAttachment{TicketID: 1, FilePath: "p"}
	repo.On("Create", mock.Anything, att).Return(errors.New("db")).Once()

	err := svc.Create(context.Background(), att)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetAndUpdateDelete_PassesThrough(t *testing.T) {
	repo := new(mockRepoTA)
	svc := NewService(repo)
	ctx := context.Background()

	exp := &postgres.TicketAttachment{ID: 10}

	repo.On("GetByID", mock.Anything, 10).Return(exp, nil).Once()
	repo.On("GetByTicketID", mock.Anything, 2).Return([]postgres.TicketAttachment{{ID: 1}}, nil).Once()
	repo.On("Update", mock.Anything, exp).Return(nil).Once()
	repo.On("Delete", mock.Anything, 10).Return(nil).Once()

	got, err := svc.GetByID(ctx, 10)
	assert.NoError(t, err)
	assert.Equal(t, exp, got)

	list, err := svc.GetByTicketID(ctx, 2)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	assert.NoError(t, svc.Update(ctx, exp))
	assert.NoError(t, svc.Delete(ctx, 10))
	repo.AssertExpectations(t)
}
