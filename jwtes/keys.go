// package to manage keys
package jwtes

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	DEFAULT_KEY_TYPE     = "ed25519"
	DEFAULT_KEY_NAME     = "key"
	DEFAULT_KEY_LOCATION = ""
)

type keyManager struct {
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
	name    string
}

func NewKeys(name string) (*keyManager, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if name == "" {
		name = DEFAULT_KEY_NAME
	}
	if err != nil {
		return nil, fmt.Errorf("creating keys: %w", err)
	}
	km := keyManager{
		privKey: priv,
		pubKey:  pub,
		name:    name,
	}
	km.writeToDisk()
	km.readFromDisk()
	return &km, nil
}

// writeToDisk
// using https://gist.github.com/rorycl/d300f3ab942fd79e6cc1f37db0c6260f
func (k *keyManager) writeToDisk() error {
	// converts a private key to PKCS #8, ASN.1 DER form.
	b, err := x509.MarshalPKCS8PrivateKey(k.privKey)
	if err != nil {
		return fmt.Errorf("marshaling key %s: %w", "privKey", err)
	}
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	}
	// write priv key to Disk
	if err := os.WriteFile(DEFAULT_KEY_LOCATION+DEFAULT_KEY_NAME, pem.EncodeToMemory(block), 0600); err != nil {
		return fmt.Errorf("writing file to disk %s: %w", DEFAULT_KEY_LOCATION+DEFAULT_KEY_NAME, err)
	}
	// converts public key to PKCS #8
	b, err = x509.MarshalPKIXPublicKey(k.pubKey)
	if err != nil {
		return fmt.Errorf("marshaling key %s: %w", "pubKey", err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}
	// write public key to disk
	if err := os.WriteFile(DEFAULT_KEY_LOCATION+DEFAULT_KEY_NAME+".pub", pem.EncodeToMemory(block), 0644); err != nil {
		return fmt.Errorf("writing file to disk %s: %w", DEFAULT_KEY_LOCATION+DEFAULT_KEY_NAME+".pub", err)
	}
	return err
}

func (k *keyManager) readFromDisk() error {
	var (
		ok bool
	)
	// read file
	file, err := os.ReadFile(DEFAULT_KEY_LOCATION + DEFAULT_KEY_NAME)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", DEFAULT_KEY_LOCATION+DEFAULT_KEY_NAME, err)
	}
	// parse file
	block, _ := pem.Decode(file)
	if block == nil {
		return fmt.Errorf("decoding pem from priv file %s", DEFAULT_KEY_LOCATION+DEFAULT_KEY_NAME)
	}
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parsing key to ec25519: %w", err)
	}
	// assign key to struct
	k.privKey, ok = privKey.(ed25519.PrivateKey)
	if !ok {
		return fmt.Errorf("error processing key - is not ec25519")
	}
	// pubfile read
	file, err = os.ReadFile(DEFAULT_KEY_LOCATION + DEFAULT_KEY_NAME + ".pub")
	if err != nil {
		return fmt.Errorf("reading file %s: %w", DEFAULT_KEY_LOCATION+DEFAULT_KEY_NAME+".pub", err)
	}
	// parse file
	block, _ = pem.Decode(file)
	if block == nil {
		return fmt.Errorf("decoding pem from pub file %s", DEFAULT_KEY_LOCATION+DEFAULT_KEY_NAME+".pub")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parsing public key to ec25519: %w", err)
	}
	k.pubKey, ok = pubKey.(ed25519.PublicKey)
	if !ok {
		return fmt.Errorf("error processing pub key - is not ec25519")
	}
	return err
}
