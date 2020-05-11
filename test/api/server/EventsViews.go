package server

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
)

func dbGetEvent(id int) (*Event, error) {
	event := &Event{}
	var err error
	err = db.Preload("Companies").Preload("Members").Preload("Manager").First(event, id).Error
	if err != nil {
		return event, errors.New("event matching query does not exist")
	}
	return event, err
}

func GetEventData(w http.ResponseWriter, r *http.Request) {
	event := r.Context().Value("event").(*Event)
	err := render.Render(w, r, &EventSerializer{Event: event})
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func GetEventList(w http.ResponseWriter, r *http.Request) {
	var events []*Event
	db.Find(&events)
	if err := render.RenderList(w, r, EventListSerializerFunc(events)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	var members []User
	var companies []Company
	err := r.ParseForm()
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	event, members, companies, err = ParsePostEvent(r.Form)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	err = db.Create(&event).Error
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	err = db.Model(&event).Association("Members").Append(members).Error
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	err = db.Model(&event).Association("Companies").Append(companies).Error
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, SuccessfulyCreated())
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	event := r.Context().Value("event").(*Event)
	var data map[string]interface{}

	err := r.ParseForm()
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	data, err = ParseEventPutData(r.Form, event)
	err = db.Model(&event).Updates(data).Error
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, SuccessfulyUpdated())
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	event := r.Context().Value("event").(*Event)

	err := db.Delete(&event).Error
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, SuccessfulyDeleted())
}