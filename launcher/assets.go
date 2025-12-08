package launcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const AssetBaseURL = "https://resources.download.minecraft.net"

type AssetObject struct {
	Hash string `json:"hash"`
	Size int    `json:"size"`
}

type Assets struct {
	Objects map[string]AssetObject `json:"objects"`
}

// DownloadAssets downloads the asset index and all referenced assets
func DownloadAssets(assetIndexDetails AssetIndex, gameDir string) error {
	// 1. Download Asset Index
	indexesDir := filepath.Join(gameDir, "assets", "indexes")
	if err := os.MkdirAll(indexesDir, 0755); err != nil {
		return err
	}

	indexPath := filepath.Join(indexesDir, assetIndexDetails.ID+".json")
	if err := downloadFile(assetIndexDetails.URL, indexPath); err != nil {
		return fmt.Errorf("failed to download asset index: %w", err)
	}

	// 2. Parse Asset Index
	indexFile, err := os.Open(indexPath)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	var assets Assets
	if err := json.NewDecoder(indexFile).Decode(&assets); err != nil {
		return err
	}

	// 3. Download Objects
	objectsDir := filepath.Join(gameDir, "assets", "objects")
	for _, obj := range assets.Objects {
		// Object path structure: /hash_prefix_2chars/full_hash
		prefix := obj.Hash[:2]
		objPath := filepath.Join(objectsDir, prefix, obj.Hash)
		objURL := fmt.Sprintf("%s/%s/%s", AssetBaseURL, prefix, obj.Hash)

		// Simple check if exists (should verify hash in production, keeping simple for now)
		if _, err := os.Stat(objPath); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(objPath), 0755); err != nil {
				return err
			}
			// TODO: Add concurrency here significantly for speed
			if err := downloadFile(objURL, objPath); err != nil {
				fmt.Printf("Failed to download asset %s: %v\n", obj.Hash, err)
				// Continue strictly? Or fail? For now log and continue.
			}
		}
	}

	return nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
