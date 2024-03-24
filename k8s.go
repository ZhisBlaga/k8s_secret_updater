package main

import (
	"context"
	"errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func initClient(namespase string) coreV1Types.SecretInterface {
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset.CoreV1().Secrets(namespase)
}

// Читает сертифткат и ключ из секрета (должны лежать в tls.crt & tls.key соответвенно
func readK8sSecret(namespace, secretName string) (map[string]string, error) {

	var secretsClient = initClient(namespace)
	log.Printf("Reading secret %s in namespace %s\n", secretName, namespace)
	secret, err := secretsClient.Get(context.TODO(), secretName, metaV1.GetOptions{})
	if err != nil {
		return nil, errors.New(err.Error())
	}
	dict := make(map[string]string)

	dict["tls.crt"] = string(secret.Data["tls.crt"])
	dict["tls.key"] = string(secret.Data["tls.key"])

	if dict["tls.crt"] == "" || dict["tls.key"] == "" {
		return nil, errors.New("tls.crt or/and tls.key is empty")
	}
	return dict, nil
}
func updateK8sSecret(namespace, secretName string, vaultData map[string]string) error {
	var secretsClient = initClient(namespace)
	secret, err := secretsClient.Get(context.TODO(), secretName, metaV1.GetOptions{})
	secret.Data["tls.key"] = []byte(vaultData["tls.key"])
	secret.Data["tls.crt"] = []byte(vaultData["tls.crt"])
	updatedSecret, err := secretsClient.Update(context.TODO(), secret, metaV1.UpdateOptions{})
	if err != nil {
		return errors.New(err.Error())
	}
	log.Println(updatedSecret.Name)
	return nil
}
