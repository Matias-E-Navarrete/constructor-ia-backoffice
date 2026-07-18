package application

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"

	financedomain "rimu/backend/internal/finance/domain"
	habitsdomain "rimu/backend/internal/habits/domain"
	"rimu/backend/internal/notes/domain"
	workoutsdomain "rimu/backend/internal/workouts/domain"
)

// ExportVault is a [PRO] use-case: renders the user's authored notes plus
// every entity manually linked to a note into real .md files, packaged as an
// actual Obsidian vault .zip (openable in the real Obsidian app).
type ExportVault struct {
	Notes    domain.Repository
	Habits   habitsdomain.Repository
	Workouts workoutsdomain.Repository
	Finance  financedomain.Repository
}

func (uc *ExportVault) Execute(ctx context.Context, userID string) ([]byte, error) {
	notes, err := uc.Notes.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	habits, err := uc.Habits.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	sessions, err := uc.Workouts.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	txs, err := uc.Finance.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for _, n := range notes {
		front := fmt.Sprintf("---\nkind: %s\ncreated: %s\nupdated: %s\n---\n\n", n.Kind, n.CreatedAt.Format("2006-01-02"), n.UpdatedAt.Format("2006-01-02"))
		writeFile(zw, n.Slug+".md", front+n.BodyMarkdown)
	}

	for _, h := range habits {
		if h.NoteSlug == nil {
			continue
		}
		body := fmt.Sprintf("---\nkind: habit\n---\n\n# %s\n\nRelacionado: [[%s]]\n", h.Name, *h.NoteSlug)
		writeFile(zw, "habito-"+domain.Slugify(h.Name)+".md", body)
	}

	for _, s := range sessions {
		if s.NoteSlug == nil {
			continue
		}
		body := fmt.Sprintf("---\nkind: workout\ndate: %s\n---\n\n# Entrenamiento %s\n\nRelacionado: [[%s]]\n", s.SessionDate.Format("2006-01-02"), s.SessionDate.Format("2006-01-02"), *s.NoteSlug)
		writeFile(zw, "entrenamiento-"+s.SessionDate.Format("2006-01-02")+".md", body)
	}

	for _, t := range txs {
		if t.NoteSlug == nil {
			continue
		}
		title := t.Description
		if title == "" {
			title = t.Category
		}
		body := fmt.Sprintf("---\nkind: finance\namount: %.2f\ntype: %s\n---\n\n# %s\n\nRelacionado: [[%s]]\n", t.Amount.Value(), t.Type, title, *t.NoteSlug)
		writeFile(zw, "finanzas-"+domain.Slugify(title)+".md", body)
	}

	writeFile(zw, ".obsidian/app.json", "{}\n")

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeFile(zw *zip.Writer, name, content string) {
	f, err := zw.Create(name)
	if err != nil {
		return
	}
	f.Write([]byte(content))
}
