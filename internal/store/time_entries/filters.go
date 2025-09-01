package time_entries

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidCategoryId    = fmt.Errorf("invalid category id")
	ErrInvalidUserId        = fmt.Errorf("invalid user id")
	ErrInvalidFromDate      = fmt.Errorf("invalid from date")
	ErrInvalidToDate        = fmt.Errorf("invalid to date")
	ErrFromDateAfterToDate  = fmt.Errorf("from date cannot be after to date")
	ErrToDateBeforeFromDate = fmt.Errorf("to date cannot be before from date")
)

type Filters struct {
	Query      string     `json:"query"`
	CategoryId []string   `json:"categoryId"`
	UserId     []string   `json:"userId"`
	FromDate   *time.Time `json:"fromDate"`
	ToDate     *time.Time `json:"toDate"`
}

func (f *Filters) Parse(r *http.Request) error {
	p := r.URL.Query()

	if p.Has("query") && p.Get("query") != "" {
		f.Query = strings.TrimSpace(p.Get("query"))
	}

	if p.Has("categoryId") && p.Get("categoryId") != "" {
		idsQuery := strings.Split(p.Get("categoryId"), ",")
		ids := []string{}
		for _, id := range idsQuery {
			if id == "" {
				continue
			}
			ids = append(ids, id)
		}
		f.CategoryId = ids
	}

	if p.Has("userId") && p.Get("userId") != "" {
		idsQuery := strings.Split(p.Get("userId"), ",")
		ids := []string{}
		for _, id := range idsQuery {
			if id == "" {
				continue
			}
			ids = append(ids, id)
		}
		f.UserId = ids
	}

	if p.Has("fromDate") && p.Get("fromDate") != "" {
		parsed, err := time.Parse(time.DateOnly, p.Get("fromDate"))
		if err != nil {
			return ErrInvalidFromDate
		}
		f.FromDate = &parsed
	}

	if p.Has("toDate") && p.Get("toDate") != "" {
		parsed, err := time.Parse(time.DateOnly, p.Get("toDate"))
		if err != nil {
			return ErrInvalidToDate
		}
		inclusive := parsed.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		f.ToDate = &inclusive
	}

	if f.FromDate != nil && f.ToDate != nil {
		if f.FromDate.After(*f.ToDate) {
			return ErrFromDateAfterToDate
		}
		if f.ToDate.Before(*f.FromDate) {
			return ErrToDateBeforeFromDate
		}
	}

	return nil
}
