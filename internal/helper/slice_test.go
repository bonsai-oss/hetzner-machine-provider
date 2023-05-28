package helper_test

import (
	"reflect"
	"testing"

	"github.com/bonsai-oss/hetzner-machine-provider/internal/helper"
)

func TestFilter(t *testing.T) {
	for _, testCase := range []struct {
		name string
		in   []int
		fn   func(int) bool
		want []int
	}{
		{
			name: "filter out even numbers",
			in:   []int{1, 2, 3, 4, 5},
			fn:   func(i int) bool { return i%2 != 0 },
			want: []int{1, 3, 5},
		},
		{
			name: "filter out numbers larger than 5",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			fn:   func(i int) bool { return i <= 5 },
			want: []int{1, 2, 3, 4, 5},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			got := helper.Filter(testCase.in, testCase.fn)
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("Filter() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	for _, testCase := range []struct {
		name string
		in   []int
		fn   func(int) int
		want []int
	}{
		{
			name: "double each number",
			in:   []int{1, 2, 3, 4, 5},
			fn:   func(i int) int { return i * 2 },
			want: []int{2, 4, 6, 8, 10},
		},
		{
			name: "square each number",
			in:   []int{1, 2, 3, 4, 5},
			fn:   func(i int) int { return i * i },
			want: []int{1, 4, 9, 16, 25},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			got := helper.Map(testCase.in, testCase.fn)
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("Map() = %v, want %v", got, testCase.want)
			}
		})
	}
}
