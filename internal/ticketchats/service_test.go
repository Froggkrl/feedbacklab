package ticketchats

import (
	"context"
	"errors"
	"innotech/internal/storage/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Create(ctx context.Context, chat *postgres.TicketChat) error {
	args := m.Called(ctx, chat)
	return args.Error(0)
}

func (m *mockRepo) GetByID(ctx context.Context, id int) (*postgres.TicketChat, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postgres.TicketChat), args.Error(1)
}

func (m *mockRepo) GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketChat, error) {
	args := m.Called(ctx, ticketID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgres.TicketChat), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, chat *postgres.TicketChat) error {
	args := m.Called(ctx, chat)
	return args.Error(0)
}

func (m *mockRepo) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreate_ValidationFails_WhenEmptyMessage(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	ctx := context.Background()
	err := svc.Create(ctx, &postgres.TicketChat{Message: ""})
	assert.Error(t, err)
}

func TestCreate_HappyPath_CallsRepo(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	ctx := context.Background()
	chat := &postgres.TicketChat{TicketID: 1, SenderID: "u", Message: "hi"}
	repo.On("Create", ctx, chat).Return(nil).Once()

	err := svc.Create(ctx, chat)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError_ReturnsError(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	ctx := context.Background()
	chat := &postgres.TicketChat{TicketID: 1, SenderID: "u", Message: "hi"}
	repo.On("Create", ctx, chat).Return(errors.New("db")).Once()

	err := svc.Create(ctx, chat)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetByID_PassesThrough(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	ctx := context.Background()
	exp := &postgres.TicketChat{ID: 2}
	repo.On("GetByID", ctx, 2).Return(exp, nil).Once()

	got, err := svc.GetByID(ctx, 2)
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
	repo.AssertExpectations(t)
}

func TestGetByTicketID_PassesThrough(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	ctx := context.Background()
	list := []postgres.TicketChat{{ID: 1}, {ID: 2}}
	repo.On("GetByTicketID", ctx, 5).Return(list, nil).Once()

	got, err := svc.GetByTicketID(ctx, 5)
	assert.NoError(t, err)
	assert.Equal(t, list, got)
	repo.AssertExpectations(t)
}

func TestUpdate_Delete_PassThrough(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	ctx := context.Background()
	c := &postgres.TicketChat{ID: 3, Message: "ok"}
	repo.On("Update", ctx, c).Return(nil).Once()
	repo.On("Delete", ctx, 3).Return(nil).Once()

	assert.NoError(t, svc.Update(ctx, c))
	assert.NoError(t, svc.Delete(ctx, 3))
	repo.AssertExpectations(t)
}

func TestService_Methods_UseContext(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	type contextKey string
	const testKey contextKey = "test"

	ctx := context.WithValue(context.Background(), testKey, "value")
	chat := &postgres.TicketChat{ID: 1, Message: "test"}

	repo.On("Create", ctx, chat).Return(nil).Once()
	repo.On("GetByID", ctx, 1).Return(chat, nil).Once()
	repo.On("GetByTicketID", ctx, 1).Return([]postgres.TicketChat{*chat}, nil).Once()
	repo.On("Update", ctx, chat).Return(nil).Once()
	repo.On("Delete", ctx, 1).Return(nil).Once()

	assert.NoError(t, svc.Create(ctx, chat))

	result, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, chat, result)

	list, err := svc.GetByTicketID(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	assert.NoError(t, svc.Update(ctx, chat))
	assert.NoError(t, svc.Delete(ctx, 1))

	repo.AssertExpectations(t)
}
