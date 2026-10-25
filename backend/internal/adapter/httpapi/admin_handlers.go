package httpapi

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		FirstName          string `json:"first_name"`
		LastName           string `json:"last_name"`
		Role               string `json:"role"`
		SpecialistApproved bool   `json:"specialist_approved"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	u, err := a.svc.CreateUserByAdmin(r.Context(), actorFrom(r.Context()), application.RegisterInput{Email: in.Email,
		Password: in.Password, FirstName: in.FirstName, LastName: in.LastName, Role: domain.Role(in.Role)}, in.SpecialistApproved)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusCreated, toUser(u))
}

func (a *API) patchUser(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "userId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var in struct {
		Status             *domain.UserStatus `json:"status"`
		SpecialistApproved *bool              `json:"specialist_approved"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	u, err := a.svc.UpdateUser(r.Context(), actorFrom(r.Context()), id, application.UserPatch{Status: in.Status, SpecialistApproved: in.SpecialistApproved})
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, toUser(u))
}

func (a *API) listUsers(w http.ResponseWriter, r *http.Request) {
	limit, err := parseLimit(r.URL.Query().Get("limit"))
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var after *uuid.UUID
	if c := r.URL.Query().Get("cursor"); c != "" {
		id, err := uuid.Parse(c)
		if err != nil {
			fail(w, r, a.log, invalid("cursor"))
			return
		}
		after = &id
	}
	us, err := a.svc.ListUsers(r.Context(), actorFrom(r.Context()), limit, after)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	if len(us) == limit {
		w.Header().Set("X-Next-Cursor", us[len(us)-1].ID.String())
	}
	writeJSON(w, http.StatusOK, mapSlice(us, toUser))
}

func toAudit(e domain.AuditEvent) auditDTO {
	return auditDTO{e.ID, e.ActorID, e.ChildID, e.Action, e.Resource, e.ResourceID, e.Outcome, e.RequestID, e.OccurredAt}
}

func (a *API) auditList(w http.ResponseWriter, r *http.Request, fetch func(limit int, cur *application.Cursor) (application.AuditPage, error)) {
	limit, err := parseLimit(r.URL.Query().Get("limit"))
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	cur, err := decodeCursor(r.URL.Query().Get("cursor"))
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	p, err := fetch(limit, cur)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	if p.Next != nil {
		w.Header().Set("X-Next-Cursor", encodeCursor(p.Next))
	}
	writeJSON(w, http.StatusOK, mapSlice(p.Items, toAudit))
}

func (a *API) systemAudit(w http.ResponseWriter, r *http.Request) {
	a.auditList(w, r, func(l int, c *application.Cursor) (application.AuditPage, error) {
		return a.svc.SystemAudit(r.Context(), actorFrom(r.Context()), l, c)
	})
}

func (a *API) childAudit(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	a.auditList(w, r, func(l int, c *application.Cursor) (application.AuditPage, error) {
		return a.svc.ChildAudit(r.Context(), actorFrom(r.Context()), id, l, c)
	})
}
