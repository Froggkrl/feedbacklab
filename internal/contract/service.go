package contract

// Service handles business logic for contract operations.
type Service struct {
	repo *Repository
}

// NewService creates a new Service instance.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetAll retrieves all contracts.
func (s *Service) GetAll() ([]Contract, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a contract by its ID.
func (s *Service) GetByID(id int) (*Contract, error) {
	return s.repo.GetByID(id)
}

// Create creates a new contract.
func (s *Service) Create(c *Contract) error {
	return s.repo.Create(c)
}

// Update updates an existing contract.
func (s *Service) Update(c *Contract) error {
	return s.repo.Update(c)
}

// Delete deletes a contract by its ID.
func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
