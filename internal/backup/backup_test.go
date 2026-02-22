package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	// Create a temp source file
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "source.txt")
	content := []byte("hello, backup!")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("failed to write source: %v", err)
	}

	dstPath := filepath.Join(tmpDir, "dest.txt")
	n, err := CopyFile(srcPath, dstPath)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	if n != int64(len(content)) {
		t.Errorf("expected %d bytes, got %d", len(content), n)
	}

	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read dest: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", got, content)
	}
}

func TestBackupDir(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create source structure:
	// srcDir/
	//   file1.txt
	//   subdir/
	//     file2.txt
	if err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("file1"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(srcDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("file2"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := BackupDir(srcDir, dstDir)
	if err != nil {
		t.Fatalf("BackupDir failed: %v", err)
	}

	if result.FilesCopied != 2 {
		t.Errorf("expected 2 files copied, got %d", result.FilesCopied)
	}

	if len(result.Errors) != 0 {
		t.Errorf("unexpected errors: %v", result.Errors)
	}

	// Verify files exist in destination
	if _, err := os.Stat(filepath.Join(dstDir, "file1.txt")); os.IsNotExist(err) {
		t.Error("file1.txt not found in backup")
	}
	if _, err := os.Stat(filepath.Join(dstDir, "subdir", "file2.txt")); os.IsNotExist(err) {
		t.Error("subdir/file2.txt not found in backup")
	}
}

func TestBackupDir_InvalidSource(t *testing.T) {
	dstDir := t.TempDir()
	_, err := BackupDir("/nonexistent/path/xyz", dstDir)
	if err == nil {
		t.Error("expected error for nonexistent source, got nil")
	}
}

func TestBackupDir_SourceIsFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "file.txt")
	os.WriteFile(srcFile, []byte("data"), 0644)

	dstDir := t.TempDir()
	_, err := BackupDir(srcFile, dstDir)
	if err == nil {
		t.Error("expected error when source is a file, got nil")
	}
}
