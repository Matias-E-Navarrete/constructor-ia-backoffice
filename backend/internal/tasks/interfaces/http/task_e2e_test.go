package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestTasks_CreateClassifyKanbanReorderComplete(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	var t1 struct{ ID string }
	status := s.Do(http.MethodPost, "/api/tasks", token, map[string]string{
		"title": "Finish Product Design Course", "category": "Estudios", "priority": "high",
	}, &t1)
	if status != http.StatusCreated || t1.ID == "" {
		t.Fatalf("create task: status=%d body=%+v", status, t1)
	}

	var t2 struct{ ID string }
	s.Do(http.MethodPost, "/api/tasks", token, map[string]string{"title": "Call to john", "category": "Trabajo"}, &t2)

	// Eisenhower classification
	status = s.Do(http.MethodPatch, "/api/tasks/"+t1.ID+"/quadrant", token, map[string]string{"quadrant": "do_now"}, nil)
	if status != http.StatusOK {
		t.Fatalf("set quadrant: expected 200, got %d", status)
	}

	var list []struct {
		ID       string
		Quadrant string
	}
	s.Do(http.MethodGet, "/api/tasks?quadrant=do_now", token, nil, &list)
	if len(list) != 1 || list[0].ID != t1.ID {
		t.Fatalf("list by quadrant: expected 1 task %s, got %+v", t1.ID, list)
	}

	// Kanban
	status = s.Do(http.MethodPatch, "/api/tasks/"+t2.ID+"/status", token, map[string]string{"status": "in_progress"}, nil)
	if status != http.StatusOK {
		t.Fatalf("set status: expected 200, got %d", status)
	}

	var kanban []struct{ ID string }
	s.Do(http.MethodGet, "/api/tasks?status=in_progress", token, nil, &kanban)
	if len(kanban) != 1 || kanban[0].ID != t2.ID {
		t.Fatalf("list by status: expected 1 task %s, got %+v", t2.ID, kanban)
	}

	// Manual reorder
	status = s.Do(http.MethodPost, "/api/tasks/reorder", token, map[string][]string{"ordered_ids": {t2.ID, t1.ID}}, nil)
	if status != http.StatusNoContent {
		t.Fatalf("reorder: expected 204, got %d", status)
	}

	var ordered []struct{ ID string }
	s.Do(http.MethodGet, "/api/tasks", token, nil, &ordered)
	if len(ordered) != 2 || ordered[0].ID != t2.ID || ordered[1].ID != t1.ID {
		t.Fatalf("expected reordered [%s, %s], got %+v", t2.ID, t1.ID, ordered)
	}

	// Complete
	var completed struct{ Status string }
	status = s.Do(http.MethodPost, "/api/tasks/"+t1.ID+"/complete", token, nil, &completed)
	if status != http.StatusOK || completed.Status != "done" {
		t.Fatalf("complete: status=%d body=%+v", status, completed)
	}
}
