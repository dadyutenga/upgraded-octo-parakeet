package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dadyutenga/upgraded-octo-parakeet/internal/backup"
)

func main() {
	src := flag.String("src", "", "Source directory to back up")
	dst := flag.String("dst", "", "Destination base directory for backups")
	flag.Parse()

	if *src == "" || *dst == "" {
		fmt.Fprintln(os.Stderr, "Usage: backup -src <source_dir> -dst <dest_dir>")
		fmt.Fprintln(os.Stderr, "\nCreates a timestamped backup of the source directory.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Backing up %s → %s ...\n", *src, *dst)

	dstDir, result, err := backup.BackupDirWithTimestamp(*src, *dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Backup failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Backup complete!\n")
	fmt.Printf("   Destination: %s\n", dstDir)
	fmt.Printf("   Files copied: %d\n", result.FilesCopied)
	fmt.Printf("   Dirs created: %d\n", result.DirsCreated)
	fmt.Printf("   Bytes copied: %d\n", result.BytesCopied)

	if len(result.Errors) > 0 {
		fmt.Printf("\n⚠️  Errors (%d):\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("   - %s\n", e)
		}
	}
}
