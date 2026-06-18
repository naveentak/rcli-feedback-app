package feedback

import (
	"context"
	"fmt"
)

type GitHubStore interface {
	CreateIssue(ctx context.Context, req SubmitRequest) (*Ticket, error)
	GetIssue(ctx context.Context, number int) (*Ticket, error)
	ListIssues(ctx context.Context, filter ListFilter) ([]Ticket, error)
	AddComment(ctx context.Context, number int, body string) error
	UpdateStatus(ctx context.Context, number int, status Status) (*Ticket, error)
}

type Service struct {
	store GitHubStore
}

func NewService(store GitHubStore) *Service {
	return &Service{store: store}
}

func (s *Service) Submit(ctx context.Context, req SubmitRequest) (*Ticket, error) {
	if err := validateSubmit(req); err != nil {
		return nil, err
	}
	return s.store.CreateIssue(ctx, req)
}

func (s *Service) Get(ctx context.Context, number int) (*Ticket, error) {
	return s.store.GetIssue(ctx, number)
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]Ticket, error) {
	if filter.App != "" && filter.App != "all" && !IsValidApp(filter.App) {
		return nil, fmt.Errorf("invalid app: %s", filter.App)
	}
	if filter.Type != "" && !IsValidType(filter.Type) {
		return nil, fmt.Errorf("invalid type: %s", filter.Type)
	}
	if filter.Status != "" && filter.Status != "open" && !IsValidStatus(filter.Status) {
		return nil, fmt.Errorf("invalid status: %s", filter.Status)
	}
	return s.store.ListIssues(ctx, filter)
}

func (s *Service) Comment(ctx context.Context, number int, body string) error {
	if body == "" {
		return fmt.Errorf("comment body is required")
	}
	return s.store.AddComment(ctx, number, body)
}

func (s *Service) UpdateStatus(ctx context.Context, number int, status Status) (*Ticket, error) {
	if !IsValidStatus(string(status)) {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	return s.store.UpdateStatus(ctx, number, status)
}

func validateSubmit(req SubmitRequest) error {
	if !IsValidApp(string(req.App)) {
		return fmt.Errorf("invalid app: %s", req.App)
	}
	if !IsValidType(string(req.Type)) {
		return fmt.Errorf("invalid type: %s", req.Type)
	}
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Description == "" {
		return fmt.Errorf("description is required")
	}
	return nil
}