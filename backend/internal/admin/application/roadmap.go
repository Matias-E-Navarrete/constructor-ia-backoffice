package application

import (
	"context"

	"rimu/backend/internal/admin/domain"
)

type ListRoadmap struct {
	Repo domain.RoadmapRepository
}

func (uc *ListRoadmap) Execute(ctx context.Context) ([]domain.RoadmapItem, error) {
	return uc.Repo.List(ctx)
}

type CreateRoadmapItemInput struct {
	Title       string
	Description string
	Kind        domain.Kind
}

type CreateRoadmapItem struct {
	Repo domain.RoadmapRepository
}

func (uc *CreateRoadmapItem) Execute(ctx context.Context, in CreateRoadmapItemInput) (*domain.RoadmapItem, error) {
	item := &domain.RoadmapItem{
		Title:       in.Title,
		Description: in.Description,
		Kind:        in.Kind,
		Status:      domain.StatusPlanned,
	}
	if err := uc.Repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

type UpdateRoadmapStatus struct {
	Repo domain.RoadmapRepository
}

func (uc *UpdateRoadmapStatus) Execute(ctx context.Context, id string, status domain.Status) (*domain.RoadmapItem, error) {
	return uc.Repo.UpdateStatus(ctx, id, status)
}
