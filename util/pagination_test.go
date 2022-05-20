package util

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func Test_NewOffsetLimit(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		offset *int
		limit  *int
	}
	tests := []struct {
		name string
		args args
		want *OffsetLimit
	}{
		{
			name: "both offset and limit are not nil",
			args: args{
				offset: pointer.ToInt(0),
				limit:  pointer.ToInt(10),
			},
			want: &OffsetLimit{
				Offset: 0,
				Limit:  10,
			},
		},
		{
			name: "both offset and limit are nil",
			want: &OffsetLimit{
				Offset: 0,
				Limit:  1000,
			},
		},
		{
			name: "offset not nil and limit nil",
			args: args{
				offset: pointer.ToInt(2),
			},
			want: &OffsetLimit{
				Offset: 2,
				Limit:  1000,
			},
		},
		{
			name: "offset nil and limit not nil",
			args: args{
				limit: pointer.ToInt(10),
			},
			want: &OffsetLimit{
				Offset: 0,
				Limit:  10,
			},
		},
	}
	for _, tt := range tests {
		assert.Equal(tt.want, NewOffsetLimit(tt.args.offset, tt.args.limit), tt.name)
	}
}
