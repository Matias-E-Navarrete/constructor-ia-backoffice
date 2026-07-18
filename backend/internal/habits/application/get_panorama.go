package application

import (
	"context"
	"time"

	"rimu/backend/internal/habits/domain"
)

const (
	panoramaHeatmapDays = 90
	panoramaRadarWeeks  = 8
)

// GetPanorama is a [PRO] use-case: heatmap + streaks + weekly consistency
// series for every one of the caller's active habits, feeding the Panorama
// dashboard (heatmap calendar, streak rings, radar chart).
type GetPanorama struct {
	Repo domain.Repository
}

func (uc *GetPanorama) Execute(ctx context.Context, userID string) ([]domain.PanoramaHabit, error) {
	habits, err := uc.Repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	result := make([]domain.PanoramaHabit, 0, len(habits))
	for _, h := range habits {
		if h.ArchivedAt != nil {
			continue
		}
		logs, err := uc.Repo.ListLogs(ctx, h.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, domain.ComputePanorama(h, logs, now, panoramaHeatmapDays, panoramaRadarWeeks))
	}
	return result, nil
}
