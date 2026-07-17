package sumrange

import "testing"

func TestSumRange(t *testing.T) {
	tests := []struct {
		name       string
		values     []int
		start, end int
		want       int
		wantErr    bool
	}{
		{"whole", []int{4, 7, 2}, 0, 3, 13, false},
		{"prefix", []int{4, 7, 2}, 0, 2, 11, false},
		{"empty middle", []int{4, 7}, 1, 1, 0, false},
		{"empty input", nil, 0, 0, 0, false},
		{"negative start", []int{4}, -1, 1, 0, true},
		{"end too large", []int{4}, 0, 2, 0, true},
		{"reversed", []int{4}, 1, 0, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SumRange(tc.values, tc.start, tc.end)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error=%v, wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr && got != tc.want {
				t.Fatalf("sum=%d, want=%d", got, tc.want)
			}
		})
	}
}
