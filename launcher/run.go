package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type LaunchOptions struct {
	Username  string
	GameDir   string
	RamMB     int
	VersionID string
}

// Launch prepares and executes the Minecraft command
func Launch(opts LaunchOptions) error {
	// 1. Get Java
	javaPath, err := EnsureJava(opts.GameDir)
	if err != nil {
		fmt.Printf("Warning: Could not auto-download Java, trying system java: %v\n", err)
		javaPath = "java"
	}

	// 2. Load Manifest & Package
	manifest, err := GetVersionManifest()
	if err != nil {
		return fmt.Errorf("failed to get manifest: %w", err)
	}

	versionURL, err := manifest.FindVersionURL(opts.VersionID)
	if err != nil {
		return err
	}

	pkg, err := GetPackage(versionURL)
	if err != nil {
		return err
	}

	// 3. Download Everything (Blocking for now, should be progress-reported)
	// TODO: report progress
	err = DownloadAssets(pkg.AssetIndex, opts.GameDir)
	if err != nil {
		return fmt.Errorf("assets error: %w", err)
	}

	cp, err := DownloadLibraries(pkg.Libraries, opts.GameDir)
	if err != nil {
		return fmt.Errorf("libs error: %w", err)
	}

	// 4. Construct Arguments
	// Add client jar to classpath
	clientJarPath := filepath.Join(opts.GameDir, "versions", pkg.ID, pkg.ID+".jar")
	if err := downloadFile(pkg.Downloads.Client.URL, clientJarPath); err != nil {
		return fmt.Errorf("client jar download failed: %w", err)
	}
	realCp := cp + string(os.PathListSeparator) + clientJarPath

	nativesDir := filepath.Join(opts.GameDir, "natives")

	// Arguments construction
	// 1.8.9 uses "minecraftArguments" string, newer versions use "arguments" object.
	// We focus on 1.8.9 here.

	args := []string{
		fmt.Sprintf("-Xmx%dM", opts.RamMB),
		fmt.Sprintf("-Djava.library.path=%s", nativesDir),
		"-cp", realCp,
		pkg.MainClass,
	}

	// Parse minecraftArguments template
	// e.g. "--username ${auth_player_name} --version ${version_name} --gameDir ${game_directory} --assetsDir ${assets_root} --assetIndex ${assets_index_name} --uuid ${auth_uuid} --accessToken ${auth_access_token} --userProperties ${user_properties} --userType ${user_type}"
	mcArgs := pkg.MinecraftArgs
	replacements := map[string]string{
		"${auth_player_name}":  opts.Username,
		"${version_name}":      pkg.ID,
		"${game_directory}":    opts.GameDir,
		"${assets_root}":       filepath.Join(opts.GameDir, "assets"),
		"${assets_index_name}": pkg.AssetIndex.ID,
		"${auth_uuid}":         "00000000-0000-0000-0000-000000000000", // Offline UUID
		"${auth_access_token}": "null",                                 // Offline Token
		"${user_properties}":   "{}",
		"${user_type}":         "legacy",
	}

	for k, v := range replacements {
		mcArgs = strings.ReplaceAll(mcArgs, k, v)
	}

	args = append(args, strings.Split(mcArgs, " ")...)

	// 5. Execute
	fmt.Printf("Executing: %s %v\n", javaPath, args)
	cmd := exec.Command(javaPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = opts.GameDir

	return cmd.Start() // Non-blocking start
}
