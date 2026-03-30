package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var rootcmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Final project for CS50",
}

var encrypt = &cobra.Command{
	Use:   "encrypt [file]",
	Short: "Encrypt a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := []byte("CS50secretpasswd")
		filename := args[0]
		err := encryptfile(key, filename)
		if err != nil {
			fmt.Println("Error Encrypting file:", err)
			return
		}
		fmt.Println("File Encrypted successfully!")
	},
}

var decrypt = &cobra.Command{
	Use:   "decrypt [file]",
	Short: "Decrypt a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := []byte("CS50secretpasswd")
		filename := args[0]
		var password []byte
		i := 0
		for i < 3 {
			fmt.Println("Password: ")
			fmt.Scanln(&password)
			if !bytes.Equal(password, key) {
				fmt.Println("Wrong password! Try again")
				i++
			} else {
				break
			}
		}
		if i == 3 {
			fmt.Println("Failed to provide the correct password after three attempts")
			return
		}
		err := decryptfile(key, filename)
		if err != nil {
			fmt.Println("Error decrypting file:", err)
			return
		}
		fmt.Println("File decrypted successfully!")
	},
}

// Encryption method - CBC encryption
func encryptfile(key []byte, filename string) error {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return fmt.Errorf("invalid key size: %d", len(key))
	}

	plaintext, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Pad plaintext to a multiple of block size
	paddingSize := aes.BlockSize - len(plaintext)%aes.BlockSize
	padding := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	plaintext = append(plaintext, padding...)

	// Generate random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	ciphertext := make([]byte, len(plaintext))
	cfb := cipher.NewCBCEncrypter(block, iv)
	cfb.CryptBlocks(ciphertext, plaintext)

	//Prepend IV to ciphertext before encoding
	combined := append(iv, ciphertext...)
	encoded_cipher := base64.StdEncoding.EncodeToString(combined)

	return os.WriteFile(filename+".enc", []byte(encoded_cipher), 0644)
}

// Decryption method - CBC decryption
func decryptfile(key []byte, filename string) error {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return fmt.Errorf("invalid key size: %d (must be 16, 24, or 32)", len(key))
	}

	ciphertext, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	decoded_cipher, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return err
	}

	if len(decoded_cipher) < aes.BlockSize*2 {
		return fmt.Errorf("ciphertext too short: file may not be encrypted or is corrupted")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Extract IV from the first 16 bytes
	iv := decoded_cipher[:aes.BlockSize]
	ciphertext = decoded_cipher[aes.BlockSize:]

	// CBC requires ciphertext to be a multiple of block size
	if len(ciphertext)%aes.BlockSize != 0 {
		return fmt.Errorf("ciphertext is not a multiple of block size: file may be corrupted")
	}

	plaintext := make([]byte, len(ciphertext))
	cfb := cipher.NewCBCDecrypter(block, iv)
	cfb.CryptBlocks(plaintext, ciphertext)

	// Remove padding
	paddingSize := int(plaintext[len(plaintext)-1])
	if paddingSize > aes.BlockSize || paddingSize == 0 {
		return fmt.Errorf("invalid padding size: %d", paddingSize)
	}
	plaintext = plaintext[:len(plaintext)-paddingSize]

	return os.WriteFile(filename, plaintext, 0644)
}

func main() {
	figure.NewColorFigure("File Encryption Tool", "", "gray", true).Print()
	rootcmd.AddCommand(encrypt)
	rootcmd.AddCommand(decrypt)

	if err := rootcmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
