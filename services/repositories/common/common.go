package common

import (
	"fmt"
	"strings"
)

type QueryOrder string

const (
	Ascending  QueryOrder = "ASC"
	Descending            = "DESC"
)

type QueryOptions struct {
	Limit       int
	Page        int
	OffAdjust   int
	Order       QueryOrder
	OrderBy     string
	WithBanned  bool
	WithSpammed bool
	WithDeleted bool
}

type QueryCapabilities struct {
	HasBanned  bool
	HasSpammed bool
	HasDeleted bool
}

func (qo QueryOptions) Query(
	query string,
	qc QueryCapabilities,
) (q string) {
	if strings.Index(query, "WHERE") == -1 {
		q = fmt.Sprintf("%s WHERE ", query)
	} else {
		q = fmt.Sprintf("%s AND ", query)
	}

	var wheres []string

	if qc.HasBanned && qo.WithBanned == false {
		wheres = append(wheres, "banned_at IS NULL")
	}
	if qc.HasSpammed && qo.WithSpammed == false {
		wheres = append(wheres, "spammed_at IS NULL")
	}
	if qc.HasDeleted && qo.WithDeleted == false {
		wheres = append(wheres, "deleted_at IS NULL")
	}

	if len(wheres) > 0 {
		q = fmt.Sprintf("%s %s", q, strings.Join(wheres, " AND "))
	}

	if qo.OrderBy != "" {
		if qo.Order == "" {
			qo.Order = Descending
		}
		q = fmt.Sprintf("%s ORDER BY %s %s", q, qo.OrderBy, string(qo.Order))
	}

	if qo.Limit > 0 {
		q = fmt.Sprintf("%s LIMIT %d", q, qo.Limit)
		if qo.Page > 0 {
			q = fmt.Sprintf("%s OFFSET %d", q, int(((qo.Page-1)*qo.Limit)-qo.OffAdjust))
		}
	}

	return q
}
