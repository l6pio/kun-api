package db

import (
	"database/sql"
	"math"
)

type Paging struct {
	Slice       []interface{}                               `json:"slice"`
	Count       int                                         `json:"count"`
	Page        int                                         `json:"page"`
	PageCount   int                                         `json:"pageCount"`
	RowsPerPage int                                         `json:"rowsPerPage"`
	DoCount     func() (*sql.Rows, error)                   `json:"-"`
	DoQuery     func(from int, size int) (*sql.Rows, error) `json:"-"`
	Convert     func(rows *sql.Rows) []interface{}          `json:"-"`
}

func (p *Paging) Do() (*Paging, error) {
	p.PageCount = 1
	p.RowsPerPage = 15

	countRows, err := p.DoCount()
	if err != nil {
		return nil, err
	}
	defer countRows.Close()

	var count int
	countRows.Next()
	if err := countRows.Scan(&count); err != nil {
		return nil, err
	}

	if err := countRows.Err(); err != nil {
		return nil, err
	}

	p.Count = count
	p.Slice = make([]interface{}, 0)

	if p.Count == 0 {
		return p, nil
	}

	if p.Page == 0 {
		dataRows, err := p.DoQuery(0, p.Count)
		if err != nil {
			return nil, err
		}
		defer dataRows.Close()

		p.Slice = p.Convert(dataRows)

		if err := dataRows.Err(); err != nil {
			return nil, err
		}
	} else {
		p.PageCount = int(math.Ceil(float64(p.Count) / float64(p.RowsPerPage)))
		if p.Page > p.PageCount {
			p.Page = p.PageCount
		}

		dataRows, err := p.DoQuery((p.Page-1)*p.RowsPerPage, p.RowsPerPage)
		if err != nil {
			return nil, err
		}
		defer dataRows.Close()

		p.Slice = p.Convert(dataRows)

		if err := dataRows.Err(); err != nil {
			return nil, err
		}
	}
	return p, nil
}
