package cgo

import "testing"

func TestNativeArtifactDir(t *testing.T) {
	tests := []struct {
		goos   string
		goarch string
		want   string
	}{
		{goos: "darwin", goarch: "arm64", want: "darwin_arm64"},
		{goos: "linux", goarch: "amd64", want: "linux_amd64"},
		{goos: "linux", goarch: "arm64", want: "linux_arm64"},
	}

	for _, tt := range tests {
		t.Run(tt.goos+"_"+tt.goarch, func(t *testing.T) {
			got, err := nativeArtifactDir(tt.goos, tt.goarch)
			if err != nil {
				t.Fatalf("nativeArtifactDir() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("nativeArtifactDir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnsupportedPlatform(t *testing.T) {
	if _, err := nativeArtifactDir("windows", "amd64"); err == nil {
		t.Fatal("expected unsupported platform error")
	}
}

func TestRequiredNativeArtifactsExist(t *testing.T) {
	tests := []struct {
		goos   string
		goarch string
	}{
		{goos: "darwin", goarch: "arm64"},
		{goos: "linux", goarch: "amd64"},
		{goos: "linux", goarch: "arm64"},
	}

	for _, tt := range tests {
		paths, err := requiredNativeArtifacts(tt.goos, tt.goarch)
		if err != nil {
			t.Fatalf("requiredNativeArtifacts(%s,%s) error = %v", tt.goos, tt.goarch, err)
		}
		if len(paths) != 2 {
			t.Fatalf("expected 2 artifact paths for %s/%s, got %d", tt.goos, tt.goarch, len(paths))
		}
		for _, path := range paths {
			if !nativeArtifactExists(path) {
				t.Fatalf("expected native artifact %q to exist", path)
			}
		}
	}
}
