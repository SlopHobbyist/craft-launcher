package integrity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManifestUpdate_ContentChangeSameVersion(t *testing.T) {
	// Setup temporary game directory
	tmpDir, err := os.MkdirTemp("", "launcher_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// --- 1. Setup Initial State (Local Manifest V1 with File A) ---
	fileA := FileInfo{Path: "mods/A.jar", Size: 4, Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", Override: true} // "test"
	localManifest := Manifest{
		Version: 1,
		Files:   []FileInfo{fileA},
	}

	// Write local manifest
	localManifestPath := filepath.Join(tmpDir, LocalManifest)
	if err := saveLocalManifest(localManifestPath, &localManifest); err != nil {
		t.Fatal(err)
	}

	// Create physical File A
	fileAPath := filepath.Join(tmpDir, fileA.Path)
	if err := os.MkdirAll(filepath.Dir(fileAPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fileAPath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// --- 2. Setup Server State (Manifest V1 with File B, File A removed) ---
	// Same version number, but different content!
	fileB := FileInfo{Path: "mods/B.jar", Size: 4, Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", Override: true} // hash not checked for mock download really, but consistent
	serverManifest := Manifest{
		Version: 1, // SAME VERSION
		Files:   []FileInfo{fileB},
	}

	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/manifest.json") {
			json.NewEncoder(w).Encode(serverManifest)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/files/mods/B.jar") {
			w.Write([]byte("test")) // Content doesn't strictly matter for download test
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	// --- 3. Run CheckAndUpdate ---
	statusLog := []string{}
	callback := func(msg string) {
		statusLog = append(statusLog, msg)
	}

	if err := CheckAndUpdate(tmpDir, server.URL, callback); err != nil {
		t.Fatalf("CheckAndUpdate failed: %v", err)
	}

	// --- 4. Verify Results ---

	// File A should be DELETED
	if _, err := os.Stat(fileAPath); !os.IsNotExist(err) {
		t.Errorf("File A should have been deleted, but exists")
	}

	// File B should be CREATED
	fileBPath := filepath.Join(tmpDir, fileB.Path)
	if _, err := os.Stat(fileBPath); os.IsNotExist(err) {
		t.Errorf("File B should have been downloaded, but is missing")
	}

	// Local Manifest should be UPDATED (to match server, essentially same version but could track new file list if we inspected it)
	updatedLocalManifest, err := loadLocalManifest(localManifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(updatedLocalManifest.Files) != 1 || updatedLocalManifest.Files[0].Path != "mods/B.jar" {
		t.Errorf("Local manifest was not updated correctly. Got: %+v", updatedLocalManifest)
	}
}
