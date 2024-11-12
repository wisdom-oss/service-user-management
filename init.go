package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"microservice/internal/config"
	_ "microservice/internal/db" // side effect import to connect to the database and parse the sql queries from it's embed
	"microservice/oidc"
	"microservice/resources"
	"microservice/utils"

	_ "github.com/wisdom-oss/go-healthcheck/client"
)

// init is executed at every startup of the microservice and is always executed
// before main
func init() {
	configureLogger()
	loadCertificates()
	validateOIDCEnvironment()
}

// configureLogger handles the configuration of the logger used in the
// microservice. it reads the logging level from the `LOG_LEVEL` environment
// variable and sets it according to the parsed logging level. if an invalid
// value is supplied or no level is supplied, the service defaults to the
// `INFO` level
func configureLogger() {
	// set the time format to unix timestamps to allow easier machine handling
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// allow the logger to create an error stack for the logs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// now use the environment variable `LOG_LEVEL` to determine the logging
	// level for the microservice.
	rawLoggingLevel, isSet := os.LookupEnv("LOG_LEVEL")

	// if the value is not set, use the info level as default.
	var loggingLevel zerolog.Level
	if !isSet {
		loggingLevel = zerolog.InfoLevel
	} else {
		// now try to parse the value of the raw logging level to a logging
		// level for the zerolog package
		var err error
		loggingLevel, err = zerolog.ParseLevel(rawLoggingLevel)
		if err != nil {
			// since an error occurred while parsing the logging level, use info
			loggingLevel = zerolog.InfoLevel
			log.Warn().Msg("unable to parse value from environment. using info")
		}
	}
	// since now a logging level is set, configure the logger
	zerolog.SetGlobalLevel(loggingLevel)
}

func validateOIDCEnvironment() {
	clientID, isSet := os.LookupEnv("OIDC_CLIENT_ID")
	if !isSet {
		log.Fatal().Msg("OIDC_CLIENT_ID environment variable not set")
	}

	clientSecret, isSet := os.LookupEnv("OIDC_CLIENT_SECRET")
	if !isSet {
		log.Fatal().Msg("OIDC_CLIENT_SECRET environment variable not set")
	}

	issuer, isSet := os.LookupEnv("OIDC_ISSUER")
	if !isSet {
		log.Fatal().Msg("OIDC_ISSUER environment variable not set")
	}

	redirectUri, isSet := os.LookupEnv("OIDC_REDIRECT_URI")
	if !isSet {
		log.Fatal().Msg("OIDC_REDIRECT_URI environment variable not set")
	}

	err := oidc.ExternalProvider.Configure(issuer, clientID, clientSecret, redirectUri)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to configure external OIDC provider information")
	}
}

func loadCertificates() {
	// test if the certificates are already present
	var tries int
loadCerts:
	signingCertificate, err := os.Open(config.SigningCertificateFilePath)
	if err != nil {
		if tries >= 3 {
			log.Fatal().Err(err).Msg("unable to open certificate file after three tries")
		}
		tries++
		if errors.Is(err, os.ErrNotExist) {
			err = utils.GenerateCertificates()
			if err != nil {
				log.Fatal().Err(err).Msg("unable to generate certificates")
			}
			goto loadCerts
		}
	}

	encryptionCertificate, err := os.Open(config.EncryptionCertificateFilePath)
	if err != nil {
		if tries >= 3 {
			log.Fatal().Err(err).Msg("unable to open certificate file after three tries")
		}
		tries++
		if errors.Is(err, os.ErrNotExist) {
			err = utils.GenerateCertificates()
			if err != nil {
				log.Fatal().Err(err).Msg("unable to generate certificates")
			}
			goto loadCerts
		}
	}

	defer signingCertificate.Close()
	certificateContents, err := io.ReadAll(signingCertificate)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to read certificate file")
	}
	privateKeyBlock, _ := pem.Decode(certificateContents)
	if privateKeyBlock.Type != "EC PRIVATE KEY" {
		log.Fatal().Msg("unsupported private key type")
	}

	ecdsaPrivateKey, err := x509.ParseECPrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse ECDSA private key")
	}

	privateKey, err := jwk.FromRaw(ecdsaPrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create jwk private key")
	}

	err = jwk.AssignKeyID(privateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set private key key id")
	}
	err = privateKey.Set(jwk.KeyUsageKey, "sig")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set private key usage type")
	}
	err = privateKey.Set(jwk.AlgorithmKey, jwa.ES256)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set private key algorithm")
	}

	publicKey, err := jwk.PublicKeyOf(privateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse public key")
	}

	resources.PublicSigningKey = publicKey
	resources.PrivateSigningKey = privateKey

	defer encryptionCertificate.Close()
	certificateContents, err = io.ReadAll(encryptionCertificate)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to read certificate file")
	}
	privateKeyBlock, _ = pem.Decode(certificateContents)
	if privateKeyBlock.Type != "EC PRIVATE KEY" {
		log.Fatal().Msg("unsupported private key type")
	}

	ecdsaPrivateKey, err = x509.ParseECPrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse ECDSA private key")
	}

	privateKey, err = jwk.FromRaw(ecdsaPrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create jwk private key")
	}

	err = jwk.AssignKeyID(privateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set private key key id")
	}
	err = privateKey.Set(jwk.KeyUsageKey, "enc")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set private key usage type")
	}
	err = privateKey.Set(jwk.AlgorithmKey, jwa.ES256)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set private key algorithm")
	}

	publicKey, err = jwk.PublicKeyOf(privateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse public key")
	}

	resources.PublicEncryptionKey = publicKey
	resources.PrivateEncryptionKey = privateKey

	resources.KeySet = jwk.NewSet()
	err = resources.KeySet.AddKey(resources.PublicEncryptionKey)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to add public encryption key to key set")
	}
	err = resources.KeySet.AddKey(resources.PublicSigningKey)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to add public signing key to key set")
	}

	log.Info().Msg("loaded certificates")
}
