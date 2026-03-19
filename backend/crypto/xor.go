package crypto

import "encoding/base64"

func XorCrypto(data, key string) string {
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ key[i%len(key)] // XOR с ключом по кругу
	}
	// encode as base64 so the result is valid UTF‑8 and safe for storing in text fields
	return base64.StdEncoding.EncodeToString(result)
}
