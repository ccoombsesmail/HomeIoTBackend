package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
)

const (
	PrivateKeyPath      string = "rsa_4096_priv.pem.txt"
	MessageKeyDelimiter string = ":::"
)

// use a variable to store the decoded private key
var PrivateKey *rsa.PrivateKey

func init() {

	// read the private key for data decryption
	keyFile, err := ioutil.ReadFile(PrivateKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
	}
	// get the private key pem block data
	block, _ := pem.Decode(keyFile)
	if block == nil {
		log.Fatal("Error reading private key")
	}
	// decode the RSA private key
	PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal("Could not decode our")
	}
}

func main() {
	// Set the router as the default one shipped with Gin

	router := gin.Default()
	hasLoaded := false

	// Serve frontend static login page file

	router.StaticFS("/hi/", http.Dir("./login/build"))

	router.POST("/", func(c *gin.Context) {

		//fmt.Printf("BODY: %s", c.Request.Header)

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {

		}
		//fmt.Printf("%s\n", string(body))

		encrypted := strings.Split(string(body), MessageKeyDelimiter)

		// decrypt the AES-encryption-key with our private key
		keyAndIv := decryptRSA(encrypted[0])

		// split the key and its components
		keyComponents := strings.Split(keyAndIv, MessageKeyDelimiter)

		// decrypt the AES encrypted data using
		// the components from before
		message := decryptAES(keyComponents[0], keyComponents[1], encrypted[1])

		//fmt.Printf("MESSAGE: %s\n", message)

		if string(bytes.TrimRight(message, "")) == "17b84c8330f86af407ec45cd1ac3e9bc183d38d3c13d64ff06fbd699ccb3c69e" {
			if hasLoaded == false {
				router.StaticFS("/dash", http.Dir("./login/build1"))
				hasLoaded = true
			}
		} else {
			c.Writer.WriteHeader(403)
		}

		//fmt.Printf("%s\n", string(result))
	})

	log.Fatal(http.ListenAndServe(
		":5000",
		handlers.CORS(handlers.AllowedHeaders(
			[]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router),
	))

}

func decryptRSA(encrypted string) string {

	// decode  encrypted string into cipher bytes
	cipheredValue, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		log.Println("error: decoding string (rsa)")
		return ""
	}
	// decrypt the data
	var out []byte
	out, err = rsa.DecryptPKCS1v15(rand.Reader, PrivateKey, cipheredValue)
	if err != nil {
		log.Println("error: reading encrypted data")
		return ""
	}
	return string(out)
}

func decryptAES(keyString string, ivString string, encrypted string) []byte {

	// decode from hex to byte
	key, _ := hex.DecodeString(keyString)
	iv, _ := hex.DecodeString(ivString)
	// decode our encrypted string into bytes
	cipheredMessage, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		log.Println("error: decoding string (aes)")
	}
	// create a new cipher block from our key
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	// cbc message must be multiple of aes blocksize
	if len(cipheredMessage) < aes.BlockSize {
		log.Println("error: ciphertext too short")
	}
	// the iv is prepended to the actual message/data
	cipherText := cipheredMessage[aes.BlockSize:]
	// decrypt the data
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	return cipherText
}
