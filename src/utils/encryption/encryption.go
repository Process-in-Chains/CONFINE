package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func EncryptSymmetricKey(symKey []byte, public string) ([]byte, error) {
	// Load the public key from .pem file
	publicKeyFile, err := os.Open(public)
	if err != nil {
		log.Fatal("Error opening public key file:", err)
	}
	defer publicKeyFile.Close()
	pemFileInfo, _ := publicKeyFile.Stat()
	pemBytes := make([]byte, pemFileInfo.Size())
	_, err = publicKeyFile.Read(pemBytes)
	if err != nil {
		log.Fatal("Error reading public key file:", err)
	}
	block, _ := pem.Decode(pemBytes)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("Error parsing public key:", err)
	}
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		log.Fatal("Invalid public key type")
	}
	encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, symKey)
	if err != nil {
		return nil, err
	}
	return encryptedKey, nil
}

func EncryptXES(filePath string, symKey []byte) ([]byte, error) {
	xesData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	// Create a new AES cipher block using the symmetric key
	block, err := aes.NewCipher(symKey)
	if err != nil {
		return nil, err
	}
	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	// Encrypt the XES data using AES in CBC mode
	mode := cipher.NewCBCEncrypter(block, iv)
	paddedData := pkcs7Padding(xesData, aes.BlockSize)
	encryptedData := make([]byte, len(paddedData))
	mode.CryptBlocks(encryptedData, paddedData)
	// Prepend the IV to the encrypted data
	encryptedData = append(iv, encryptedData...)
	return encryptedData, nil
}

func EncryptDataWithSymetric(data []byte, symKey []byte) ([]byte, error) {

	// Create a new AES cipher block using the symmetric key
	block, err := aes.NewCipher(symKey)
	if err != nil {
		return nil, err
	}
	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	// Encrypt the XES data using AES in CBC mode
	mode := cipher.NewCBCEncrypter(block, iv)
	paddedData := pkcs7Padding(data, aes.BlockSize)
	encryptedData := make([]byte, len(paddedData))
	mode.CryptBlocks(encryptedData, paddedData)
	// Prepend the IV to the encrypted data
	encryptedData = append(iv, encryptedData...)
	return encryptedData, nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
func DecryptSymmetricKey(encryptedKey []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	decryptedKey, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedKey)
	if err != nil {
		return nil, err
	}
	return decryptedKey, nil
}
func DecryptXES(encryptedData []byte, symKey []byte) ([]byte, error) {
	// Extract the IV from the encrypted data
	iv := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]
	// Create a new AES cipher block using the symmetric key
	block, err := aes.NewCipher(symKey)
	if err != nil {
		return nil, err
	}
	// Decrypt the XES data using AES in CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)
	decryptedData := make([]byte, len(encryptedData))
	mode.CryptBlocks(decryptedData, encryptedData)

	// Remove the PKCS7 padding from the decrypted data
	decryptedData = pkcs7Unpadding(decryptedData)
	return decryptedData, nil
}
func pkcs7Unpadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
func LoadPublicKeyFromFile(path string) rsa.PublicKey {
	// Read the contents of the PEM file
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading PEM file:", err)
	}
	// Decode the PEM data
	block, _ := pem.Decode(pemData)
	if block == nil {
		fmt.Println("Failed to decode PEM block")
	}
	// Parse the DER-encoded public key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing public key:", err)
	}
	// Assert the type of the public key to RSA
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		fmt.Println("Public key is not an RSA key")
	}
	return *rsaPubKey
}
func LoadPrivateKeyFromFile(keyPath string) (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
func GenerateRandomDecryptionToken() []byte {
	symKeyString := ""
	var symKey []byte
	for strings.ContainsRune(symKeyString, '/') || symKeyString == "" {
		symKey, _ = generateAESKey()
		symKeyString = base64.StdEncoding.EncodeToString(symKey)
	}
	return symKey
}
func generateAESKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
func PublicKeyToString(objectPublicKey *rsa.PublicKey) string {

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&objectPublicKey)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Public key:\n" + base64.StdEncoding.EncodeToString(pubKeyBytes))
	return base64.StdEncoding.EncodeToString(pubKeyBytes)
}
func ParsePublicKeyToString(path string) (string, error) {
	// Read the contents of the PEM file
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading PEM file:", err)
		return "", err
	}
	// Decode the PEM data
	block, _ := pem.Decode(pemData)
	if block == nil {
		fmt.Println("Failed to decode PEM block")
		return "", err
	}
	// Parse the DER-encoded public key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing public key:", err)
		return "", err
	}
	// Assert the type of the public key to RSA
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		fmt.Println("Public key is not an RSA key")
		return "", err
	}
	a := *rsaPubKey
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&a)
	if err != nil {
		return "", err
	}
	fmt.Println("Public key:\n" + base64.StdEncoding.EncodeToString(pubKeyBytes))
	return base64.StdEncoding.EncodeToString(pubKeyBytes), nil

}

func GetReadablePublicKey(pubKey any) ([]byte, bool) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		fmt.Printf("Failed to marshal public key: %s\n", err)
		return nil, true
	}
	print(string(pubBytes))
	// Create a PEM block
	pubPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}

	// Encode the PEM block to a readable format
	pubPemBytes := pem.EncodeToMemory(pubPem)
	return pubPemBytes, false
}
