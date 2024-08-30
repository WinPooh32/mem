package mem

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	type sIn struct {
		a, b, c int
	}
	type sOut struct {
		a, b int
	}
	type args struct {
		a   iter.Seq2[string, sIn]
		prj func(string, sIn) (string, sOut)
	}
	tests := []struct {
		name     string
		args     args
		wantKeys []string
		wantVals []sOut
	}{
		{
			name: "projection of a, b fields",
			args: args{
				a:   iterate([]string{"a", "b", "c"}, []sIn{{1, 2, 3}, {3, 4, 5}, {6, 7, 8}}),
				prj: func(k string, s sIn) (string, sOut) { return "_" + k, sOut{s.a, s.b} },
			},
			wantKeys: []string{"_a", "_b", "_c"},
			wantVals: []sOut{{1, 2}, {3, 4}, {6, 7}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			gotKeys := []string{}
			gotVals := []sOut{}

			for k, v := range Project(tt.args.a, tt.args.prj) {
				gotKeys = append(gotKeys, k)
				gotVals = append(gotVals, v)
			}

			assert.Equal(tt.wantKeys, gotKeys)
			assert.Equal(tt.wantVals, gotVals)
		})
	}
}

func TestProject_panics_nil_prj(t *testing.T) {
	a := iterate([]string{"a", "b", "c"}, []int{1, 2, 3})

	f := func() {
		_ = Project[string, int, string, int](a, nil)
	}

	assert.PanicsWithValue(t, "prj is nil!", f)
}
