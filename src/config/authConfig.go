package config

import "os"

var JwtEncodingKey = ""
var JwtIssuer = "FGK_PASMAS_backend"

func LoadAuthConfig() {
    encKey := os.Getenv("JWT_ENCODING")
    if encKey == "" {
        panic("JWT_ENCODING is not set")
    }
    JwtEncodingKey = encKey

    issuer := os.Getenv("JWT_ISSUER")
    if issuer != "" {
        JwtIssuer = issuer
    }
}
