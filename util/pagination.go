package util

const (
	DefaultOffset = 0
	DefaultLimit  = 1000
)

// OffsetLimit is used to retrieve the records in the db by page.
type OffsetLimit struct {
	Offset int
	Limit  int
}

// NewOffsetLimit uses offset/DefaultOffset and limit/DefaultLimit
// to generate an OffsetLimit.
func NewOffsetLimit(offset *int, limit *int) *OffsetLimit {
	var o, l int = DefaultOffset, DefaultLimit

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
