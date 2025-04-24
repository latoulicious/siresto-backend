package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	// Generate 32 random bytes (256 bits)
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal("Error generating random bytes:", err)
	}

	// Encode to base64
	key := base64.StdEncoding.EncodeToString(bytes)

	fmt.Println("Generated JWT Secret Key:")
	fmt.Println("-------------------------")
	fmt.Println(key)
	fmt.Println("-------------------------")
	fmt.Println("Add this to your .env file as:")
	fmt.Println("JWT_SECRET_KEY=" + key)
}
