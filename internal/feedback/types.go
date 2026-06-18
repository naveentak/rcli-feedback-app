package feedback

import "time"

type App string

const (
	AppRClip      App = "rclip"
	AppBoKa       App = "boka"
	AppThxBud     App = "thxbud"
	AppMamZo      App = "mamzo"
	AppGlassCourt App = "glasscourt"
)

type TicketType string

const (
	TypeBug            TicketType = "bug"
	TypeFeatureRequest TicketType = "feature-request"
)

type Status string

const (
	StatusOpen       Status = "open"
	StatusTriaged    Status = "triaged"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

var ValidApps = []App{AppRClip, AppBoKa, AppThxBud, AppMamZo, AppGlassCourt}
var ValidTypes = []TicketType{TypeBug, TypeFeatureRequest}
var ValidStatuses = []Status{StatusOpen, StatusTriaged, StatusInProgress, StatusDone}

type SubmitRequest struct {
	App         App        `json:"app" binding:"required"`
	Type        TicketType `json:"type" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Reporter    string     `json:"reporter,omitempty"`
}

type Ticket struct {
	Number      int        `json:"number"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	State       string     `json:"state"`
	URL         string     `json:"url"`
	App         App        `json:"app"`
	Type        TicketType `json:"type"`
	Status      Status     `json:"status"`
	Labels      []string   `json:"labels"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Comments    int        `json:"comments"`
}

type ListFilter struct {
	App    string
	Status string
	Type   string
	State  string // open, closed, all
}

func (a App) Label() string {
	return "source:" + string(a)
}

func (t TicketType) Label() string {
	return "type:" + string(t)
}

func (s Status) Label() string {
	return "status:" + string(s)
}

func IsValidApp(app string) bool {
	for _, a := range ValidApps {
		if string(a) == app {
			return true
		}
	}
	return false
}

func IsValidType(t string) bool {
	for _, vt := range ValidTypes {
		if string(vt) == t {
			return true
		}
	}
	return false
}

func IsValidStatus(s string) bool {
	for _, vs := range ValidStatuses {
		if string(vs) == s {
			return true
		}
	}
	return false
}