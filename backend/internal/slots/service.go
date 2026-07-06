package slots

import "context"

type Repository interface {
	List(ctx context.Context, filters ListFilters) ([]Slot, error)
	GetByID(ctx context.Context, id string) (*Slot, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, filters ListFilters) ([]Slot, error) {
	return s.repo.List(ctx, filters)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Slot, error) {
	return s.repo.GetByID(ctx, id)
}
