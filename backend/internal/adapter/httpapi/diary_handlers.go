package httpapi

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

type sessionBody struct {
	SpecialistID *uuid.UUID `json:"specialist_id"`
	Title        *string    `json:"title"`
	StartsAt     *string    `json:"starts_at"`
	EndsAt       *string    `json:"ends_at"`
	Mode         *string    `json:"mode"`
	Status       *string    `json:"status"`
}

func (a *API) listSessions(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err == nil {
		var ss []domain.SupportSession
		if ss, err = a.svc.ListSessions(r.Context(), actorFrom(r.Context()), id); err == nil {
			writeJSON(w, http.StatusOK, mapSlice(ss, toSession))
			return
		}
	}
	fail(w, r, a.log, err)
}

func (a *API) createSession(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var b sessionBody
	if err := decode(r, &b); err != nil {
		fail(w, r, a.log, err)
		return
	}
	if b.SpecialistID == nil || b.Title == nil || b.StartsAt == nil || b.EndsAt == nil || b.Mode == nil || b.Status == nil {
		fail(w, r, a.log, invalid("body"))
		return
	}
	st, err1 := parseTime("starts_at", *b.StartsAt)
	en, err2 := parseTime("ends_at", *b.EndsAt)
	if err1 != nil || err2 != nil {
		fail(w, r, a.log, invalid("starts_at|ends_at"))
		return
	}
	s, err := a.svc.CreateSession(r.Context(), actorFrom(r.Context()), id, application.SessionInput{SpecialistID: *b.SpecialistID,
		Title: *b.Title, StartsAt: st, EndsAt: en, Mode: domain.SessionMode(*b.Mode), Status: domain.SessionStatus(*b.Status)})
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusCreated, toSession(s))
}

func (a *API) patchSession(w http.ResponseWriter, r *http.Request) {
	ids, err := pathIDs(r, "childId", "sessionId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var b sessionBody
	if err := decode(r, &b); err != nil || b.SpecialistID != nil {
		if err == nil {
			err = invalid("specialist_id")
		}
		fail(w, r, a.log, err)
		return
	}
	p := application.SessionPatch{Title: b.Title}
	if b.StartsAt != nil {
		t, err := parseTime("starts_at", *b.StartsAt)
		if err != nil {
			fail(w, r, a.log, err)
			return
		}
		p.StartsAt = &t
	}
	if b.EndsAt != nil {
		t, err := parseTime("ends_at", *b.EndsAt)
		if err != nil {
			fail(w, r, a.log, err)
			return
		}
		p.EndsAt = &t
	}
	if b.Mode != nil {
		m := domain.SessionMode(*b.Mode)
		p.Mode = &m
	}
	if b.Status != nil {
		s := domain.SessionStatus(*b.Status)
		p.Status = &s
	}
	s, err := a.svc.UpdateSession(r.Context(), actorFrom(r.Context()), ids[0], ids[1], p)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, toSession(s))
}

type entryBody struct {
	Type       string `json:"type"`
	Visibility string `json:"visibility"`
	Content    string `json:"content"`
	ObservedAt string `json:"observed_at"`
}

func (b entryBody) input() (application.EntryInput, error) {
	t, err := parseTime("observed_at", b.ObservedAt)
	return application.EntryInput{Type: domain.EntryType(b.Type), Visibility: domain.Visibility(b.Visibility), Content: b.Content, ObservedAt: t}, err
}

func (a *API) listEntries(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	q := r.URL.Query()
	limit, err := parseLimit(q.Get("limit"))
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	cur, err := decodeCursor(q.Get("cursor"))
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	vis := domain.Visibility(q.Get("visibility"))
	if vis != "" && !vis.Valid() || len(q.Get("q")) > 100 {
		fail(w, r, a.log, invalid("query"))
		return
	}
	p, err := a.svc.ListEntries(r.Context(), actorFrom(r.Context()), id, application.ListQuery{Query: q.Get("q"), Visibility: vis, Limit: limit, After: cur})
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	if p.Next != nil {
		w.Header().Set("X-Next-Cursor", encodeCursor(p.Next))
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(p.Total))
	writeJSON(w, http.StatusOK, mapSlice(p.Items, toEntry))
}

func (a *API) createEntry(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "childId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var b entryBody
	if err := decode(r, &b); err != nil {
		fail(w, r, a.log, err)
		return
	}
	in, _ := b.input() // zero observed_at is reported by the use case after access checks
	e, err := a.svc.CreateEntry(r.Context(), actorFrom(r.Context()), id, in)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusCreated, toEntry(e))
}

func (a *API) editEntry(w http.ResponseWriter, r *http.Request) {
	ids, err := pathIDs(r, "childId", "entryId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	var b entryBody
	if err := decode(r, &b); err != nil {
		fail(w, r, a.log, err)
		return
	}
	in, _ := b.input()
	e, err := a.svc.EditEntry(r.Context(), actorFrom(r.Context()), ids[0], ids[1], in)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, toEntry(e))
}

func (a *API) listVersions(w http.ResponseWriter, r *http.Request) {
	ids, err := pathIDs(r, "childId", "entryId")
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	vs, err := a.svc.ListVersions(r.Context(), actorFrom(r.Context()), ids[0], ids[1])
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(vs, func(v domain.DiaryEntryVersion) versionDTO {
		return versionDTO{v.ID, v.EntryID, v.VersionNumber, v.ChangedBy, v.Type, v.Visibility, v.Content, v.ObservedAt, v.ChangedAt}
	}))
}
