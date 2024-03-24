package main

type Tls struct {
	Namespaces      []string `yaml:"namespaces,flow"`
	SecretName      string   `yaml:"secretName"`
	VaultSecretName string   `yaml:"vaultSecretName"`
	VaultPath       string   `yaml:"vaultPath"`
}

type Config struct {
	MinDaysBeforeUpdateCert int `yaml:"minDaysBeforeUpdateCert"`
	Vault                   struct {
		Server string `yaml:"server"`
		Token  string `yaml:"token"`
	} `yaml:"vault"`
	Tls struct {
		ListSecrets map[string]Tls `yaml:"tls,inline"`
	}
}
