package domain

import (
	"context"
	"time"
)

type Kind string

const (
	KindFeature Kind = "feature"
	KindBug     Kind = "bug"
	KindTask    Kind = "task"
)

type Status string

const (
	StatusPlanned    Status = "planned"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RoadmapItem struct {
	ID          string
	Title       string
	Description string
	Kind        Kind
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoadmapRepository interface {
	Create(ctx context.Context, item *RoadmapItem) error
	List(ctx context.Context) ([]RoadmapItem, error)
	UpdateStatus(ctx context.Context, id string, status Status) (*RoadmapItem, error)
}
