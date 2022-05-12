package db

const (
	DefaultOffset = 0
	DefaultLimit  = 1000
)

// OffsetLimit is used to retrieve the records in the db by page.
type OffsetLimit struct {
	Offset int64
	Limit  int64
}

// NewOffsetLimit uses offset/DefaultOffset and limit/DefaultLimit
// to generate an OffsetLimit.
func NewOffsetLimit(offset *int64, limit *int64) *OffsetLimit {
	var o, l int64 = DefaultOffset, DefaultLimit

	if offset != nil {
		o = *offset
	}

	if limit != nil {
		l = *limit
	}

	return &OffsetLimit{
		Offset: o,
		Limit:  l,
	}
}
