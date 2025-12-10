package tickets

import (
	"context"
	"errors"
	"innotech/internal/storage/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, t *postgres.Ticket) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *mockRepository) GetByID(ctx context.Context, id int) (*postgres.Ticket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postgres.Ticket), args.Error(1)
}

func (m *mockRepository) GetAll(ctx context.Context) ([]postgres.Ticket, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgres.Ticket), args.Error(1)
}

func (m *mockRepository) Update(ctx context.Context, t *postgres.Ticket) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *mockRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_Create_SetsStatusAndCallsRepo(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepository)
	svc := NewService(repo)

	tIn := &postgres.Ticket{Title: "t1", Message: "m1"}

	repo.On("Create", mock.Anything, tIn).Run(func(_ mock.Arguments) {
	}).Return(nil).Once()

	err := svc.Create(ctx, tIn)
	assert.NoError(t, err)
	assert.Equal(t, "open", tIn.Status)
	repo.AssertExpectations(t)
}

func TestService_Create_RepoError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepository)
	svc := NewService(repo)

	tIn := &postgres.Ticket{Title: "t2"}
	repo.On("Create", mock.Anything, tIn).Return(errors.New("db error")).Once()

	err := svc.Create(ctx, tIn)
	assert.Error(t, err)
	assert.Equal(t, "open", tIn.Status)
	repo.AssertExpectations(t)
}

func TestService_GetByID_ReturnsTicket(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepository)
	svc := NewService(repo)

	exp := &postgres.Ticket{ID: 1, Title: "t"}
	repo.On("GetByID", mock.Anything, 1).Return(exp, nil).Once()

	got, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
	repo.AssertExpectations(t)
}

func TestService_GetAll_ReturnsList(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepository)
	svc := NewService(repo)

	list := []postgres.Ticket{{ID: 1}, {ID: 2}}
	repo.On("GetAll", mock.Anything).Return(list, nil).Once()

	got, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, list, got)
	repo.AssertExpectations(t)
}

func TestService_Update_PassesThrough(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepository)
	svc := NewService(repo)

	tIn := &postgres.Ticket{ID: 5, Title: "t"}
	repo.On("Update", mock.Anything, tIn).Return(nil).Once()

	err := svc.Update(ctx, tIn)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_Delete_PassesThrough(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepository)
	svc := NewService(repo)

	repo.On("Delete", mock.Anything, 7).Return(nil).Once()

	err := svc.Delete(ctx, 7)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
