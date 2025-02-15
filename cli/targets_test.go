package cli //nolint: testpackage // Need access to internals

import "testing"

func TestIsValidTarget(t *testing.T) {
	tests := []struct {
		name   string
		target string
		want   bool
	}{
		{
			name:   "vscode exists",
			target: "visualstudiocode",
			want:   true,
		},
		{
			name:   "macos exists",
			target: "macos",
			want:   true,
		},
		{
			name:   "python exists",
			target: "python",
			want:   true,
		},
		{
			name:   "nonsense doesnt exist",
			target: "nonsense",
			want:   false,
		},
		{
			name:   "blah doesnt exist",
			target: "blah",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidTarget(tt.target); got != tt.want {
				t.Errorf("got %v, wanted %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsValidTarget(b *testing.B) {
	for range b.N {
		IsValidTarget("go")
	}
}
