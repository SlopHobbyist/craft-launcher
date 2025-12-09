package launcher

import "time"

// VersionManifest represents the main version list from Mojang
type VersionManifest struct {
	Latest   Latest    `json:"latest"`
	Versions []Version `json:"versions"`
}

type Latest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type Version struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	URL         string    `json:"url"`
	Time        time.Time `json:"time"`
	ReleaseTime time.Time `json:"releaseTime"`
}

// Package represents the specific version.json (e.g., 1.8.9.json)
type Package struct {
	AssetIndex    AssetIndex `json:"assetIndex"`
	Assets        string     `json:"assets"`
	Downloads     Downloads  `json:"downloads"`
	ID            string     `json:"id"`
	Libraries     []Library  `json:"libraries"`
	MainClass     string     `json:"mainClass"`
	MinecraftArgs string     `json:"minecraftArguments"` // Legacy (1.8.9)
	Type          string     `json:"type"`
}

type AssetIndex struct {
	ID        string `json:"id"`
	Sha1      string `json:"sha1"`
	Size      int    `json:"size"`
	TotalSize int    `json:"totalSize"`
	URL       string `json:"url"`
}

type Downloads struct {
	Client        DownloadInfo `json:"client"`
	Server        DownloadInfo `json:"server"`
	WindowsServer DownloadInfo `json:"windows_server,omitempty"`
}

type DownloadInfo struct {
	Sha1 string `json:"sha1"`
	Size int    `json:"size"`
	URL  string `json:"url"`
}

type Library struct {
	Downloads LibDownloads      `json:"downloads"`
	Name      string            `json:"name"`
	Natives   map[string]string `json:"natives,omitempty"`
	Rules     []Rule            `json:"rules,omitempty"`
}

type LibDownloads struct {
	Artifact    *Artifact            `json:"artifact,omitempty"`
	Classifiers map[string]*Artifact `json:"classifiers,omitempty"` // Natives often here
}

type PkgMainClass struct {
	Client string `json:"client"`
	Server string `json:"server"`
}

// Fabric Structs

type FabricLoaderResponse struct {
	Loader       FabricLoaderMeta `json:"loader"`
	Intermediary FabricLoaderMeta `json:"intermediary"`
	LaunchMeta   FabricLaunchMeta `json:"launcherMeta"`
}

type FabricLoaderMeta struct {
	Version string `json:"version"`
	Maven   string `json:"maven"`
}

type FabricLaunchMeta struct {
	Libraries FabricLibraries `json:"libraries"`
	MainClass FabricMainClass `json:"mainClass"`
}

type FabricLibraries struct {
	Client []FabricLibrary `json:"client"`
	Common []FabricLibrary `json:"common"`
}

type FabricLibrary struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	MD5  string `json:"md5"`
	SHA1 string `json:"sha1"`
	Size int64  `json:"size"`
}

type FabricMainClass struct {
	Client string `json:"client"`
	Server string `json:"server"`
}

type Artifact struct {
	Path string `json:"path"`
	Sha1 string `json:"sha1"`
	Size int    `json:"size"`
	URL  string `json:"url"`
}

type Rule struct {
	Action string `json:"action"`
	OS     OS     `json:"os,omitempty"`
}

type OS struct {
	Name string `json:"name"`
}
