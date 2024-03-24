/*
Программа для обновления секретов в кластере кубернетес
Значение для обновления берет из vault хранилища

*/

package main

import (
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"log"
)

var configuration = Config{}
var secretsClient coreV1Types.SecretInterface

func main() {

	//read configuration
	cfg, err := NewConfig("config.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}
	configuration = *cfg

	//Пробегаемся по словарю с конфига и проверяем секреты
	for a, b := range configuration.Tls.ListSecrets {
		log.Println("Checking ", a, "block")
		for _, ns := range b.Namespaces {
			log.Println("Work with ", ns, "namespace")
			secretsClient = initClient(ns)
			secrets, err := readK8sSecret(ns, b.SecretName)
			if err != nil {
				log.Println(err)
				continue
			}
			days, err := checkCertExpiredTime(secrets["tls.crt"])
			if err != nil {
				log.Println(err)
			}
			if days <= configuration.MinDaysBeforeUpdateCert {
				log.Println("Secret ", b.SecretName, "is need update!")
				log.Println("Read secret from vault")
				data, err := readVaultSecret(b.VaultSecretName, b.VaultPath)
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
				err = updateK8sSecret(ns, b.SecretName, data)
				if err != nil {
					log.Println(err)
				}
			}
		}

		//create vault client
		//var secrets = map[string]string{}
		//secrets, err = readK8sSecret("argocd", "argocd-secret")
		//if err != nil {
		//	log.Println(err)
		//}
		//var daysBefore, error = checkCertExpiredTime(secrets["tls.crt"])
		//if error != nil {
		//	log.Println(error)
		//}
		//fmt.Println(daysBefore)
		//var keyPem = string(secret.Data["tls.key"])

		//fmt.Println(cert)
		//fmt.Println(cert.Leaf.NotAfter)
		//fmt.Println(cert.Leaf.NotBefore)

	}
}