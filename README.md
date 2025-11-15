<p align="center">
  <img width="500" height="500" alt="Image" src="https://github.com/user-attachments/assets/11b193f0-b9ed-45e9-a7fa-923f0c8c47ae" />
</p>

This is a helper program for the [BlindFold](https://chromewebstore.google.com/detail/blindfold/cmgjlocmnppfpknaipdfodjhbplnhimk?hl=ko&utm_source=ext_sidebar) Chrome extension's E2EE configuration that helps store the device key in the Windows Credential Manager on Windows and in the macOS Keychain.

## Download

- [Windows](https://github.com/personalconnect/blindfold-keeper/releases/download/v0.0.1/Blindfold-Keeper.exe)
- [MacOS](https://github.com/personalconnect/blindfold-keeper/releases/download/v0.0.1/Blindfold-Keeper.pkg)

## Output

### MacOS

- /Library/Application\ Support/Google/Chrome/NativeMessagingHosts/com.blindfold.keeper.json
- /Library/Application Support/Blindfold/mac-blindfold-keeper

### Windows (64bit)

- C:\Program Files\Blindfold
  - com.blindfold.keeper.json
  - unins000.exe
  - unins000.dat
  - win-blindfold-keeper.exe

### Windows (32bit)

- C:\Program Files (x86)\Blindfold
  - com.blindfold.keeper.json
  - unins000.exe
  - unins000.dat
  - win-blindfold-keeper.exe