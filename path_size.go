package pathsize

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
	"path/filepath"
	"strings"
)

func NewApp() *cli.Command {
	return &cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory",
		UsageText: "hexlet-path-size [global options] <path>",
		Flags: []cli.Flag{
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
			all := cmd.Bool("all")
			size, err := GetSize(path, all)

			if err != nil {
				return err
			}

			human := cmd.Bool("human")
			formatted := FormatSize(size, human)
			fmt.Printf("%s\t%s\n", formatted, path)

			return nil
		},
	}
}

func GetSize(path string, all bool) (int64, error) {
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
