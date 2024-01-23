package main

import "testing"

func Test_transformURLFor(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`{{ url_for('static', path='favicon.ico') }}`, `static/favicon.ico`},
		{`{{ url_for('static', path='/css/install.css') }}`, `static/css/install.css`},
		{`{{ url_for('install_license') }}`, `install_license`},
		{`{{ nothing }}`, ``},
	}
	for _, tt := range tests {
		if got := transformURLFor(tt.input); got != tt.want {
			t.Errorf("Input: %s, Expected: %s, Got: %s", tt.input, tt.want, got)
		}
	}
}
