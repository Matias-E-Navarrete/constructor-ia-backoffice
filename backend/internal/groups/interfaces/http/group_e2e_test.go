package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestGroups_FamilySharingAndCoachingAccess(t *testing.T) {
	s := testutil.NewServer(t)

	coachToken := s.RegisterAndLogin("coach@example.com", "password123")
	clientToken := s.RegisterAndLogin("client@example.com", "password123")
	strangerToken := s.RegisterAndLogin("stranger@example.com", "password123")

	var client struct{ ID string }
	s.Do(http.MethodGet, "/api/me", clientToken, nil, &client)

	// --- family sharing ------------------------------------------------------------
	var family struct{ ID string }
	status := s.Do(http.MethodPost, "/api/groups", coachToken, map[string]string{
		"name": "Familia", "kind": "family",
	}, &family)
	if status != http.StatusCreated || family.ID == "" {
		t.Fatalf("create family group: status=%d body=%+v", status, family)
	}

	status = s.Do(http.MethodPost, "/api/groups/"+family.ID+"/members", coachToken, map[string]string{
		"email": "client@example.com",
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("invite family member: expected 201, got %d", status)
	}

	status = s.Do(http.MethodPost, "/api/finance/transactions", coachToken, map[string]interface{}{
		"type": "expense", "amount": 20, "category": "casa", "tx_date": "2026-07-01", "group_id": family.ID,
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("record shared transaction: expected 201, got %d", status)
	}

	var sharedForMember []struct{ ID string }
	status = s.Do(http.MethodGet, "/api/finance/transactions?groupId="+family.ID, clientToken, nil, &sharedForMember)
	if status != http.StatusOK || len(sharedForMember) != 1 {
		t.Fatalf("member views shared transaction: status=%d len=%d", status, len(sharedForMember))
	}

	status = s.Do(http.MethodGet, "/api/finance/transactions?groupId="+family.ID, strangerToken, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("stranger views shared transaction: expected 403, got %d", status)
	}

	// --- coaching --------------------------------------------------------------
	status = s.Do(http.MethodPost, "/api/groups", coachToken, map[string]string{
		"name": "Coaching", "kind": "coaching",
	}, nil)
	if status != http.StatusForbidden {
		t.Fatalf("create coaching group on free plan: expected 403, got %d", status)
	}

	s.Upgrade(coachToken)

	var coaching struct{ ID string }
	status = s.Do(http.MethodPost, "/api/groups", coachToken, map[string]string{
		"name": "Coaching", "kind": "coaching",
	}, &coaching)
	if status != http.StatusCreated || coaching.ID == "" {
		t.Fatalf("create coaching group on pro plan: status=%d body=%+v", status, coaching)
	}

	status = s.Do(http.MethodPost, "/api/groups/"+coaching.ID+"/members", coachToken, map[string]string{
		"email": "client@example.com",
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("invite client to coaching group: expected 201, got %d", status)
	}

	status = s.Do(http.MethodPost, "/api/habits", clientToken, map[string]string{"name": "Correr"}, nil)
	if status != http.StatusCreated {
		t.Fatalf("client creates habit: expected 201, got %d", status)
	}

	var clientHabits []struct{ Name string }
	status = s.Do(http.MethodGet, "/api/habits?userId="+client.ID, coachToken, nil, &clientHabits)
	if status != http.StatusOK || len(clientHabits) != 1 {
		t.Fatalf("coach views client habits: status=%d len=%d", status, len(clientHabits))
	}

	status = s.Do(http.MethodGet, "/api/habits?userId="+client.ID, strangerToken, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("non-coach views client habits: expected 403, got %d", status)
	}
}
