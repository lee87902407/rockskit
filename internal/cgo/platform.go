package cgo

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func nativeArtifactDir(goos, goarch string) (string, error) {
	switch {
	case goos == "darwin" && goarch == "arm64":
		return "darwin_arm64", nil
	case goos == "linux" && goarch == "amd64":
		return "linux_amd64", nil
	case goos == "linux" && goarch == "arm64":
		return "linux_arm64", nil
	default:
		return "", fmt.Errorf("unsupported platform %s/%s", goos, goarch)
	}
}

func requiredNativeArtifacts(goos, goarch string) ([]string, error) {
	dir, err := nativeArtifactDir(goos, goarch)
	if err != nil {
		return nil, err
	}
	baseDir := repoRoot()
	return []string{
		filepath.Join(baseDir, "native", dir, "include", "rocksdb", "c.h"),
		filepath.Join(baseDir, "native", dir, "librocksdb.a"),
	}, nil
}

func nativeArtifactExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
