package model

import "testing"

func TestUpload_IsImage(t *testing.T) {
	cases := []struct {
		mime string
		want bool
	}{
		{"image/jpeg", true},
		{"image/jpg", true},
		{"image/png", true},
		{"image/gif", true},
		{"image/webp", true},
		{"image/bmp", true},
		{"image/svg+xml", false},
		{"text/plain", false},
	}

	for _, tt := range cases {
		u := &Upload{MimeType: tt.mime}
		if got := u.IsImage(); got != tt.want {
			t.Fatalf("IsImage(%q) = %v, want %v", tt.mime, got, tt.want)
		}
	}
}

func TestUpload_IsVideo(t *testing.T) {
	cases := []struct {
		mime string
		want bool
	}{
		{"video/mp4", true},
		{"video/mpeg", true},
		{"video/quicktime", true},
		{"video/x-msvideo", true},
		{"video/webm", true},
		{"video/ogg", false},
		{"image/png", false},
	}

	for _, tt := range cases {
		u := &Upload{MimeType: tt.mime}
		if got := u.IsVideo(); got != tt.want {
			t.Fatalf("IsVideo(%q) = %v, want %v", tt.mime, got, tt.want)
		}
	}
}

func TestUpload_IsAudio(t *testing.T) {
	cases := []struct {
		mime string
		want bool
	}{
		{"audio/mpeg", true},
		{"audio/wav", true},
		{"audio/ogg", true},
		{"audio/webm", true},
		{"audio/aac", true},
		{"audio/flac", false},
		{"video/mp4", false},
	}

	for _, tt := range cases {
		u := &Upload{MimeType: tt.mime}
		if got := u.IsAudio(); got != tt.want {
			t.Fatalf("IsAudio(%q) = %v, want %v", tt.mime, got, tt.want)
		}
	}
}

func TestUpload_GetSizeInMB(t *testing.T) {
	u := &Upload{FileSize: 1048576} // 1MB
	got := u.GetSizeInMB()
	if got < 0.99 || got > 1.01 {
		t.Fatalf("GetSizeInMB() = %f, want around 1.0", got)
	}
}
