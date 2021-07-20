package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math"
	"strings"
)

type Paging struct {
	Slice       []interface{} `json:"slice"`
	Count       int           `json:"count"`
	Page        int           `json:"page"`
	PageCount   int           `json:"pageCount"`
	RowsPerPage int           `json:"rowsPerPage"`
}

func (p *Paging) Do(
	getCount func() (int, error),
	getSlice func() ([]interface{}, error),
	page int,
) (*Paging, error) {
	p.Page = page
	p.PageCount = 1
	p.RowsPerPage = 15

	count, err := getCount()
	if err != nil {
		return nil, err
	}

	p.Count = count
	p.Slice = make([]interface{}, 0)

	if p.Count == 0 {
		return p, nil
	}

	if p.Page == 0 {
		slice, err := getSlice()
		if err != nil {
			return nil, err
		}
		p.Slice = slice
	} else {
		p.PageCount = int(math.Ceil(float64(p.Count) / float64(p.RowsPerPage)))
		if p.Page > p.PageCount {
			p.Page = p.PageCount
		}

		slice, err := getSlice()
		if err != nil {
			return nil, err
		}
		p.Slice = slice
	}
	return p, nil
}

func (p *Paging) DoQuery(query *mgo.Query, page int, order string) (*Paging, error) {
	return p.Do(
		func() (int, error) {
			return query.Count()
		},
		func() ([]interface{}, error) {
			var slice []interface{}

			if p.Page == 0 {
				p.Page = 1
				if err := query.Skip(0).Limit(p.Count).Sort(order).All(&slice); err != nil {
					return nil, err
				}
			} else {
				if err := query.Skip((p.Page - 1) * p.RowsPerPage).Limit(p.RowsPerPage).Sort(order).All(&slice); err != nil {
					return nil, err
				}
			}
			return slice, nil
		},
		page,
	)
}

func (p *Paging) DoPipe(col *mgo.Collection, stages []bson.M, page int, order string) (*Paging, error) {
	return p.Do(
		func() (int, error) {
			var countMap map[string]int

			err := col.Pipe(append(stages,
				bson.M{"$group": bson.M{"_id": "null", "count": bson.M{"$sum": 1}}},
				bson.M{"$project": bson.M{"_id": 0}},
			)).One(&countMap)
			if err != nil {
				return 0, err
			}
			return countMap["count"], nil
		},
		func() ([]interface{}, error) {
			var slice []interface{}

			if order != "" {
				stages = append(stages, bson.M{"$sort": ToSort(order)})
			}

			if p.Page == 0 {
				p.Page = 1
				err := col.Pipe(stages).All(&slice)
				if err != nil {
					return nil, err
				}
			} else {
				stages = append(stages,
					bson.M{"$skip": (p.Page - 1) * p.RowsPerPage},
					bson.M{"$limit": p.RowsPerPage},
				)
				err := col.Pipe(stages).All(&slice)
				if err != nil {
					return nil, err
				}
			}
			return slice, nil
		},
		page,
	)
}

func ToSort(order string) bson.M {
	if strings.HasPrefix(order, "-") {
		return bson.M{order[1:]: -1}
	} else {
		return bson.M{order: 1}
	}
}
