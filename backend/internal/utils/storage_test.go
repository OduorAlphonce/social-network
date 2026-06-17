package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveImageAllowsRequiredImageTypes(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		ext  string
	}{
		{
			name: "jpeg",
			data: []byte{0xff, 0xd8, 0xff, 0xdb, 0x00, 0x43, 0x00, 0x08, 0x06, 0x06, 0x07, 0xff, 0xd9},
			ext:  ".jpg",
		},
		{
			name: "png",
			data: []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52},
			ext:  ".png",
		},
		{
			name: "gif",
			data: []byte("GIF89a\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\xff\xff\xff,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02D\x01\x00;"),
			ext:  ".gif",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir() + string(os.PathSeparator)
			path, err := SaveImage(bytes.NewReader(tt.data), dir)
			if err != nil {
				t.Fatalf("SaveImage returned error: %v", err)
			}
			if !strings.HasSuffix(path, tt.ext) {
				t.Fatalf("path = %q, want suffix %q", path, tt.ext)
			}
			if _, err := os.Stat(filepath.Clean(path)); err != nil {
				t.Fatalf("saved file not found: %v", err)
			}
		})
	}
}

func TestSaveImageRejectsUnsupportedType(t *testing.T) {
	_, err := SaveImage(bytes.NewReader([]byte("plain text is not an image")), t.TempDir()+string(os.PathSeparator))
	if err == nil {
		t.Fatal("expected unsupported file type to be rejected")
	}
}
