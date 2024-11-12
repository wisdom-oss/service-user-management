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
	defer db.Redis.Del(context.Background(), "ums-is-generating-jwk")

	privateSigningKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	privateSigningKeyBytes, err := x509.MarshalECPrivateKey(privateSigningKey)
	if err != nil {
		return err
	}

	privateSigningBlock := pem.Block{}
	privateSigningBlock.Type = "EC PRIVATE KEY"
	privateSigningBlock.Bytes = privateSigningKeyBytes

	err = os.MkdirAll("./.certs", 0600)
	if err != nil {
		return err
	}

	signingCertificateFile, err := os.OpenFile(config.SigningCertificateFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	err = pem.Encode(signingCertificateFile, &privateSigningBlock)
	if err != nil {
		return err
	}

	err = signingCertificateFile.Close()
	if err != nil {
		return err
	}

	privateEncryptionKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	privateEncryptionKeyBytes, err := x509.MarshalECPrivateKey(privateEncryptionKey)
	if err != nil {
		return err
	}

	privateEncryptionBlock := pem.Block{}
	privateEncryptionBlock.Type = "EC PRIVATE KEY"
	privateEncryptionBlock.Bytes = privateEncryptionKeyBytes

	encryptionCertificateFile, err := os.OpenFile(config.EncryptionCertificateFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	err = pem.Encode(encryptionCertificateFile, &privateEncryptionBlock)
	if err != nil {
		return err
	}

	err = encryptionCertificateFile.Close()
	return err
}
