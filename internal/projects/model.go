package projects

import "time"

// Project captures a collaborative workspace owned by a user.
type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	Members   []string  `json:"members"`
	CreatedAt time.Time `json:"created_at"`
}
