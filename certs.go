package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"time"
)

// Проверяет сколько времени (в днях) действует сертификат
func checkCertExpiredTime(certPem string) (int, error) {

	block, _ := pem.Decode([]byte(certPem))
	if block == nil {
		return 0, errors.New("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return 0, errors.New("\"failed to parse certificate: \" + err.Error()")
	}
	var days = int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
	log.Print("Days before expired cert is:", days)
	return days, nil

}
