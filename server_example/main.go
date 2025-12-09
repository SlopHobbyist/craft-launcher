package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	ServerURL       = "http://192.168.1.48:8090"
	LocalModpackDir = "./local_modpack"
	VersionFile     = "./local_modpack/.client_version"
)

type Manifest struct {
	Version int        `json:"version"`
	Files   []FileInfo `json:"files"`
}

type FileInfo struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
}

func main() {
	// Ensure modpack directory exists
	os.MkdirAll(LocalModpackDir, 0755)

	// Get current local version
	localVersion := getLocalVersion()
	fmt.Printf("Local version: %d\n", localVersion)

	// Fetch server manifest
	manifest, err := fetchManifest()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server version: %d\n", manifest.Version)

	// Check if update needed
	if manifest.Version > localVersion {
		fmt.Println("Update required! Downloading files...")
		updateModpack(manifest)
		saveLocalVersion(manifest.Version)
		fmt.Println("Update complete!")
	} else if manifest.Version == localVersion {
		// Verify integrity even on same version
		fmt.Println("Verifying file integrity...")
		verified := verifyIntegrity(manifest)
		if !verified {
			fmt.Println("Tampering detected! Re-downloading all files...")
			updateModpack(manifest)
		} else {
			fmt.Println("All files verified. No update needed.")
		}
	}
}

func fetchManifest() (*Manifest, error) {
	resp, err := http.Get(ServerURL + "/manifest.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var manifest Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func updateModpack(manifest *Manifest) {
	// Remove all existing files (clean slate)
	cleanModpackDir()

	// Download all files
	for i, file := range manifest.Files {
		fmt.Printf("[%d/%d] Downloading %s\n", i+1, len(manifest.Files), file.Path)
		if err := downloadFile(file); err != nil {
			fmt.Printf("Error downloading %s: %v\n", file.Path, err)
			continue
		}

		// Verify checksum immediately after download
		if !verifyFileChecksum(file) {
			fmt.Printf("WARNING: Checksum mismatch for %s\n", file.Path)
		}
	}

	// Save checksums for future verification
	saveManifest(manifest)
}

func downloadFile(file FileInfo) error {
	// Create directory structure
	localPath := filepath.Join(LocalModpackDir, file.Path)
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Download file
	url := fmt.Sprintf("%s/files/%s", ServerURL, file.Path)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Save to disk
	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func verifyFileChecksum(file FileInfo) bool {
	localPath := filepath.Join(LocalModpackDir, file.Path)
	f, err := os.Open(localPath)
	if err != nil {
		return false
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false
	}

	checksum := hex.EncodeToString(h.Sum(nil))
	return checksum == file.Checksum
}

func verifyIntegrity(manifest *Manifest) bool {
	allValid := true

	for _, file := range manifest.Files {
		if !verifyFileChecksum(file) {
			fmt.Printf("Tamper detected: %s\n", file.Path)
			allValid = false
		}
	}

	return allValid
}

func cleanModpackDir() {
	// Remove all files except .client_version
	filepath.Walk(LocalModpackDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == LocalModpackDir {
			return nil
		}
		if !info.IsDir() && !strings.HasSuffix(path, ".client_version") && !strings.HasSuffix(path, ".manifest.json") {
			os.Remove(path)
		}
		return nil
	})
}

func getLocalVersion() int {
	data, err := os.ReadFile(VersionFile)
	if err != nil {
		return 0
	}

	var version int
	fmt.Sscanf(string(data), "%d", &version)
	return version
}

func saveLocalVersion(version int) {
	os.WriteFile(VersionFile, []byte(fmt.Sprintf("%d", version)), 0644)
}

func saveManifest(manifest *Manifest) {
	data, _ := json.MarshalIndent(manifest, "", "  ")
	os.WriteFile(filepath.Join(LocalModpackDir, ".manifest.json"), data, 0644)
}
