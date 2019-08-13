package geohash

import (
	"testing"
)

type TestCase struct {
}

type ValidPointTestcase struct {
	p    Point
	want error
}

type EncodingTestcase struct {
	p         Point
	precision int
	want      string
	wantErr   error
}

func TestValidPoint(t *testing.T) {
	for _, c := range vpTestcases {
		got := ValidPoint(c.p)

		if got != c.want {
			t.Errorf("Incorrect point validation for '%v'; expected '%v', got '%v'", c.p, c.want, got)
		}
	}
}

func TestEncodeUsingPrecision(t *testing.T) {
	for _, c := range encTestcases {
		got, gotErr := EncodeUsingPrecision(c.p, c.precision)
		if got != c.want || gotErr != c.wantErr {
			t.Errorf("Incorrect encoding for '%v' and '%v'; expected '%v' and '%v', got '%v' and '%v'", c.p, c.precision, c.want, c.wantErr, got, gotErr)
		}

	}
}
