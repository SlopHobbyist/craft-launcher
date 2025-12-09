package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	BundledDir   = "bundled"
	IntegrityDir = "launcher/integrity"
	EmbedDir     = IntegrityDir + "/embedded_assets"
	GenFile      = IntegrityDir + "/loader.go"
)

// Whitelisted definitions - these can be modified by user
var Whitelist = []string{
	"options.txt",
	"optionsof.txt",
	"optionsshaders.txt",
	"usercache.json",
	"logs",
	"crash-reports",
	"saves",
	"resourcepacks",
	"screenshots",
	"shaderpacks",
	"versions",  // managed by launcher
	"assets",    // managed by launcher
	"libraries", // managed by launcher
	"natives",   // managed by launcher
	"launcher_profiles.json",
}

type FileAsset struct {
	Path   string // Relative path in game dir, e.g. "mods/foo.jar"
	Hash   string // SHA256
	GoName string // Safe variable name for embedding
}

func main() {
	if _, err := os.Stat(BundledDir); os.IsNotExist(err) {
		fmt.Println("No bundled/ directory found. Skipping integrity generation.")
		// Generate empty loader to allow compilation
		generateEmptyLoader()
		return
	}

	fmt.Println("Generating integrity assets from bundled/...")

	// 1. Clean previous assets
	os.RemoveAll(EmbedDir)
	os.MkdirAll(EmbedDir, 0755)
	os.MkdirAll(IntegrityDir, 0755)

	var assets []FileAsset
	counter := 0

	err := filepath.WalkDir(BundledDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Skip .git or filtered dirs
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Rel path from bundled/
		relPath, err := filepath.Rel(BundledDir, path)
		if err != nil {
			return err
		}

		// Normalize path to forward slashes for consistency
		relPath = filepath.ToSlash(relPath)

		// Filter OS junk files
		baseName := filepath.Base(relPath)
		if baseName == ".DS_Store" || baseName == "Thumbs.db" || baseName == ".AppleDouble" || baseName == ".LSOverride" || strings.HasPrefix(baseName, "._") {
			return nil
		}

		// Check filter
		if isWhitelisted(relPath) {
			// fmt.Printf("Skipping whitelisted: %s\n", relPath)
			return nil
		}

		if relPath == "README.txt" {
			return nil
		}

		// Process File
		// Copy to EmbedDir with unique name
		goName := fmt.Sprintf("Asset_%d", counter)
		counter++
		destPath := filepath.Join(EmbedDir, goName)

		hash, err := copyAndHash(path, destPath)
		if err != nil {
			return fmt.Errorf("failed to copy %s: %w", path, err)
		}

		assets = append(assets, FileAsset{
			Path:   relPath,
			Hash:   hash,
			GoName: goName,
		})
		fmt.Printf("Protected: %s\n", relPath)

		return nil
	})

	if err != nil {
		panic(err)
	}

	if err := generateLoader(assets); err != nil {
		panic(err)
	}

	fmt.Printf("Integrity generation complete. %d protected files.\n", len(assets))
}

func isWhitelisted(path string) bool {
	for _, w := range Whitelist {
		// Exact match or subdirectory match
		if path == w || strings.HasPrefix(path, w+"/") {
			return true
		}
	}
	return false
}

func copyAndHash(src, dest string) (string, error) {
	s, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer s.Close()

	d, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer d.Close()

	hasher := sha256.New()
	multiWriter := io.MultiWriter(d, hasher)

	if _, err := io.Copy(multiWriter, s); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func generateLoader(assets []FileAsset) error {
	tmpl := `package integrity

import (
	_ "embed"
)

type ProtectedFile struct {
	Path    string
	Hash    string
	Content []byte
}

{{ range . }}
//go:embed embedded_assets/{{ .GoName }}
var {{ .GoName }} []byte
{{ end }}

func GetProtectedFiles() []ProtectedFile {
	return []ProtectedFile{
{{ range . }}
		{
			Path:    "{{ .Path }}",
			Hash:    "{{ .Hash }}",
			Content: {{ .GoName }},
		},
{{ end }}
	}
}
`
	t, err := template.New("loader").Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(GenFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, assets)
}

func generateEmptyLoader() {
	os.MkdirAll(IntegrityDir, 0755)
	content := `package integrity

type ProtectedFile struct {
	Path    string
	Hash    string
	Content []byte
}

func GetProtectedFiles() []ProtectedFile {
	return []ProtectedFile{}
}
`
	os.WriteFile(GenFile, []byte(content), 0644)
}
