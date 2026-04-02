package cgo

import "testing"

func TestParseByteSizeSupportsExtendedUnits(t *testing.T) {
	tests := []struct {
		input string
		want  uint64
	}{
		{input: "1K", want: 1024},
		{input: "1KB", want: 1024},
		{input: "1KiB", want: 1024},
		{input: "1M", want: 1024 * 1024},
		{input: "1MB", want: 1024 * 1024},
		{input: "1MiB", want: 1024 * 1024},
		{input: "1G", want: 1024 * 1024 * 1024},
		{input: "1GB", want: 1024 * 1024 * 1024},
		{input: "1GiB", want: 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		got, err := parseByteSize(tt.input)
		if err != nil {
			t.Fatalf("parseByteSize(%q) error = %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("parseByteSize(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestBuildArtifactsSkipsUnsetValuesAndGroupsBlobSettings(t *testing.T) {
	cfg := &Config{
		WriteBufferSize:           "",
		WriteBufferNumber:         0,
		WriteBufferNumberToMerge:  0,
		MaxBackgroundJobs:         0,
		CompressionType:           "none",
		BottommostCompressionType: "none",
		LRUSize:                   "64MB",
		FilterBitsPerKey:          10,
		BlockSize:                 "4KB",
		BlobFileSize:              "128MB",
		MinBlobSize:               "4KB",
		EnableBlobFiles:           false,
		EnableBlobGC:              true,
	}

	artifacts, err := buildArtifacts(cfg, true, false)
	if err != nil {
		t.Fatalf("buildArtifacts() error = %v", err)
	}
	if artifacts.options == nil || artifacts.blockOptions == nil || artifacts.cache == nil {
		t.Fatal("expected artifacts to be created even when many config fields are unset")
	}
	if artifacts.rateLimiter != nil {
		t.Fatal("expected empty RateBytesPerSec to skip rate limiter creation")
	}
	artifacts.Close()
}

func TestBuildArtifactsUsesNoCachePathWhenLRUSizeIsZero(t *testing.T) {
	cfg := &Config{
		CompressionType:           "none",
		BottommostCompressionType: "none",
		LRUSize:                   "0",
		RateBytesPerSec:           "",
	}

	artifacts, err := buildArtifacts(cfg, true, false)
	if err != nil {
		t.Fatalf("buildArtifacts() error = %v", err)
	}
	defer artifacts.Close()
	if artifacts.blockOptions == nil {
		t.Fatal("expected block options to exist")
	}
	if artifacts.cache != nil {
		t.Fatal("expected no-cache path to skip cache creation")
	}
}
