package db

import (
	"gopkg.in/mgo.v2"
	"math"
)

type Paging struct {
	Slice       []interface{} `json:"slice"`
	Count       int           `json:"count"`
	Page        int           `json:"page"`
	PageCount   int           `json:"pageCount"`
	RowsPerPage int           `json:"rowsPerPage"`
}

func (p *Paging) Do(query *mgo.Query, page int, order string) (*Paging, error) {
	p.Page = page
	p.PageCount = 1
	p.RowsPerPage = 15

	count, err := query.Count()
	if err != nil {
		return nil, err
	}

	p.Count = count
	p.Slice = make([]interface{}, 0)

	if p.Count == 0 {
		return p, nil
	}

	if p.Page == 0 {
		var slice []interface{}
		if err := query.Skip(0).Limit(p.Count).Sort(order).All(&slice); err != nil {
			return nil, err
		}
		p.Slice = slice
	} else {
		p.PageCount = int(math.Ceil(float64(p.Count) / float64(p.RowsPerPage)))
		if p.Page > p.PageCount {
			p.Page = p.PageCount
		}

		var slice []interface{}
		if err := query.Skip((p.Page - 1) * p.RowsPerPage).Limit(p.RowsPerPage).Sort(order).All(&slice); err != nil {
			return nil, err
		}
		p.Slice = slice
	}
	return p, nil
}
