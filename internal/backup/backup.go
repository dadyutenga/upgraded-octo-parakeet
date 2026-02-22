package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Result holds statistics from a backup operation.
type Result struct {
	FilesCopied int
	DirsCreated int
	BytesCopied int64
	Errors      []string
}

// CopyFile copies a single file from src to dst, preserving permissions.
func CopyFile(src, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("open source: %w", err)
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return 0, fmt.Errorf("stat source: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return 0, fmt.Errorf("create destination: %w", err)
	}
	defer dstFile.Close()

	n, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return n, fmt.Errorf("copy data: %w", err)
	}

	return n, nil
}

// BackupDir copies all files from srcDir to dstDir recursively.
func BackupDir(srcDir, dstDir string) (*Result, error) {
	result := &Result{}

	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return result, fmt.Errorf("resolve source path: %w", err)
	}

	info, err := os.Stat(srcDir)
	if err != nil {
		return result, fmt.Errorf("stat source directory: %w", err)
	}
	if !info.IsDir() {
		return result, fmt.Errorf("source is not a directory: %s", srcDir)
	}

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("walk error: %s: %v", path, err))
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("rel path error: %s: %v", path, err))
			return nil
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			if err := os.MkdirAll(dstPath, info.Mode()); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("mkdir error: %s: %v", dstPath, err))
				return nil
			}
			result.DirsCreated++
			return nil
		}

		n, err := CopyFile(path, dstPath)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("copy error: %s: %v", path, err))
			return nil
		}
		result.FilesCopied++
		result.BytesCopied += n

		return nil
	})

	return result, err
}

// BackupDirWithTimestamp creates a timestamped backup of srcDir inside dstBase.
func BackupDirWithTimestamp(srcDir, dstBase string) (string, *Result, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	baseName := filepath.Base(srcDir)
	dstDir := filepath.Join(dstBase, fmt.Sprintf("%s_%s", baseName, timestamp))

	result, err := BackupDir(srcDir, dstDir)
	return dstDir, result, err
}
