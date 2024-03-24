/*
Программа для обновления секретов в кластере кубернетес
Значение для обновления берет из vault хранилища

*/

package main

import (
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"log"
	"os"
	"time"
)

var configuration = Config{}
var secretsClient coreV1Types.SecretInterface

func main() {

	//read configuration
	cfg, err := NewConfig(os.Getenv("PWD") + "/config.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}
	configuration = *cfg

	for {
		for blockName, blockValues := range configuration.Tls.ListSecrets {
			log.Println("Checking ", blockName, "block")
			for _, ns := range blockValues.Namespaces {
				log.Println("Work with ", ns, "namespace")
				secretsClient = initClient(ns)
				secrets, err := readK8sSecret(ns, blockValues.SecretName)
				if err != nil {
					log.Println(err)
					continue
				}
				days, err := checkCertExpiredTime(secrets["tls.crt"])
				if err != nil {
					log.Println(err)
				}
				if days <= configuration.MinDaysBeforeUpdateCert {
					log.Println("Secret ", blockValues.SecretName, "is need update!")
					log.Println("Read secret from vault")
					data, err := readVaultSecret(blockValues.VaultSecretName, blockValues.VaultPath)
					if err != nil {
						log.Println(err)
					}
					cert := data["tls.crt"]
					vaultCertExpiredTime, err := checkCertExpiredTime(cert)
					if err != nil {
						log.Println(err)
					}
					if cert == secrets["tls.crt"] {
						log.Println("The certificate in vault is equal to the certificate in k8s. Skip...")
						continue
					}
					if vaultCertExpiredTime <= configuration.MinDaysBeforeUpdateCert {
						log.Println("Cert in vault expired by ", vaultCertExpiredTime, "days. Skip...")
						continue
					}
					log.Println("Update cert in k8s")
					err = updateK8sSecret(ns, blockValues.SecretName, data)
					if err != nil {
						log.Println(err)
					}
				} else {
					log.Println("The secret ", blockValues.SecretName, " does not require renewal")
				}
			}

		}
		time.Sleep(time.Duration(configuration.TimeToSleep) * time.Minute)
	}
}
