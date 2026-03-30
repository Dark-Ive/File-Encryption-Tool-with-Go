# 🔐 File Encryption Tool

A command-line file encryption/decryption tool built in Go using **AES-256 CBC** encryption. Built as a final project for CS50.

---

## Features

- 🔒 AES-CBC encryption with PKCS7 padding
- 🔑 Password-protected decryption
- 📄 Encrypts any file type
- 🖼️ ASCII art banner via `go-figure`
- ⚡ Simple CLI interface powered by `cobra`

---

## Requirements

- Go 1.18+

---

## Installation

```bash
# Clone the repository
git clone https://github.com/Dark-Ive/file-encryption-tool.git
cd file-encryption-tool

# Install dependencies
go mod tidy

# Build the binary
go build -o encrypt .
```

---

## Usage

### Encrypt a file

```bash
./encrypt encrypt <filename>
```

This creates an encrypted file with a `.enc` extension.

**Example:**
```bash
./encrypt encrypt secret.txt
# Output: secret.txt.enc
```

---

### Decrypt a file

```bash
./encrypt decrypt <filename>
```

You will be prompted to enter the password (up to 3 attempts).

**Example:**
```bash
./encrypt decrypt secret.txt.enc
# Prompts: Password:
# Output: Decrypted file written back to secret.txt.enc
```


---

## How It Works

1. **Encryption:**
   - Reads the file contents
   - Applies PKCS7 padding to align with AES block size (16 bytes)
   - Generates a random 16-byte IV (Initialization Vector)
   - Encrypts using AES-CBC
   - Prepends the IV to the ciphertext and Base64-encodes the result
   - Saves as `<filename>.enc`

2. **Decryption:**
   - Reads and Base64-decodes the `.enc` file
   - Extracts the IV from the first 16 bytes
   - Decrypts using AES-CBC with the provided key
   - Strips PKCS7 padding
   - Writes the plaintext back to the file

---

## Dependencies

| Package | Purpose |
|---|---|
| [`github.com/spf13/cobra`](https://github.com/spf13/cobra) | CLI framework |
| [`github.com/common-nighthawk/go-figure`](https://github.com/common-nighthawk/go-figure) | ASCII art banner |

---

## Security Notes

> ⚠️ This project is for **educational purposes** (CS50 final project).

- The encryption key is currently hardcoded as `CS50secretpasswd`. In a production tool, keys should be derived from user passwords using a KDF like **PBKDF2** or **Argon2**.
- CBC mode is secure when used correctly with a random IV (as done here), but **AES-GCM** is generally preferred for modern applications as it also provides authentication.

---

## License

MIT License - feel free to use, modify, and distribute.

---

Built with ❤️ as a CS50 Final Project.
