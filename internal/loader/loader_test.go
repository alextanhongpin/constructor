package loader_test

import (
	"testing"

	"github.com/alextanhongpin/constructor/internal/loader"
)

func TestTrimExtension(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"", ""},
		{"/path/to/file.ext", "/path/to/file"},
		{"/path/to/file", "/path/to/file"},
	}

	for _, tt := range tests {
		if got := loader.TrimExtension(tt.in); got != tt.out {
			t.Errorf("want TrimExtension(%q) = %q, got %q", tt.in, tt.out, got)
		}
	}
}

func TestPackagePath(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"", ""},
		{"/path/to/github.com/organization/yourpackage/main.go", "github.com/organization/yourpackage/main"},
	}
	for _, tt := range tests {
		if got := loader.PackagePath(loader.DomainGithub, tt.in); got != tt.out {
			t.Errorf("want PackagePath(%q) = %q, got %q", tt.in, tt.out, got)
		}
	}
}

func TestPackageName(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"", ""},
		{"/path/to/github.com/organization/yourpackage/main.go", "github.com/organization/yourpackage/main"},
	}
	for _, tt := range tests {
		if got := loader.PackagePath(loader.DomainGithub, tt.in); got != tt.out {
			t.Errorf("want PackageName(%q) = %q, got %q", tt.in, tt.out, got)
		}
	}
}
