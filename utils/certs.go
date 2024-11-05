package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	"microservice/internal/config"
	"microservice/internal/db"
)

func GenerateCertificates() error {
	otherServiceGenerating, _ := db.Redis.Get(context.Background(), "ums-is-generating-jwk").Bool()
	if otherServiceGenerating {
		for {
			stillGenerating, _ := db.Redis.Get(context.Background(), "ums-is-generating-jwk").Bool()
			if stillGenerating {
				time.Sleep(250 * time.Millisecond)
				continue
			}
			break
		}
		return nil
	}
	_ = db.Redis.Set(context.Background(), "ums-is-generating-jwk", true, 0)
	//defer db.Redis.Del(context.Background(), "ums-is-generating-jwk")

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}

	privateBlock := pem.Block{}
	privateBlock.Type = "EC PRIVATE KEY"
	privateBlock.Bytes = privateKeyBytes

	err = os.MkdirAll("./.certs", 0600)
	certificateFile, err := os.OpenFile(config.CertificateFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	defer certificateFile.Close()
	if err != nil {
		return err
	}

	err = pem.Encode(certificateFile, &privateBlock)
	if err != nil {
		return err
	}
	return nil
}
