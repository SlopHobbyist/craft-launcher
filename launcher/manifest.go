package launcher

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const ManifestURL = "https://piston-meta.mojang.com/mc/game/version_manifest.json"

// GetVersionManifest fetches the list of all Minecraft versions
func GetVersionManifest() (*VersionManifest, error) {
	resp, err := http.Get(ManifestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var manifest VersionManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// GetPackage fetches the specific version.json (e.g. for 1.8.9)
func GetPackage(url string) (*Package, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pkg Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}

// FindVersionURL returns the URL for a specific version ID (e.g. "1.8.9")
func (vm *VersionManifest) FindVersionURL(id string) (string, error) {
	for _, v := range vm.Versions {
		if v.ID == id {
			return v.URL, nil
		}
	}
	return "", fmt.Errorf("version %s not found", id)
}
