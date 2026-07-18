package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestInbox_CaptureListPinDelete(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	var item struct {
		ID     string
		Pinned bool
	}
	status := s.Do(http.MethodPost, "/api/inbox", token, map[string]string{"content": "Read e-mail"}, &item)
	if status != http.StatusCreated || item.ID == "" {
		t.Fatalf("capture: status=%d body=%+v", status, item)
	}

	var list []struct{ ID string }
	status = s.Do(http.MethodGet, "/api/inbox", token, nil, &list)
	if status != http.StatusOK || len(list) != 1 {
		t.Fatalf("list: status=%d len=%d", status, len(list))
	}

	status = s.Do(http.MethodPatch, "/api/inbox/"+item.ID+"/pin", token, map[string]bool{"pinned": true}, nil)
	if status != http.StatusNoContent {
		t.Fatalf("pin: expected 204, got %d", status)
	}

	status = s.Do(http.MethodDelete, "/api/inbox/"+item.ID, token, nil, nil)
	if status != http.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d", status)
	}

	list = nil
	s.Do(http.MethodGet, "/api/inbox", token, nil, &list)
	if len(list) != 0 {
		t.Fatalf("expected empty inbox after delete, got %+v", list)
	}
}
