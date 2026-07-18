package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestNotes_CreateLinkGraphExportGating(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	var note struct {
		Slug string `json:"slug"`
	}
	status := s.Do(http.MethodPost, "/api/notes", token, map[string]string{
		"title": "Meta de habito", "body": "Notas sobre mi meta",
	}, &note)
	if status != http.StatusCreated || note.Slug == "" {
		t.Fatalf("create note: status=%d body=%+v", status, note)
	}

	status = s.Do(http.MethodPost, "/api/habits", token, map[string]interface{}{
		"name": "Meditar", "note_slug": note.Slug,
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("create linked habit: expected 201, got %d", status)
	}

	status = s.Do(http.MethodGet, "/api/notes/graph", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("graph on free plan: expected 403, got %d", status)
	}
	status = s.Do(http.MethodGet, "/api/notes/export", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("export on free plan: expected 403, got %d", status)
	}

	s.Upgrade(token)

	var graph struct {
		Nodes []struct {
			ID   string `json:"id"`
			Kind string `json:"kind"`
		} `json:"nodes"`
		Edges []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"edges"`
	}
	status = s.Do(http.MethodGet, "/api/notes/graph", token, nil, &graph)
	if status != http.StatusOK {
		t.Fatalf("graph on pro plan: expected 200, got %d", status)
	}

	foundEdge := false
	for _, e := range graph.Edges {
		if e.To == note.Slug {
			foundEdge = true
		}
	}
	if !foundEdge {
		t.Fatalf("expected an edge pointing to note %q in graph %+v", note.Slug, graph)
	}

	status = s.Do(http.MethodGet, "/api/notes/export", token, nil, nil)
	if status != http.StatusOK {
		t.Fatalf("export on pro plan: expected 200, got %d", status)
	}
}
