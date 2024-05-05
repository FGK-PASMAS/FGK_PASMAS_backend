package config

import "os"

type TlsConfig struct {
    CertPath string
    KeyPath string
}

func LoadTlsConfig() *TlsConfig {
    tlsConfig := new(TlsConfig)
    tlsConfig.CertPath = os.Getenv("TLS_CERT_PATH")
    tlsConfig.KeyPath = os.Getenv("TLS_KEY_PATH")

    if tlsConfig.CertPath == "" || tlsConfig.KeyPath == ""{
        return nil
    }

    return tlsConfig
}
