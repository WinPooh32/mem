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
		a iter.Seq2[string, int]
		b iter.Seq2[string, int]
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "a has decreasing order",
			args: args{
				a: iterate([]string{"c", "a", "b"}, []int{1, 2, 3}),
				b: iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
			},
			want: "sequences must be ordered!",
		},
		{
			name: "a decreasing from the middle",
			args: args{
				a: iterate([]string{"a", "b", "d", "c"}, []int{1, 2, 3, 4}),
				b: iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
			},
			want: "sequences must be ordered!",
		},
		{
			name: "a has less value at the end",
			args: args{
				a: iterate([]string{"a", "b", "c", "a"}, []int{1, 2, 3, 4}),
				b: iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
			},
			want: "sequences must be ordered!",
		},
		{
			name: "b has decreasing order",
			args: args{
				a: iterate([]string{"a", "b", "c"}, []int{1, 2, 3}),
				b: iterate([]string{"c", "a", "b"}, []int{1, 2, 3}),
			},
			want: "sequences must be ordered!",
		},
		{
			name: "b decreasing from the middle",
			args: args{
				a: iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
				b: iterate([]string{"a", "b", "d", "c"}, []int{1, 2, 3, 4}),
			},
			want: "sequences must be ordered!",
		},
		{
			name: "b has less value at the end",
			args: args{
				a: iterate([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4}),
				b: iterate([]string{"a", "b", "c", "a"}, []int{1, 2, 3, 4}),
			},
			want: "sequences must be ordered!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.PanicsWithValue(t, tt.want, func() {
				for range Merge(tt.args.a, tt.args.b) {
				}
			})
		})
	}
}
