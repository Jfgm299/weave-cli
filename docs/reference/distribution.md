# Distribution Baseline

## Artifact naming and versioning

Release binaries MUST follow semver and include target OS/arch in the filename:

- `weave_<semver>_<os>_<arch>.tar.gz`
- `weave_<semver>_<os>_<arch>.zip` (optional)

Example:

- `weave_0.1.0_darwin_arm64.tar.gz`

## Signed checksums

For each release:

1. Generate `checksums.txt` for all artifacts.
2. Sign the checksum file producing `checksums.txt.sig`.
3. Publish both files with release artifacts.

## Verification

Users should verify checksum and signature before install when possible.
