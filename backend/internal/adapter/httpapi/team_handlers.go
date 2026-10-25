package httpapi

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

type teamDTO struct {
	UserID       uuid.UUID `json:"user_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `json:"role"`
	IsPrimary    bool      `json:"is_primary"`
	Relationship string    `json:"relationship,omitempty"`
	AssignedAt   any       `json:"assigned_at,omitempty"`
	EndedAt      any       `json:"ended_at,omitempty"`
}

func (a *API) team(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	ms, err := a.svc.Team(r.Context(), actorFrom(r.Context()), id)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(ms, func(m application.TeamMember) teamDTO {
		d := teamDTO{UserID: m.UserID, FirstName: m.FirstName, LastName: m.LastName, Role: string(m.Role), IsPrimary: m.IsPrimary, Relationship: string(m.Relationship)}
		if m.AssignedAt != nil {
			d.AssignedAt = m.AssignedAt
		}
		if m.EndedAt != nil {
			d.EndedAt = m.EndedAt
		}
		return d
	}))
}

func (a *API) specialists(w http.ResponseWriter, r *http.Request) {
	us, err := a.svc.Specialists(r.Context(), actorFrom(r.Context()))
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(us, func(u domain.User) map[string]any {
		return map[string]any{"id": u.ID, "first_name": u.FirstName, "last_name": u.LastName}
	}))
}
