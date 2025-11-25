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

## API Reference

DragPass Keeper communicates with the Chrome extension via Native Messaging protocol. All messages use an **envelope pattern** for better type safety and extensibility.

### Message Format

**Request (Envelope Pattern):**
```json
{
  "action": "action_name",
  "payload": {
    // action-specific fields
  }
}
```

**Success Response:**
```json
{
  "success": true,
  "data": {
    // action-specific response data
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "error message"
}
```

---

### Health Check

#### `ping` - Health Check

Check if the DragPass Keeper is running and responsive.

**Request:**
```json
{
  "action": "ping"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "version": "0.0.5",
    "hash": "binary_sha256_hash",
    "path": "/path/to/dragpass-keeper"
  }
}
```

---

### Device Key Management

#### `savedevicekey` - Save Device Key

Stores the device encryption key in the OS keystore.

**Request:**
```json
{
  "action": "savedevicekey",
  "payload": {
    "key": "base64_encoded_device_key"
  }
}
```

**Response:**
```json
{
  "success": true
}
```

---

#### `getdevicekey` - Get Device Key

Retrieves the stored device encryption key.

**Request:**
```json
{
  "action": "getdevicekey"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "key": "base64_encoded_device_key"
  }
}
```

---

#### `deletedevicekey` - Delete Device Key

Removes the device encryption key from the keystore.

**Request:**
```json
{
  "action": "deletedevicekey"
}
```

**Response:**
```json
{
  "success": true
}
```

---

### Keypair Management

#### `generatekeypair` - Generate RSA Keypair

Generates a new RSA-2048 keypair for the Helper. Requires server signature verification.

**Request:**
```json
{
  "action": "generatekeypair",
  "payload": {
    "challenge_token": "server_provided_challenge_token",
    "signature": "base64_server_signature"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "publickey": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
  }
}
```

**Notes:**
- Verifies the signature using the server's public key
- Deletes existing session code and keypair before generating new one
- Stores both private and public keys in the OS keystore

---

#### `getpublickey` - Get Helper Public Key

Retrieves the Helper's public key.

**Request:**
```json
{
  "action": "getpublickey"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "publickey": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
  }
}
```

---

#### `getserverpubkey` - Get Server Public Key

Retrieves the server's public key that is stored in the OS keystore.

**Request:**
```json
{
  "action": "getserverpubkey"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "publickey": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
  }
}
```

**Notes:**
- The server public key is hardcoded in the binary and initialized on first run
- This key is used to verify signatures from the server
- Stored in OS-native keystore for retrieval

---

### Session Code Management

#### `savesessioncode` - Save Encrypted Session Code

Decrypts and stores the session code. Used during signup.

**Request:**
```json
{
  "action": "savesessioncode",
  "payload": {
    "encrypted_session_code": "base64_encrypted_session_code",
    "signature": "base64_server_signature"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "session_code": "decrypted_session_code"
  }
}
```

**Process:**
1. Verifies signature using server's public key
2. Decrypts the session code using Helper's private key (RSA-OAEP with SHA-256)
3. Stores the decrypted session code in the OS keystore
4. Returns the decrypted session code

---

#### `getsessioncode` - Get Session Code

Retrieves the stored session code.

**Request:**
```json
{
  "action": "getsessioncode"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "session_code": "stored_session_code"
  }
}
```

---

### Signup Flow

#### `signalias` - Sign User Alias

Signs the user alias with Helper's private key. Used during signup.

**Request:**
```json
{
  "action": "signalias",
  "payload": {
    "alias": "user_alias"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "signature": "base64_signature",
    "publickey": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
  }
}
```

**Process:**
1. Signs the alias using Helper's private key (RSA PKCS#1 v1.5 with SHA-256)
2. Returns the signature and Helper's public key

---

### Login Flow

#### `signaliaswithtimestamp` - Sign Alias with Timestamp

Signs the user alias with current timestamp. Used for login authentication.

**Request:**
```json
{
  "action": "signaliaswithtimestamp",
  "payload": {
    "alias": "user_alias"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "signature": "base64_signature",
    "timestamp": 1234567890
  }
}
```

**Process:**
1. Generates current Unix timestamp
2. Creates payload: `"alias:timestamp"`
3. Signs the payload using Helper's private key
4. Returns the signature and timestamp

---

#### `signchallengetoken` - Sign Challenge Token

Verifies and signs a challenge token. Used for login verification.

**Request:**
```json
{
  "action": "signchallengetoken",
  "payload": {
    "challenge_token": "server_challenge_token",
    "signature": "base64_server_signature"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "signature": "base64_helper_signature"
  }
}
```

**Process:**
1. Verifies the server's signature on the challenge token using server's public key
2. Signs the challenge token using Helper's private key
3. Returns the Helper's signature

---

## Cryptographic Details

### Key Formats
- **RSA Key Size**: 2048 bits
- **Private Key Format**: PKCS#8 PEM
- **Public Key Format**: PKIX PEM

### Algorithms
- **Signature Algorithm**: RSA PKCS#1 v1.5 with SHA-256
- **Encryption Algorithm**: RSA-OAEP with SHA-256
- **Hash Function**: SHA-256

### Key Storage Locations

**macOS Keychain:**
```
Service: com.dragpass.keeper
Items:
- DragPassServerPublicKey
- DragPassKeeperPrivateKey
- DragPassKeeperPublicKey
- DeviceKey
- SessionCode
```

**Linux Secret Service:**
```
Collection: default keyring
Schema: com.dragpass.keeper
Items:
- DragPassServerPublicKey
- DragPassKeeperPrivateKey
- DragPassKeeperPublicKey
- DeviceKey
- SessionCode
```

**Windows Credential Manager:**
```
Target Prefix: com.dragpass.keeper
Credentials:
- DragPassServerPublicKey
- DragPassKeeperPrivateKey
- DragPassKeeperPublicKey
- DeviceKey
- SessionCode
```
