/*
Package main - создание сертификата X.509, а также пары приватного и публичного ключей в configs/.
Пример из: Спринт 9 → Тема 2/4: Генерация случайных чисел → Урок 1/3: Пакеты math/rand и crypto/rand.
*/

package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	dir := "./configs"
	certPEMFile, err := os.Create(filepath.Join(dir, "tls.pem"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = certPEMFile.Write(certPEM.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	privateKeyPEMFile, err := os.Create(filepath.Join(dir, "tls.key"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = privateKeyPEMFile.Write(privateKeyPEM.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}
