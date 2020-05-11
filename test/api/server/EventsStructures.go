package server

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type EventSerializer struct {
	*Event
}

type EventListSerializer struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	DateStart time.Time `json:"date_start"`
	DateFin time.Time `json:"date_fin"`
	PlaceName string `json:"place_name"`
}

func EventListSerializerFunc(events []*Event) []render.Renderer {
	var list []render.Renderer
	for _, event := range events {
		list = append(list, &EventListSerializer{
			ID: event.ID,
			Name: event.Name,
			DateStart: event.DateStart,
			DateFin: event.DateFin,
			PlaceName: event.PlaceName,
		})
	}
	return list
}

func (a *EventSerializer) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *EventListSerializer) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}


func ParsePostEvent(data url.Values) (Event, []User, []Company, error) {
	var a Event
	var err error
	var manager User
	var members []User
	var companies []Company
	var manId int
	var membersIds []int
	var membersIdsString []string
	var companiesIds []int
	var companiesIdsString []string

	const layout = "2006-01-02T15:04:05.000Z"

	if data["name"] == nil ||
		data["description"] == nil ||
		data["datestart"] == nil ||
		data["datefin"] == nil ||
		data["placename"] == nil ||
		data["placelat"] == nil ||
		data["placelon"] == nil ||
		data["manager"] == nil ||
		data["members"] == nil ||
		data["companies"] == nil {
		return a, nil, nil, errors.New("missing required Event fields")
	}

	if len(data["name"]) != 1 ||
		len(data["description"]) != 1 ||
		len(data["datestart"]) != 1 ||
		len(data["datefin"]) != 1 ||
		len(data["placename"]) != 1 ||
		len(data["placelat"]) != 1 ||
		len(data["placelon"]) != 1 ||
		len(data["manager"]) != 1 ||
		len(data["members"]) != 1 ||
		len(data["companies"]) != 1 {
		return a, nil, nil, errors.New("multiple forms submitted")
	}

	a.Name = data["name"][0]
	a.Description = data["description"][0]
	a.PlaceName = data["placename"][0]

	a.DateStart, err = time.Parse(layout, data["datestart"][0])
	if err != nil {
		return a, nil, nil, errors.New("invalid values submitted0")
	}

	a.DateFin, err = time.Parse(layout, data["datefin"][0])
	if err != nil {
		return a, nil, nil, errors.New("invalid values submitted1")
	}

	a.PlaceLat, err = strconv.ParseFloat(data["placelat"][0], 64)
	if err != nil {
		return a, nil, nil, errors.New("invalid values submitted2")
	}

	a.PlaceLon, err = strconv.ParseFloat(data["placelon"][0], 64)
	if err != nil {
		return a, nil, nil, errors.New("invalid values submitted3")
	}

	manId, err = strconv.Atoi(data["manager"][0])
	if err != nil {
		return a, nil, nil, errors.New("invalid values submitted4")
	}
	err = db.First(&manager, manId).Error
	if err != nil {
		return a, nil, nil, errors.New("manager user does not exist5")
	}
	a.Manager = manager

	membersIdsString = strings.Fields(data["members"][0])
	membersIds, err = StringArrToInt(membersIdsString)
	if err != nil {
		return a, nil, nil, errors.New("invalid values submitted6")
	}
	err = db.Where("id IN (?)", membersIds).Find(&members).Error
	if err != nil {
		return a, nil, nil, errors.New("member user do not exist")
	}


	companiesIdsString = strings.Fields(data["companies"][0])
	companiesIds, err = StringArrToInt(companiesIdsString)
	if err != nil {
		return a, nil, nil, errors.New("invalid values submitted7")
	}
	err = db.Where("id IN (?)", companiesIds).Find(&companies).Error
	if err != nil {
		return a, nil, nil, errors.New("copmanies do not exist")
	}

	return a, members, companies, nil
}



func ParseEventPutData(data url.Values, event *Event) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	var err error
	var manager User
	var companies []Company
	var companiesIds []int
	var companiesIdsString []string
	var manId int
	var members []User
	var membersIds []int
	var membersIdsString []string
	const layout = "2006-01-02T15:04:05.000Z"

	for k, v := range data {
		if k == "date_start" || k == "date_fin" {
			res[k], err = time.Parse(layout, data[k][0])
			if err != nil {
				return res, errors.New("invalid values submitted1")
			}
		} else if k == "manager" {

			manId, err = strconv.Atoi(data["manager"][0])
			if err != nil {
				return res, errors.New("invalid values submitted4")
			}
			err = db.First(&manager, manId).Error
			if err != nil {
				return res, errors.New("manager user does not exist5")
			}
			res[k] = manager

		} else if k == "companies" {


			companiesIdsString = data[k]
			companiesIds, err = StringArrToInt(companiesIdsString)
			if err != nil {
				return res, errors.New("invalid values submitted7")
			}
			err = db.Where("id IN (?)", companiesIds).Find(&companies).Error
			if err != nil {
				return res, errors.New("company do not exist")
			}
			err = db.Model(&event).Association("Companies").Clear().Error
			err = db.Model(&event).Association("Companies").Append(companies).Error
			if err != nil {
				return res, errors.New("company does not exist")
			}

		} else if k == "members" {


			membersIdsString = strings.Fields(data["members"][0])
			membersIds, err = StringArrToInt(membersIdsString)
			if err != nil {
				return res, errors.New("invalid values submitted6")
			}
			err = db.Where("id IN (?)", membersIds).Find(&members).Error
			if err != nil {
				return res, errors.New("member user do not exist")
			}

			err = db.Model(&event).Association("Members").Clear().Error
			err = db.Model(&event).Association("Members").Append(members).Error
			if err != nil {
				return res, errors.New("member does not exist")
			}

		} else {
			res[k] = v[0]
		}

	}

	return res, nil
}