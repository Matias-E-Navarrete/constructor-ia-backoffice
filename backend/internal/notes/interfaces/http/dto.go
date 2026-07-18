package http

import (
	"rimu/backend/internal/notes/application"
	"rimu/backend/internal/notes/domain"
)

type createNoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	Kind  string `json:"kind,omitempty"`
}

type updateNoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type noteResponse struct {
	Slug  string   `json:"slug"`
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Kind  string   `json:"kind"`
	Links []string `json:"links"`
}

func toNoteResponse(n domain.Note) noteResponse {
	return noteResponse{Slug: n.Slug, Title: n.Title, Body: n.BodyMarkdown, Kind: string(n.Kind), Links: n.Links()}
}

type graphNodeResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Kind  string `json:"kind"`
}

type graphEdgeResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type graphResponse struct {
	Nodes []graphNodeResponse `json:"nodes"`
	Edges []graphEdgeResponse `json:"edges"`
}

func toGraphResponse(g application.Graph) graphResponse {
	nodes := make([]graphNodeResponse, 0, len(g.Nodes))
	for _, n := range g.Nodes {
		nodes = append(nodes, graphNodeResponse{ID: n.ID, Title: n.Title, Kind: n.Kind})
	}
	edges := make([]graphEdgeResponse, 0, len(g.Edges))
	for _, e := range g.Edges {
		edges = append(edges, graphEdgeResponse{From: e.From, To: e.To})
	}
	return graphResponse{Nodes: nodes, Edges: edges}
}
