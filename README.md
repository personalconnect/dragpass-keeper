<p align="center">
  <img width="500" height="500" alt="Image" src="https://github.com/user-attachments/assets/11b193f0-b9ed-45e9-a7fa-923f0c8c47ae" />
</p>

This helper program supports the E2EE configuration for the [BlindFold Chrome Extension](https://chromewebstore.google.com/detail/blindfold/cmgjlocmnppfpknaipdfodjhbplnhimk?hl=ko&utm_source=ext_sidebar) **v1.0.2** by securely storing device keys in OS-native encrypted vaults:

- **macOS**: Keychain
- **Linux**: Secret Service API (GNOME Keyring / KDE Wallet)
- **Windows**: Credential Manager

## Download

Download the latest release from the [Releases page](https://github.com/personalconnect/dragpass-keeper/releases).

### Available Packages

- **macOS**:
  - `dragpass-keeper-macos-x86_64.pkg` (Intel)
  - `dragpass-keeper-macos-arm64.pkg` (Apple Silicon)
- **Linux**:
  - `dragpass-keeper-linux-x86_64.deb` (x86_64/amd64)
  - `dragpass-keeper-linux-arm64.deb` (ARM64)
- **Windows**: `dragpass-keeper.exe` (x64 installer)

## Verifying Downloads

All release packages are signed with GPG for security. We strongly recommend verifying the integrity of downloaded files.

### 1. Import the Public Key

```bash
# Download and import the public key
curl https://raw.githubusercontent.com/personalconnect/dragpass-keeper/main/GPG_PUBLIC_KEY.asc | gpg --import
```

Or import manually from [GPG_KEYSPUBLIC_KEY.asc](GPG_PUBLIC_KEY.asc).

**Key Fingerprint**: `66DF 4017 8A5F 6F66 EAAF 318A 3FC4 1856 9192 8FDC`

### 2. Verify the Signature

```bash
# For macOS (Intel)
gpg --verify dragpass-keeper-macos-x86_64.pkg.sig dragpass-keeper-macos-x86_64.pkg

# For macOS (Apple Silicon)
gpg --verify dragpass-keeper-macos-arm64.pkg.sig dragpass-keeper-macos-arm64.pkg

# For Linux (x86_64)
gpg --verify dragpass-keeper-linux-x86_64.deb.sig dragpass-keeper-linux-x86_64.deb

# For Linux (ARM64)
gpg --verify dragpass-keeper-linux-arm64.deb.sig dragpass-keeper-linux-arm64.deb

# For Windows
gpg --verify dragpass-keeper.exe.sig dragpass-keeper.exe
```

You should see output like:
```
gpg: Good signature from "JinHyeok Hong <vjinhyeokv@gmail.com>" [ultimate]
```

## Installation Output

### macOS

After installing the `.pkg` file, the following files are created:

- `/Library/Application Support/DragPass/dragpass-keeper` - Main executable
- `/Library/Application Support/Google/Chrome/NativeMessagingHosts/com.dragpass.keeper.json` - Chrome Native Messaging manifest

**Key Storage**: macOS Keychain

### Linux

After installing the `.deb` file, the following files are created:

- `/opt/dragpass/dragpass-keeper` - Main executable
- `/etc/opt/chrome/native-messaging-hosts/com.dragpass.keeper.json` - Chrome manifest
- `/etc/chromium/native-messaging-hosts/com.dragpass.keeper.json` - Chromium manifest

**Key Storage**: Secret Service API (GNOME Keyring / KDE Wallet)

### Windows

After running the `.exe` installer, the following files are created:

**64-bit System:**
- `C:\Program Files\DragPass\`
  - `dragpass-keeper.exe` - Main executable
  - `com.dragpass.keeper.json` - Chrome Native Messaging manifest
  - `unins000.exe` - Uninstaller
  - `unins000.dat` - Uninstaller data

**32-bit System:**
- `C:\Program Files (x86)\DragPass\`
  - `dragpass-keeper.exe` - Main executable
  - `com.dragpass.keeper.json` - Chrome Native Messaging manifest
  - `unins000.exe` - Uninstaller
  - `unins000.dat` - Uninstaller data

**Key Storage**: Windows Credential Manager