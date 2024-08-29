package mem

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func iterate[K comparable, V any](ks []K, vs []V) iter.Seq2[K, V] {
	if len(ks) != len(vs) {
		panic("ks and vs has different lengths")
	}
	return func(yield func(k K, v V) bool) {
		for i := range len(ks) {
			if !yield(ks[i], vs[i]) {
				return
			}
		}
	}
}

func TestMerge_2(t *testing.T) {
	type args struct {
		a iter.Seq2[string, int]
		b iter.Seq2[string, int]
	}
	tests := []struct {
		name string
		args args
		want iter.Seq2[string, int]
	}{
		{
			name: "a and b lengths are equal",
			args: args{
				a: iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				b: iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
			},
			want: iterate([]string{"a", "a", "b", "b", "c", "c"}, []int{1, 1, 2, 2, 3, 3}),
		},
		{
			name: "a shorter than b",
			args: args{
				a: iterate([]string{"a"}, []int{1}),
				b: iterate([]string{"b", "c"}, []int{2, 3}),
			},
			want: iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
		},
		{
			name: "b shorter than a",
			args: args{
				a: iterate([]string{"b", "c"}, []int{2, 3}),
				b: iterate([]string{"a"}, []int{1}),
			},
			want: iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
		},
		{
			name: "a is empty",
			args: args{
				a: iterate([]string{}, []int{}),
				b: iterate([]string{"b", "c"}, []int{2, 3}),
			},
			want: iterate([]string{"b", "c"}, []int{2, 3}),
		},
		{
			name: "b is empty",
			args: args{
				a: iterate([]string{"a"}, []int{1}),
				b: iterate([]string{}, []int{}),
			},
			want: iterate([]string{"a"}, []int{1}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextWant, stopWant := iter.Pull2(tt.want)
			defer stopWant()

			gotKeys := []string{}

			for k, v := range Merge(tt.args.a, tt.args.b) {
				wantK, wantV, ok := nextWant()
				require.True(t, ok)

				assert.Equal(t, wantK, k)
				assert.Equal(t, wantV, v)

				gotKeys = append(gotKeys, k)
			}

			assert.IsNonDecreasing(t, gotKeys)
		})
	}
}

func TestMerge_2_panics(t *testing.T) {
	type args struct {
		seqs []iter.Seq2[string, int]
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "a has decreasing order",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"c", "a", "b"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
			}},
			want: "sequences must be ordered!",
		},
		{
			name: "a decreasing from the middle",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a", "b", "d", "c"}, []int{1, 2, 3, 4}),
				iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
			}},
			want: "sequences must be ordered!",
		},
		{
			name: "a has less value at the end",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a", "b", "c", "a"}, []int{1, 2, 3, 4}),
				iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
			}},
			want: "sequences must be ordered!",
		},
		{
			name: "b has decreasing order",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"c", "a", "b"}, []int{1, 2, 3}),
			}},
			want: "sequences must be ordered!",
		},
		{
			name: "b decreasing from the middle",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
				iterate([]string{"a", "b", "d", "c"}, []int{1, 2, 3, 4}),
			}},
			want: "sequences must be ordered!",
		},
		{
			name: "b has less value at the end",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
				iterate([]string{"a", "b", "c", "a"}, []int{1, 2, 3, 4}),
			}},
			want: "sequences must be ordered!",
		},
		{
			name: "1st has decreasing order",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"c", "a", "b"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
			}},
			want: "sequences must be ordered!",
		},
		{
			name: "1st decreasing from the middle",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a", "b", "d", "c"}, []int{1, 2, 3, 4}),
				iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
				iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
			}},
			want: "sequences must be ordered!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.PanicsWithValue(t, tt.want, func() {
				for range Merge(tt.args.seqs[0], tt.args.seqs[1], tt.args.seqs[2:]...) {
				}
			})
		})
	}
}

func TestMerge_N(t *testing.T) {
	type args struct {
		seqs []iter.Seq2[string, int]
	}
	tests := []struct {
		name     string
		args     args
		wantKeys []string
		wantVals []int
	}{
		{
			name: "N = 3 lengths are equal",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
			}},
			wantKeys: []string{"a", "a", "a", "b", "b", "b", "c", "c", "c"},
			wantVals: []int{1, 1, 1, 2, 2, 2, 3, 3, 3},
		},
		{
			name: "N = 4 lengths are equal",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
			}},
			wantKeys: []string{"a", "a", "a", "a", "b", "b", "b", "b", "c", "c", "c", "c"},
			wantVals: []int{1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3},
		},
		{
			name: "N = 4 increasing lengths order",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a"}, []int{1}),
				iterate([]string{"a", "b"}, []int{1, 2}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
			}},
			wantKeys: []string{"a", "a", "a", "a", "b", "b", "b", "c", "c", "d"},
			wantVals: []int{1, 1, 1, 1, 2, 2, 2, 3, 3, 4},
		},
		{
			name: "N = 4 decreasing lengths order",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b"}, []int{1, 2}),
				iterate([]string{"a"}, []int{1}),
			}},
			wantKeys: []string{"a", "a", "a", "a", "b", "b", "b", "c", "c", "d"},
			wantVals: []int{1, 1, 1, 1, 2, 2, 2, 3, 3, 4},
		},
		{
			name: "N = 4 shuffled lengths order",
			args: args{[]iter.Seq2[string, int]{
				iterate([]string{"a"}, []int{1}),
				iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				iterate([]string{"a", "b"}, []int{1, 2}),
				iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
			}},
			wantKeys: []string{"a", "a", "a", "a", "b", "b", "b", "c", "c", "d"},
			wantVals: []int{1, 1, 1, 1, 2, 2, 2, 3, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeys := []string{}
			gotVals := []int{}

			for k, v := range Merge(tt.args.seqs[0], tt.args.seqs[1], tt.args.seqs[2:]...) {
				gotKeys = append(gotKeys, k)
				gotVals = append(gotVals, v)
			}

			if assert.Equal(t, tt.wantKeys, gotKeys) && assert.Equal(t, tt.wantVals, gotVals) {
				assert.IsNonDecreasing(t, gotKeys)
			}
		})
	}
}
