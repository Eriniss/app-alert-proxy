package logto

import "testing"

func TestAlertEnabled(t *testing.T) {
	cases := []struct {
		name string
		data map[string]any
		want bool
	}{
		{"bool true", map[string]any{"isAlert": true}, true},
		{"string true", map[string]any{"isAlert": "true"}, true},
		{"bool false", map[string]any{"isAlert": false}, false},
		{"키 없음", map[string]any{}, false},
		{"nil 맵", nil, false},
		{"엉뚱한 타입", map[string]any{"isAlert": 1}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := User{CustomData: c.data}.AlertEnabled()
			if got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}
