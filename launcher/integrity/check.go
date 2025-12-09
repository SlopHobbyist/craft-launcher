package integrity

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// VerifyAndRestore checks all protected files and restores them if compromised
func VerifyAndRestore(gameDir string) ([]string, error) {
	files := GetProtectedFiles()
	var restored []string

	for _, file := range files {
		// Target path on disk
		targetPath := filepath.Join(gameDir, file.Path)

		needsRestore := false

		// check existence
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			needsRestore = true
		} else {
			// check hash
			currentHash, err := hashFile(targetPath)
			if err != nil {
				// if we can't read it, we should probably restore it
				fmt.Printf("Error hashing %s: %v. Restoring.\n", targetPath, err)
				needsRestore = true
			} else if currentHash != file.Hash {
				fmt.Printf("Integrity mismatch for %s. Expected %s, got %s\n", file.Path, file.Hash, currentHash)
				needsRestore = true
			}
		}

		if needsRestore {
			fmt.Printf("Restoring protected file: %s\n", file.Path)
			if err := restoreFile(targetPath, file.Content); err != nil {
				return restored, fmt.Errorf("failed to restore %s: %w", file.Path, err)
			}
			restored = append(restored, file.Path)
		}
	}

	return restored, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func restoreFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	// Use 0644 for files
	return os.WriteFile(path, content, 0644)
}
