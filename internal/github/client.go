package github

import (
	"context"
	"fmt"
	"strings"

	gh "github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"

	"github.com/rcli/feedback/internal/feedback"
)

type Client struct {
	owner  string
	repo   string
	client *gh.Client
}

func NewClient(token, owner, repo string) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return &Client{
		owner:  owner,
		repo:   repo,
		client: gh.NewClient(tc),
	}
}

func (c *Client) CreateIssue(ctx context.Context, req feedback.SubmitRequest) (*feedback.Ticket, error) {
	labels := []string{
		req.App.Label(),
		req.Type.Label(),
	}

	body := req.Description
	if req.Reporter != "" {
		body = fmt.Sprintf("**Reporter:** %s\n\n%s", req.Reporter, req.Description)
	}

	issue, _, err := c.client.Issues.Create(ctx, c.owner, c.repo, &gh.IssueRequest{
		Title:  gh.String(req.Title),
		Body:   gh.String(body),
		Labels: &labels,
	})
	if err != nil {
		return nil, fmt.Errorf("create issue: %w", err)
	}

	return c.toTicket(issue), nil
}

func (c *Client) GetIssue(ctx context.Context, number int) (*feedback.Ticket, error) {
	issue, _, err := c.client.Issues.Get(ctx, c.owner, c.repo, number)
	if err != nil {
		return nil, fmt.Errorf("get issue #%d: %w", number, err)
	}
	return c.toTicket(issue), nil
}

func (c *Client) ListIssues(ctx context.Context, filter feedback.ListFilter) ([]feedback.Ticket, error) {
	opts := &gh.IssueListByRepoOptions{
		ListOptions: gh.ListOptions{PerPage: 100},
	}

	switch filter.State {
	case "closed":
		opts.State = "closed"
	case "all":
		opts.State = "all"
	default:
		opts.State = "open"
	}

	if filter.App != "" && filter.App != "all" {
		opts.Labels = []string{feedback.App(filter.App).Label()}
	}

	var labels []string
	if filter.Type != "" {
		labels = append(labels, feedback.TicketType(filter.Type).Label())
	}
	if filter.Status != "" && filter.Status != "open" {
		labels = append(labels, feedback.Status(filter.Status).Label())
	}
	if len(labels) > 0 {
		if len(opts.Labels) > 0 {
			opts.Labels = append(opts.Labels, labels...)
		} else {
			opts.Labels = labels
		}
	}

	var tickets []feedback.Ticket
	for {
		issues, resp, err := c.client.Issues.ListByRepo(ctx, c.owner, c.repo, opts)
		if err != nil {
			return nil, fmt.Errorf("list issues: %w", err)
		}
		for _, issue := range issues {
			if issue.PullRequestLinks != nil {
				continue
			}
			tickets = append(tickets, *c.toTicket(issue))
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return tickets, nil
}

func (c *Client) AddComment(ctx context.Context, number int, body string) error {
	_, _, err := c.client.Issues.CreateComment(ctx, c.owner, c.repo, number, &gh.IssueComment{
		Body: gh.String(body),
	})
	if err != nil {
		return fmt.Errorf("add comment to #%d: %w", number, err)
	}
	return nil
}

func (c *Client) UpdateStatus(ctx context.Context, number int, status feedback.Status) (*feedback.Ticket, error) {
	issue, _, err := c.client.Issues.Get(ctx, c.owner, c.repo, number)
	if err != nil {
		return nil, fmt.Errorf("get issue #%d: %w", number, err)
	}

	var keep []string
	for _, label := range issue.Labels {
		name := label.GetName()
		if !strings.HasPrefix(name, "status:") {
			keep = append(keep, name)
		}
	}
	keep = append(keep, status.Label())

	updated, _, err := c.client.Issues.Edit(ctx, c.owner, c.repo, number, &gh.IssueRequest{
		Labels: &keep,
	})
	if err != nil {
		return nil, fmt.Errorf("update issue #%d: %w", number, err)
	}

	if status == feedback.StatusDone {
		updated, _, err = c.client.Issues.Edit(ctx, c.owner, c.repo, number, &gh.IssueRequest{
			State: gh.String("closed"),
		})
		if err != nil {
			return nil, fmt.Errorf("close issue #%d: %w", number, err)
		}
	}

	return c.toTicket(updated), nil
}

func (c *Client) toTicket(issue *gh.Issue) *feedback.Ticket {
	t := &feedback.Ticket{
		Number:    issue.GetNumber(),
		Title:     issue.GetTitle(),
		Body:      issue.GetBody(),
		State:     issue.GetState(),
		URL:       issue.GetHTMLURL(),
		CreatedAt: issue.GetCreatedAt().Time,
		UpdatedAt: issue.GetUpdatedAt().Time,
		Comments:  issue.GetComments(),
	}

	for _, label := range issue.Labels {
		name := label.GetName()
		t.Labels = append(t.Labels, name)

		switch {
		case strings.HasPrefix(name, "source:"):
			t.App = feedback.App(strings.TrimPrefix(name, "source:"))
		case strings.HasPrefix(name, "type:"):
			t.Type = feedback.TicketType(strings.TrimPrefix(name, "type:"))
		case strings.HasPrefix(name, "status:"):
			t.Status = feedback.Status(strings.TrimPrefix(name, "status:"))
		}
	}

	if t.Status == "" && t.State == "open" {
		t.Status = feedback.StatusOpen
	}

	return t
}