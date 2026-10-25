package httpapi

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

func (a *API) listChildren(w http.ResponseWriter, r *http.Request) {
	cs, err := a.svc.ListChildren(r.Context(), actorFrom(r.Context()))
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(cs, toChild))
}

func (a *API) createChild(w http.ResponseWriter, r *http.Request) {
	var in struct {
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		BirthDate     string `json:"birth_date"`
		PolicyVersion string `json:"policy_version"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	bd, err := time.Parse(time.DateOnly, in.BirthDate)
	if err != nil {
		fail(w, r, a.log, invalid("birth_date"))
		return
	}
	c, err := a.svc.CreateChild(r.Context(), actorFrom(r.Context()), application.ChildInput{FirstName: in.FirstName,
		LastName: in.LastName, BirthDate: bd, PolicyVersion: in.PolicyVersion})
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusCreated, toChild(c))
}

func (a *API) getChild(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err == nil {
		var c domain.Child
		if c, err = a.svc.GetChild(r.Context(), actorFrom(r.Context()), id); err == nil {
			writeJSON(w, http.StatusOK, toChild(c))
			return
		}
	}
	fail(w, r, a.log, err)
}

func (a *API) patchChild(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var in struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		BirthDate *string `json:"birth_date"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	p := application.ChildPatch{FirstName: in.FirstName, LastName: in.LastName}
	if in.BirthDate != nil {
		bd, err := time.Parse(time.DateOnly, *in.BirthDate)
		if err != nil {
			fail(w, r, a.log, invalid("birth_date"))
			return
		}
		p.BirthDate = &bd
	}
	c, err := a.svc.UpdateChild(r.Context(), actorFrom(r.Context()), id, p)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, toChild(c))
}

func (a *API) deleteChild(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err == nil {
		err = a.svc.DeleteChild(r.Context(), actorFrom(r.Context()), id)
	}
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) linkGuardian(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var in struct {
		GuardianID   uuid.UUID `json:"guardian_id"`
		Relationship string    `json:"relationship"`
		IsPrimary    bool      `json:"is_primary"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	if in.IsPrimary { // exactly one primary guardian per child (invariant)
		fail(w, r, a.log, &domain.ValidationError{Fields: []domain.FieldError{{Field: "is_primary", Code: "primary_exists"}}})
		return
	}
	g, err := a.svc.LinkGuardian(r.Context(), actorFrom(r.Context()), id, domain.ChildGuardian{GuardianID: in.GuardianID, Relationship: domain.Relationship(in.Relationship)})
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"child_id": g.ChildID, "guardian_id": g.GuardianID, "relationship": g.Relationship, "is_primary": g.IsPrimary, "created_at": g.CreatedAt})
}

func assignmentJSON(as domain.ChildSpecialist) map[string]any {
	return map[string]any{"child_id": as.ChildID, "specialist_id": as.SpecialistID, "assigned_by": as.AssignedBy, "assigned_at": as.AssignedAt, "ended_at": as.EndedAt}
}

func (a *API) assign(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var in struct {
		SpecialistID uuid.UUID `json:"specialist_id"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	as, err := a.svc.AssignSpecialist(r.Context(), actorFrom(r.Context()), id, in.SpecialistID)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusCreated, assignmentJSON(as))
}

func (a *API) endAssignment(w http.ResponseWriter, r *http.Request) {
	ids, err := pathIDs(r, "childId", "specialistId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var in struct {
		EndedAt string `json:"ended_at"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	at, err := parseTime("ended_at", in.EndedAt)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	as, err := a.svc.EndAssignment(r.Context(), actorFrom(r.Context()), ids[0], ids[1], at)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, assignmentJSON(as))
}

func (a *API) export(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	ex, err := a.svc.ExportChild(r.Context(), actorFrom(r.Context()), id)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	out := []any{toChild(ex.Child)}
	for _, s := range ex.Sessions {
		out = append(out, toSession(s))
	}
	for _, e := range ex.Entries {
		out = append(out, toEntry(e))
	}
	w.Header().Set("Content-Disposition", `attachment; filename="child-export.json"`)
	writeJSON(w, http.StatusOK, out)
}
