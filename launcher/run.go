package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type LaunchOptions struct {
	Username       string
	GameDir        string
	RamMB          int
	VersionID      string
	StatusCallback func(string)
}

// Launch prepares and executes the Minecraft command
func Launch(opts LaunchOptions) (*exec.Cmd, error) {
	report := func(msg string) {
		if opts.StatusCallback != nil {
			opts.StatusCallback(msg)
		}
	}

	// 1. Get Java
	report("Checking Java...")
	javaPath, err := EnsureJava(opts.GameDir)
	if err != nil {
		fmt.Printf("Warning: Could not auto-download Java, trying system java: %v\n", err)
		javaPath = "java"
	}

	// 2. Load Manifest & Package
	report("Fetching Version Manifest...")
	manifest, err := GetVersionManifest()
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest: %w", err)
	}

	versionURL, err := manifest.FindVersionURL(opts.VersionID)
	if err != nil {
		return nil, err
	}

	report("Fetching Package Info...")
	pkg, err := GetPackage(versionURL)
	if err != nil {
		return nil, err
	}

	// 3. Download Everything (Blocking for now, should be progress-reported)
	// TODO: report progress
	report("Downloading Assets...")
	err = DownloadAssets(pkg.AssetIndex, opts.GameDir)
	if err != nil {
		return nil, fmt.Errorf("assets error: %w", err)
	}

	report("Downloading Libraries...")
	cp, err := DownloadLibraries(pkg.Libraries, opts.GameDir)
	if err != nil {
		return nil, fmt.Errorf("libs error: %w", err)
	}

	// 4. Construct Arguments
	// Add client jar to classpath
	report("Downloading Client Jar...")
	clientJarPath := filepath.Join(opts.GameDir, "versions", pkg.ID, pkg.ID+".jar")
	if err := downloadFile(pkg.Downloads.Client.URL, clientJarPath); err != nil {
		return nil, fmt.Errorf("client jar download failed: %w", err)
	}
	realCp := cp + string(os.PathListSeparator) + clientJarPath

	nativesDir := filepath.Join(opts.GameDir, "natives")

	// Apply M1 Patches if needed
	report("Checking for Native Patches...")
	if err := PatchNatives(nativesDir); err != nil {
		fmt.Printf("Warning: Failed to apply M1 patches: %v\n", err)
	}

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
	report("Launching...")
	fmt.Printf("Executing: %s %v\n", javaPath, args)
	cmd := exec.Command(javaPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = opts.GameDir

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}
