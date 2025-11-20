package pathsize

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSize_File(t *testing.T) {
	size1, err := GetSize("testdata/file1.txt", false)
	assert.NoError(t, err)
	size2, err := GetSize("testdata/file2.txt", false)
	assert.NoError(t, err)
	assert.Equal(t, size1, size2)
}

func TestGetSize_Dir(t *testing.T) {
	size1, err := GetSize("testdata/dir1/a.txt", false)
	assert.NoError(t, err)
	size2, err := GetSize("testdata/dir1/b.txt", false)
	assert.NoError(t, err)
	size3, err := GetSize("testdata/dir1", false)
	assert.NoError(t, err)
	assert.Equal(t, size1+size2, size3)
}

func TestGetSize_NonExistentDir(t *testing.T) {
	_, err := GetSize("testdata/nonexistent.txt", false)
	assert.Error(t, err)
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name  string
		size  int64
		human bool
		want  string
	}{
		{
			name:  "non human small",
			size:  123,
			human: false,
			want:  "123B",
		},
		{
			name:  "human small less than 1KB",
			size:  512,
			human: true,
			want:  "512B",
		},
		{
			name:  "human exactly 1KB",
			size:  1024,
			human: true,
			want:  "1.0KB",
		},
		{
			name:  "human exactly 1MB",
			size:  1024 * 1024,
			human: true,
			want:  "1.0MB",
		},
		{
			name:  "human example from task",
			size:  1234567,
			human: true,
			want:  "1.2MB",
		},
		{
			name:  "zero size non human",
			size:  0,
			human: false,
			want:  "0B",
		},
		{
			name:  "zero size human",
			size:  0,
			human: true,
			want:  "0B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSize(tt.size, tt.human)
			if got != tt.want {
				t.Fatalf("FormatSize(%d, %v) = %q, want %q",
					tt.size, tt.human, got, tt.want)
			}
		})
	}
}
