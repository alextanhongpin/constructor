package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

const DomainGithub = "github.com"

func RelativePath(rel string) string {
	path, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("loader: failed to get working directory: %w", err))
	}
	path = filepath.Join(path, rel)
	return path
}

// TrimExtension removes the file extension if exists.
func TrimExtension(path string) string {
	if ext := filepath.Ext(path); ext != "" {
		return path[:len(path)-len(ext)]
	}
	return path
}

// PackagePath returns the package name of the given path.
// If the package is hosted elsewhere other than Github, specify
// the package path.
func PackagePath(domain, path string) string {
	if domain == "" || path == "" {
		return ""
	}
	path = TrimExtension(path)
	path = strings.TrimRight(path, "/")
	return path[strings.Index(path, domain):]
}

// PackageName returns the package name of the given path.
func PackageName(domain, path string) string {
	return filepath.Base(PackagePath(domain, path))
}

// LoadPackage loads the package at the given path.
func LoadPackage(path string) *packages.Package {
	mode := packages.NeedName | packages.NeedTypes | packages.NeedImports
	cfg := &packages.Config{
		Mode: mode,
	}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		panic(fmt.Errorf("loader: failed to load package: %w", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}
	return pkgs[0]
}
