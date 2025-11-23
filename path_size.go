package code

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

// NewApp creates and configures the hexlet-path-size CLI application.
func NewApp() *cli.Command {
	return &cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
		UsageText: "hexlet-path-size [global options]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "recursive size of directories",
			},
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit)",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() < 1 {
				return fmt.Errorf("no path")
			}

			path := cmd.Args().Get(0)
			recursive := cmd.Bool("recursive")
			human := cmd.Bool("human")
			all := cmd.Bool("all")

			result, err := GetPathSize(path, recursive, human, all)
			if err != nil {
				return err
			}

			fmt.Printf("%s\t%s\n", result, path)
			return nil
		},
	}
}

// calculateSize walks the given path and returns its total size in bytes.
func calculateSize(path string, recursive, all bool) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var size int64

	for _, entry := range entries {
		name := entry.Name()

		if !all && strings.HasPrefix(name, ".") {
			continue
		}

		childPath := filepath.Join(path, name)

		if entry.IsDir() {
			if !recursive {
				continue
			}

			childSize, err := calculateSize(childPath, recursive, all)
			if err != nil {
				return 0, err
			}

			size += childSize
			continue
		}

		childInfo, err := os.Lstat(childPath)
		if err != nil {
			return 0, err
		}

		if childInfo.Mode().IsRegular() {
			size += childInfo.Size()
		}
	}

	return size, nil
}

// GetPathSize calculates size of the given path.
// It supports human-readable output and flags for including hidden files
// and recursive directory traversal.
func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := calculateSize(path, recursive, all)
	if err != nil {
		return "", err
	}

	return FormatSize(size, human), nil
}

// FormatSize converts raw size in bytes to a formatted string.
// When human is true, it uses units B, KB, MB, GB, etc.
func FormatSize(size int64, human bool) string {
	if !human {
		return fmt.Sprintf("%dB", size)
	}

	const unit = 1024.0
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}

	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}

	value := float64(size)
	i := 0

	for value >= unit && i < len(units)-1 {
		value /= unit
		i++
	}

	return fmt.Sprintf("%.1f%s", value, units[i])
}
