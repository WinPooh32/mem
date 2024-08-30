package mem

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	type args struct {
		a iter.Seq2[string, int]
		b iter.Seq2[string, float32]
	}
	tests := []struct {
		name     string
		args     args
		wantKeys []string
		wantVals []Tuple[int, float32]
	}{
		{
			name: "all joined",
			args: args{
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c"}, []float32{-1.0, -2.0, -3.0}),
			},
			wantKeys: []string{"a", "b", "c"},
			wantVals: []Tuple[int, float32]{
				{1, []float32{-1.0}},
				{2, []float32{-2.0}},
				{3, []float32{-3.0}},
			},
		},
		{
			name: "all joined with miltiple right values",
			args: args{
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c", "b"}, []float32{-1.0, -2.0, -3.0, -4.0}),
			},
			wantKeys: []string{"a", "b", "c"},
			wantVals: []Tuple[int, float32]{
				{1, []float32{-1.0}},
				{2, []float32{-2.0, -4}},
				{3, []float32{-3.0}},
			},
		},
		{
			name: "nothing joined",
			args: args{
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"d", "e", "f", "g"}, []float32{-1.0, -2.0, -3.0, -4.0}),
			},
			wantKeys: []string{"a", "b", "c"},
			wantVals: []Tuple[int, float32]{
				{1, nil},
				{2, nil},
				{3, nil},
			},
		},
		{
			name: "some joined",
			args: args{
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"d", "e", "f", "a"}, []float32{-1.0, -2.0, -3.0, -4.0}),
			},
			wantKeys: []string{"a", "b", "c"},
			wantVals: []Tuple[int, float32]{
				{1, []float32{-4.0}},
				{2, nil},
				{3, nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			gotKeys := []string{}
			gotVals := []Tuple[int, float32]{}

			for k, v := range Lookup(tt.args.a, tt.args.b) {
				gotKeys = append(gotKeys, k)
				gotVals = append(gotVals, v)
			}

			assert.Equal(tt.wantKeys, gotKeys)
			assert.Equal(tt.wantVals, gotVals)
		})
	}
}
