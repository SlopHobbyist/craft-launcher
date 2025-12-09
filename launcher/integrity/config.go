package integrity

// ServerURL is the endpoint for fetching the modpack manifest.
// It is injected at build time via -ldflags.
// Default is empty or loopback for safety.
var ServerURL = "http://127.0.0.1:8090"
