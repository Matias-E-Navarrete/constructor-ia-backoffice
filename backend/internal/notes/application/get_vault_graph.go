package application

import (
	"context"
	"strings"

	financedomain "rimu/backend/internal/finance/domain"
	habitsdomain "rimu/backend/internal/habits/domain"
	"rimu/backend/internal/notes/domain"
	workoutsdomain "rimu/backend/internal/workouts/domain"
)

type GraphNode struct {
	ID    string
	Title string
	Kind  string // "note" | "habit" | "workout" | "finance"
}

type GraphEdge struct {
	From string
	To   string
}

type Graph struct {
	Nodes []GraphNode
	Edges []GraphEdge
}

// GetVaultGraph is a [PRO] use-case: it builds the backlink graph out of the
// user's authored notes (via [[wikilinks]]) plus every habit/workout/
// transaction that was manually linked to a note via note_slug.
type GetVaultGraph struct {
	Notes    domain.Repository
	Habits   habitsdomain.Repository
	Workouts workoutsdomain.Repository
	Finance  financedomain.Repository
}

func (uc *GetVaultGraph) Execute(ctx context.Context, userID string) (Graph, error) {
	notes, err := uc.Notes.ListByUser(ctx, userID)
	if err != nil {
		return Graph{}, err
	}

	graph := Graph{}
	bySlug := make(map[string]domain.Note, len(notes))
	byTitle := make(map[string]string, len(notes)) // lowercase title -> slug
	for _, n := range notes {
		graph.Nodes = append(graph.Nodes, GraphNode{ID: n.Slug, Title: n.Title, Kind: "note"})
		bySlug[n.Slug] = n
		byTitle[strings.ToLower(n.Title)] = n.Slug
	}

	for _, n := range notes {
		for _, link := range n.Links() {
			target, ok := bySlug[domain.Slugify(link)]
			if !ok {
				if slug, ok2 := byTitle[strings.ToLower(link)]; ok2 {
					target = bySlug[slug]
				} else {
					continue
				}
			}
			graph.Edges = append(graph.Edges, GraphEdge{From: n.Slug, To: target.Slug})
		}
	}

	habits, err := uc.Habits.ListByUser(ctx, userID)
	if err != nil {
		return Graph{}, err
	}
	for _, h := range habits {
		if h.NoteSlug == nil {
			continue
		}
		id := "habit:" + h.ID
		graph.Nodes = append(graph.Nodes, GraphNode{ID: id, Title: h.Name, Kind: "habit"})
		graph.Edges = append(graph.Edges, GraphEdge{From: id, To: *h.NoteSlug})
	}

	sessions, err := uc.Workouts.ListByUser(ctx, userID)
	if err != nil {
		return Graph{}, err
	}
	for _, s := range sessions {
		if s.NoteSlug == nil {
			continue
		}
		id := "workout:" + s.ID
		graph.Nodes = append(graph.Nodes, GraphNode{ID: id, Title: "Entrenamiento " + s.SessionDate.Format("2006-01-02"), Kind: "workout"})
		graph.Edges = append(graph.Edges, GraphEdge{From: id, To: *s.NoteSlug})
	}

	txs, err := uc.Finance.ListByUser(ctx, userID)
	if err != nil {
		return Graph{}, err
	}
	for _, t := range txs {
		if t.NoteSlug == nil {
			continue
		}
		id := "finance:" + t.ID
		title := t.Description
		if title == "" {
			title = t.Category
		}
		graph.Nodes = append(graph.Nodes, GraphNode{ID: id, Title: title, Kind: "finance"})
		graph.Edges = append(graph.Edges, GraphEdge{From: id, To: *t.NoteSlug})
	}

	return graph, nil
}
