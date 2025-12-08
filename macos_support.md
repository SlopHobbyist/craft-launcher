# macOS ARM64 (Apple Silicon) Support - Critical Implementation Notes

## DO NOT MODIFY - Working Configuration

This launcher successfully runs Minecraft 1.8.9 on macOS ARM64 (M1/M2/M3) Macs. The implementation is **simple and proven to work**. Do not "improve" or "modernize" it without testing.

## What Makes It Work

### 1. Native Library Patching ONLY

The M1 fix consists of **only replacing native `.dylib` files** in the natives directory:

- `liblwjgl.dylib` - ARM64-compiled LWJGL native library
- `libopenal.dylib` - ARM64-compiled OpenAL native library

**Source**: [GreeniusGenius/m1-prism-launcher-hack-1.8.9](https://github.com/GreeniusGenius/m1-prism-launcher-hack-1.8.9)

### 2. Leave the Classpath Alone

**DO NOT** filter, modify, or replace LWJGL JAR files in the classpath. The standard LWJGL 2.9.x JAR files from Mojang work fine on ARM64. Only the **native libraries** need replacement.

### 3. NO -XstartOnFirstThread Flag

**CRITICAL**: Do NOT add the `-XstartOnFirstThread` JVM flag on macOS ARM64 with LWJGL 2.

- This flag **causes crashes** with the message: `signal: trace/BPT trap`
- This is the opposite of typical macOS LWJGL advice (which applies to LWJGL 3+)
- LWJGL 2.9.x behaves differently than LWJGL 3.x

## What NOT to Do

### ❌ Do NOT Use "Fat Jars"

Previous broken implementations attempted to:
- Download `lwjglfat.jar` and `lwjgl_util.jar`
- Filter out vanilla LWJGL jars from classpath
- Replace them with ARM64-compiled fat jars

**This approach causes crashes.** The native library replacement alone is sufficient.

### ❌ Do NOT Add -XstartOnFirstThread

```go
// WRONG - This causes crashes on ARM64 with LWJGL 2
if runtime.GOOS == "darwin" {
    args = append(args, "-XstartOnFirstThread")
}
```

This flag is needed for LWJGL 3+ but **breaks LWJGL 2** on Apple Silicon.

### ❌ Do NOT Manipulate the Classpath for M1

The vanilla classpath from Mojang works perfectly:
- `lwjgl-2.9.2.jar`
- `lwjgl-platform-2.9.2-natives-osx.jar`
- `lwjgl_util-2.9.2.jar`

**Leave them in the classpath.**

### ❌ Do NOT Remove "Incompatible" Libraries Preemptively

Previous implementations tried to remove:
- `libjinput-osx.dylib`
- `libtwitchsdk.dylib`

These don't need to be removed. The Twitch SDK will fail gracefully (see log: "Couldn't initialize twitch stream"), which is expected and harmless.

## The Complete Working Implementation

```go
// patch_m1.go - COMPLETE FILE
func PatchNatives(nativesDir string) error {
    // Only run on macOS ARM64
    if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
        return nil
    }

    fmt.Println("Applying Apple Silicon (M1/M2) patches to natives...")

    // 1. Download patched liblwjgl.dylib
    lwjglPath := filepath.Join(nativesDir, "liblwjgl.dylib")
    if err := downloadFile(PatchedLwjglUrl, lwjglPath); err != nil {
        return fmt.Errorf("failed to patch liblwjgl.dylib: %w", err)
    }

    // 2. Download patched libopenal.dylib
    openalPath := filepath.Join(nativesDir, "libopenal.dylib")
    if err := downloadFile(PatchedOpenalUrl, openalPath); err != nil {
        return fmt.Errorf("failed to patch libopenal.dylib: %w", err)
    }

    return nil
}
```

In `run.go`:
```go
// Call after extracting natives, before launching
if err := PatchNatives(nativesDir); err != nil {
    fmt.Printf("Warning: Failed to apply M1 patches: %v\n", err)
}

// Standard JVM args - NO special macOS flags for LWJGL 2
args := []string{
    fmt.Sprintf("-Xmx%dM", opts.RamMB),
    fmt.Sprintf("-Djava.library.path=%s", nativesDir),
    "-cp", realCp,  // Use the vanilla classpath unchanged
    pkg.MainClass,
}
```

## Why This Works

1. **Java JARs are platform-independent** - LWJGL's Java bytecode runs fine on ARM64 JVMs
2. **Native libraries are platform-specific** - The `.dylib` files must be compiled for ARM64
3. **The patched dylibs come from official LWJGL sources** - Compiled from [lwjgl.org/source](https://www.lwjgl.org/source) by Tanmay
4. **LWJGL 2 threading model differs from LWJGL 3** - It doesn't require `-XstartOnFirstThread`

## Testing Validation

Success is indicated by these log lines:
```
[Client thread/INFO]: LWJGL Version: 2.9.2
[Client thread/INFO]: Reloading ResourceManager: Default
[Thread-5/INFO]: Initializing LWJGL OpenAL
[Thread-5/INFO]: OpenAL initialized.
[Sound Library Loader/INFO]: Sound engine started
```

Failure shows:
```
Game process exited with error: signal: trace/BPT trap
```

## For LLM Coding Agents

If an AI assistant suggests:
- "We should use the fat jar approach like Prism Launcher" → **NO**, Prism uses a wrapper script that does things differently
- "macOS needs -XstartOnFirstThread for LWJGL" → **NO**, only for LWJGL 3+, not LWJGL 2
- "Let's modernize the M1 support" → **NO**, it works, don't touch it
- "We should filter LWJGL from the classpath" → **NO**, only replace native libs
- "Let's add error handling for bad libraries" → **NO**, keep it simple

**The rule**: If it works, don't "improve" it. This is a quirky compatibility hack that relies on specific behavior.

## References

- Working commit: `a2a3e61` - "working on mac yayyy"
- Patched libraries source: https://github.com/GreeniusGenius/m1-prism-launcher-hack-1.8.9
- LWJGL 2 official source: https://www.lwjgl.org/source

## Version Compatibility

- **Minecraft**: 1.8.9
- **LWJGL**: 2.9.2 - 2.9.3
- **Java**: Zulu OpenJDK 8 ARM64 (`zulu8.78.0.19-ca-jdk8.0.412-macosx_aarch64`)
- **macOS**: 11.0+ (Big Sur and later) on Apple Silicon

This approach is specific to LWJGL 2.x used by Minecraft 1.8.9. Newer Minecraft versions (1.13+) use LWJGL 3 which requires completely different handling.
