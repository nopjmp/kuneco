package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base32"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"cwtch.im/cwtch/connectivity/tor"
	"golang.org/x/net/proxy"
)

func generateOnionAddress() (string, string) {
	const bits = 10
	const PEM = "RSA PRIVATE KEY"
	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	pub, _ := asn1.Marshal(key.PublicKey)
	hashBytes := sha1.Sum(pub)
	hash := base32.StdEncoding.EncodeToString(hashBytes[:])
	exportedPriv := &pem.Block{
		Type:  PEM,
		Bytes: x509.MarshalPKCS1PrivateKey(key)}
	privateKey := pem.EncodeToMemory(exportedPriv)
	return strings.ToLower(hash[:16] + ".onion"), string(privateKey)
}

func startTor(torPath string) error {
	torrc := path.Join(savePath, "tor", "torrc")

	if _, err := os.Stat(torrc); os.IsNotExist(err) {
		log.Printf("> Writing torrc to: %v\n", torrc)
		file, err := os.Create(torrc)
		if err != nil {
			return err
		}
		fmt.Fprintf(file, "SOCKSPort %d\nControlPort %d\nCookieAuthentication 0\nSafeSocks 1\nHiddenServiceDir %v\nHiddenServicePort 80 127.0.0.1:48486", 9050, 9051, savePath+loadedUserProfile.Username)
		file.Close()
	}
	tm, err := tor.NewTorManager(9050, 9051, torPath, torrc)
	if err != nil {
		return err
	}

	app.TorManager = tm
	return nil

}

func sendTorRequest(method string, URL string) ([]byte, error) {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	timeout := time.Duration(1 * time.Second)
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport, Timeout: timeout}

	httpTransport.Dial = dialer.Dial
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return buf.Bytes(), nil
	}
	log.Println(URL, "UNAVAILABLE")
	return nil, nil
}
