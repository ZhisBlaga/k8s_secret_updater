/*
Программа для обновления секретов в кластере кубернетес
Значение для обновления берет из vault хранилища

*/

package main

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"os"
	"time"
)

var configuration = Config{}
var secretsClient coreV1Types.SecretInterface
var log = logrus.New()

func init() {
	//read configuration
	cfg, err := NewConfig(os.Getenv("PWD") + "/config.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}
	configuration = *cfg

	//init logger
	// Log as JSON instead of the default ASCII formatter.
	log.Formatter = new(prefixed.TextFormatter)
	log.Level = logrus.DebugLevel

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

}

func main() {

	for {
		for blockName, blockValues := range configuration.Tls.ListSecrets {
			log.Info("Processing tls block ", blockName)
			for _, ns := range blockValues.Namespaces {
				logger := log.WithFields(logrus.Fields{"block": blockName, "ns": ns})
				secretsClient = initClient(ns)
				secrets, err := readK8sSecret(ns, blockValues.SecretName)
				if err != nil {
					logger.Error(err)
					continue
				}
				days, err := checkCertExpiredTime(secrets["tls.crt"])
				if err != nil {
					logger.Error(err)
				}
				if days <= configuration.MinDaysBeforeUpdateCert {
					logger.Info("Secret ", blockValues.SecretName, "is need update!")

					logger.Debug("Attempting to receive data from vault")
					data, err := readVaultSecret(blockValues.VaultSecretName, blockValues.VaultPath)
					if err != nil {
						logger.Error(err)
					}
					logger.Debug("Data retrieved from vault")
					cert := data["tls.crt"]
					vaultCertExpiredTime, err := checkCertExpiredTime(cert)
					if err != nil {
						logger.Error(err)
					}
					if cert == secrets["tls.crt"] {
						logger.Info("The certificate in vault is equal to the certificate in k8s. Skip...")
						continue
					}
					if vaultCertExpiredTime <= configuration.MinDaysBeforeUpdateCert {
						logger.Info("Cert in vault expired by ", vaultCertExpiredTime, " days. Skip...")
						continue
					}
					logger.Info("Attempting updating the secret in k8s")
					err = updateK8sSecret(ns, blockValues.SecretName, data)
					if err != nil {
						logger.Error(err)
						continue
					}
					logger.Info("Secret updated")
				} else {
					logger.Info("The secret ", blockValues.SecretName, " does not require renewal")
				}
			}

		}
		time.Sleep(time.Duration(configuration.TimeToSleep) * time.Minute)
	}
}
