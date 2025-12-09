package launcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	LegacyFabricMetaURL = "https://meta.legacyfabric.net/v2/versions/loader/1.8.9"
)

// GetFabricMeta fetches the loader metadata for 1.8.9
func GetFabricMeta() (*FabricLoaderResponse, error) {
	resp, err := http.Get(LegacyFabricMetaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch fabric meta: %s", resp.Status)
	}

	var data []FabricLoaderResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no fabric versions found")
	}

	// Return the first one (usually the stable/latest one)
	return &data[0], nil
}

// DownloadFabricLibraries downloads the required fabric libraries
func DownloadFabricLibraries(meta *FabricLoaderResponse, gameDir string) ([]string, error) {
	libsDir := filepath.Join(gameDir, "libraries")
	var cp []string

	// Helper to process a list of libraries
	processLibs := func(libs []FabricLibrary) error {
		for _, lib := range libs {
			// Convert maven coordinates to path
			parts := strings.Split(lib.Name, ":")
			if len(parts) != 3 {
				fmt.Printf("Skipping invalid fabric lib: %s\n", lib.Name)
				continue
			}
			group := parts[0]
			artifact := parts[1]
			version := parts[2]

			relPath := fmt.Sprintf("%s/%s/%s/%s-%s.jar",
				strings.ReplaceAll(group, ".", "/"),
				artifact,
				version,
				artifact,
				version,
			)

			destPath := filepath.Join(libsDir, relPath)
			absPath, _ := filepath.Abs(destPath)
			cp = append(cp, absPath)

			if _, err := os.Stat(destPath); err == nil {
				continue // Already exists
			}

			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}

			// Construct download URL
			downloadURL := lib.URL + relPath
			if lib.URL == "" {
				downloadURL = "https://maven.fabricmc.net/" + relPath
			}

			fmt.Printf("Downloading Fabric Lib: %s\n", lib.Name)
			if err := downloadFile(downloadURL, destPath); err != nil {
				return fmt.Errorf("failed to download %s: %w", lib.Name, err)
			}
		}
		return nil
	}

	// 1. Download Common Libraries
	if err := processLibs(meta.LaunchMeta.Libraries.Common); err != nil {
		return nil, err
	}

	// 2. Download Client Libraries
	if err := processLibs(meta.LaunchMeta.Libraries.Client); err != nil {
		return nil, err
	}

	// Helper for maven string
	downloadMaven := func(mavenStr, repoBase string) (string, error) {
		parts := strings.Split(mavenStr, ":")
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid maven: %s", mavenStr)
		}
		group, artifact, version := parts[0], parts[1], parts[2]
		relPath := fmt.Sprintf("%s/%s/%s/%s-%s.jar",
			strings.ReplaceAll(group, ".", "/"),
			artifact,
			version,
			artifact,
			version,
		)
		destPath := filepath.Join(libsDir, relPath)
		absPath, _ := filepath.Abs(destPath)

		if _, err := os.Stat(destPath); err == nil {
			return absPath, nil
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return "", err
		}

		url := repoBase + relPath
		fmt.Printf("Downloading Maven Art: %s\n", mavenStr)
		if err := downloadFile(url, destPath); err != nil {
			return "", err
		}
		return absPath, nil
	}

	// 3. Download Intermediary
	// Intermediary URL usually: https://maven.legacyfabric.net/
	interPath, err := downloadMaven(meta.Intermediary.Maven, "https://maven.legacyfabric.net/")
	if err != nil {
		return nil, fmt.Errorf("failed to download intermediary: %w", err)
	}
	cp = append(cp, interPath)

	// 4. Download Loader
	// Loader URL usually: https://maven.fabricmc.net/
	loaderPath, err := downloadMaven(meta.Loader.Maven, "https://maven.fabricmc.net/")
	if err != nil {
		return nil, fmt.Errorf("failed to download loader: %w", err)
	}
	cp = append(cp, loaderPath)

	return cp, nil
}
