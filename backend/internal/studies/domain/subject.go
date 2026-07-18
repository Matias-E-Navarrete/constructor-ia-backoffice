package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("studies: not found")
	ErrNotOwner = errors.New("studies: you do not own this subject")
)

// Subject is a materia (course/subject) the user is studying.
type Subject struct {
	ID        string
	UserID    string
	Name      string
	NoteSlug  *string
	CreatedAt time.Time
}

// StudySession is one logged study block against a subject.
type StudySession struct {
	ID              string
	UserID          string
	SubjectID       string
	SessionDate     time.Time
	DurationMinutes int
	Topic           string
	CreatedAt       time.Time
}

type Repository interface {
	CreateSubject(ctx context.Context, s *Subject) error
	ListSubjectsByUser(ctx context.Context, userID string) ([]Subject, error)
	FindSubjectByID(ctx context.Context, id string) (*Subject, error)
	DeleteSubject(ctx context.Context, userID, id string) error

	LogSession(ctx context.Context, s *StudySession) error
	ListSessionsBySubject(ctx context.Context, subjectID string) ([]StudySession, error)
	ListSessionsByUser(ctx context.Context, userID string) ([]StudySession, error)
}
