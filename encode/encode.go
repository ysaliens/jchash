package encode

//Hash libraries
import "crypto/sha512"
import "encoding/base64"

func Encode(pass string) string {
	//Create new Sha512 instance
	hash := sha512.New()
	
	//SHA-512 password & output as byte slice
	hash.Write([]byte(pass))
	passSHA512 := hash.Sum(nil)
	
	// Encode to base 64
	hashed := base64.StdEncoding.EncodeToString([]byte(passSHA512))
	
	return hashed
}