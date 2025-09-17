# Asset Ranking System

The system implements a ranked asset selection that:

1. Ignores Source Archives
   Filters out source.tar.gz, v1.54.2.tar.gz, and similar source code archives
   Prevents downloading source code instead of binaries
2. Priority Ranking (highest to lowest)
   Standalone Binaries (+100): .exe files or platform-specific binaries
   Archives (+50): .tar.gz, .zip, .dmg files
   Deb Packages (+25): .deb files for Ubuntu/Debian systems
3. Platform Compatibility
   Must match current architecture (amd64/x86_64/x64)
   Must be compatible with current OS (linux/darwin/windows)
   Bonus points for exact OS match (+5)
   Special Ubuntu bonus for .deb packages (+15)
4. Best Asset Selection
   SelectBestAsset() function ranks all available assets
   Returns the asset with the highest rank
   Used in GithubDownloadLatestRelease() instead of first-match logic
   The system successfully selected golangci-lint-1.54.2-linux-amd64.tar.gz (rank 65) over less optimal options like the .deb package (rank 55) or standalone binary (rank 15), and completely ignored source archives (rank 0).
