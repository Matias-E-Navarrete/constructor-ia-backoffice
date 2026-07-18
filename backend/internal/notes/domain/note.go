package domain

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
)

var ErrNotFound = errors.New("notes: not found")
var ErrSlugTaken = errors.New("notes: slug already used")

type Kind string

const (
	KindFreeform Kind = "freeform"
	KindProfile  Kind = "profile"
)

type Note struct {
	ID           string
	UserID       string
	Slug         string
	Title        string
	BodyMarkdown string
	Kind         Kind
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var wikilinkPattern = regexp.MustCompile(`\[\[([^\]]+)\]\]`)

// Links extracts every [[wikilink]] target referenced in the note body, the
// core Obsidian mechanic this vault is built around.
func (n Note) Links() []string {
	matches := wikilinkPattern.FindAllStringSubmatch(n.BodyMarkdown, -1)
	links := make([]string, 0, len(matches))
	for _, m := range matches {
		links = append(links, strings.TrimSpace(m[1]))
	}
	return links
}

func Slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

type Repository interface {
	Create(ctx context.Context, n *Note) error
	Update(ctx context.Context, n *Note) error
	FindBySlug(ctx context.Context, userID, slug string) (*Note, error)
	ListByUser(ctx context.Context, userID string) ([]Note, error)
	Delete(ctx context.Context, userID, slug string) error
}
