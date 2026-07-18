package domain

import "time"

// HeatmapDay is one cell of the GitHub-style completion calendar.
type HeatmapDay struct {
	Date      string
	Completed bool
}

// PanoramaHabit bundles everything the [PRO] Panorama dashboard needs for a
// single habit: a heatmap calendar, streaks, and a per-week consistency
// series (used to draw the radar chart across habits).
type PanoramaHabit struct {
	HabitID       string
	Name          string
	CurrentStreak int
	BestStreak    int
	Heatmap       []HeatmapDay
	WeeklyRates   []float64 // oldest to newest
}

// ComputePanorama derives the heatmap/streak/radar data from raw log rows.
// Kept as a pure function, like ComputeStats, so it's unit-testable without a
// database. days controls the heatmap window; weeks controls how many
// 7-day buckets feed the radar's weekly consistency series.
func ComputePanorama(h Habit, logs []HabitLog, asOf time.Time, days, weeks int) PanoramaHabit {
	byDate := make(map[string]bool, len(logs))
	for _, l := range logs {
		byDate[l.LogDate.Format("2006-01-02")] = l.Completed
	}

	heatmap := make([]HeatmapDay, 0, days)
	start := asOf.AddDate(0, 0, -(days - 1))
	for d := start; !d.After(asOf); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		heatmap = append(heatmap, HeatmapDay{Date: key, Completed: byDate[key]})
	}

	currentStreak := 0
	for d := asOf; ; d = d.AddDate(0, 0, -1) {
		if byDate[d.Format("2006-01-02")] {
			currentStreak++
			continue
		}
		break
	}

	bestStreak, run := 0, 0
	for _, day := range heatmap {
		if day.Completed {
			run++
			if run > bestStreak {
				bestStreak = run
			}
		} else {
			run = 0
		}
	}

	weeklyRates := make([]float64, weeks)
	for w := 0; w < weeks; w++ {
		weekEnd := asOf.AddDate(0, 0, -7*w)
		completed := 0
		for i := 0; i < 7; i++ {
			d := weekEnd.AddDate(0, 0, -i)
			if byDate[d.Format("2006-01-02")] {
				completed++
			}
		}
		weeklyRates[weeks-1-w] = float64(completed) / 7.0
	}

	return PanoramaHabit{
		HabitID:       h.ID,
		Name:          h.Name,
		CurrentStreak: currentStreak,
		BestStreak:    bestStreak,
		Heatmap:       heatmap,
		WeeklyRates:   weeklyRates,
	}
}
