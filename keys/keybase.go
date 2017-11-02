package keys

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/tendermint/go-crypto/nano"
)

// XXX Lets use go-crypto/bcrypt and ascii encoding directly in here without
// further wrappers around a store or DB.
// Copy functions from: https://github.com/tendermint/mintkey/blob/master/cmd/mintkey/common.go
//
// dbKeybase combines encyption and storage implementation to provide
// a full-featured key manager
type dbKeybase struct {
	db    dbm.DB
	codec Codec
}

func New(db dbm.DB, codec Codec) dbKeybase {
	return dbKeybase{
		db:    db,
		codec: codec,
	}
}

var _ Keybase = dbKeybase{}

// Create adds a new key to the storage engine, returning error if
// another key already stored under this name
//
// algo must be a supported go-crypto algorithm: ed25519, secp256k1
func (kb dbKeybase) Create(name, passphrase, algo string) (Info, string, error) {
	// 128-bits are the all the randomness we can make use of
	secret := crypto.CRandBytes(16)
	key, err := generate(algo, secret)
	if err != nil {
		return Info{}, "", err
	}

	err = kb.es.Put(name, passphrase, key)
	if err != nil {
		return Info{}, "", err
	}

	// we append the type byte to the serialized secret to help with recovery
	// ie [secret] = [secret] + [type]
	typ := key.Bytes()[0]
	secret = append(secret, typ)

	seed, err := kb.codec.BytesToWords(secret)
	phrase := strings.Join(seed, " ")
	return info(name, key), phrase, err
}

// Recover takes a seed phrase and tries to recover the private key.
//
// If the seed phrase is valid, it will create the private key and store
// it under name, protected by passphrase.
//
// Result similar to New(), except it doesn't return the seed again...
func (kb dbKeybase) Recover(name, passphrase, seedphrase string) (Info, error) {
	words := strings.Split(strings.TrimSpace(seedphrase), " ")
	secret, err := kb.codec.WordsToBytes(words)
	if err != nil {
		return Info{}, err
	}

	// secret is comprised of the actual secret with the type appended
	// ie [secret] = [secret] + [type]
	l := len(secret)
	secret, typ := secret[:l-1], secret[l-1]

	key, err := generateByType(typ, secret)
	if err != nil {
		return Info{}, err
	}

	// d00d, it worked!  create the bugger....
	err = kb.es.Put(name, passphrase, key)
	return info(name, key), err
}

// List loads the keys from the storage and enforces alphabetical order
func (kb dbKeybase) List() ([]Info, error) {
	res, err := kb.es.List()
	sort.Slice(res, func(a, b int) bool { return res[a].Name < res[b].Name })
	return res, err
}

// Get returns the public information about one key
func (kb dbKeybase) Get(name string) (Info, error) {
	_, info, err := kb.es.store.Get(name)
	return info, err
}

// Sign will modify the Signable in order to attach a valid signature with
// this public key
//
// If no key for this name, or the passphrase doesn't match, returns an error
func (kb dbKeybase) Sign(name, passphrase string, msg []byte) (sig crypto.Signature, pk crypto.PubKey, err error) {
	var key crypto.PrivKey
	key, _, err = kb.es.Get(name, passphrase)
	if err != nil {
		return
	}
	sig = key.Sign(msg)
	pk = key.PubKey()
	return
}

// Export decodes the private key with the current password, encodes
// it with a secure one-time password and generates a sequence that can be
// Imported by another dbKeybase
//
// This is designed to copy from one device to another, or provide backups
// during version updates.
func (kb dbKeybase) Export(name, oldpass, transferpass string) ([]byte, error) {
	key, _, err := kb.es.Get(name, oldpass)
	if err != nil {
		return nil, err
	}

	res, err := kb.es.coder.Encrypt(key, transferpass)
	return res, err
}

// Import accepts bytes generated by Export along with the same transferpass
// If they are valid, it stores the password under the given name with the
// new passphrase.
func (kb dbKeybase) Import(name, newpass, transferpass string, data []byte) error {
	key, err := kb.es.coder.Decrypt(data, transferpass)
	if err != nil {
		return err
	}

	return kb.es.Put(name, newpass, key)
}

// Delete removes key forever, but we must present the
// proper passphrase before deleting it (for security)
func (kb dbKeybase) Delete(name, passphrase string) error {
	// verify we have the proper password before deleting
	_, _, err := kb.es.Get(name, passphrase)
	if err != nil {
		return err
	}
	return kb.es.Delete(name)
}

// Update changes the passphrase with which a already stored key is encoded.
//
// oldpass must be the current passphrase used for encoding, newpass will be
// the only valid passphrase from this time forward
func (kb dbKeybase) Update(name, oldpass, newpass string) error {
	key, _, err := kb.es.Get(name, oldpass)
	if err != nil {
		return err
	}

	// we must delete first, as Putting over an existing name returns an error
	kb.Delete(name, oldpass)

	return kb.es.Put(name, newpass, key)
}

func generate(algo string, secret []byte) (crypto.PrivKey, error) {
	switch algo {
	case crypto.NameEd25519:
		return crypto.GenPrivKeyEd25519FromSecret(secret).Wrap(), nil
	case crypto.NameSecp256k1:
		return crypto.GenPrivKeySecp256k1FromSecret(secret).Wrap(), nil
	case nano.NameLedgerEd25519:
		return nano.NewPrivKeyLedgerEd25519Ed25519()
	default:
		err := errors.Errorf("Cannot generate keys for algorithm: %s", algo)
		return crypto.PrivKey{}, err
	}
}

func generateByType(typ byte, secret []byte) (crypto.PrivKey, error) {
	switch typ {
	case crypto.TypeEd25519:
		return generate(crypto.NameEd25519, secret)
	case crypto.TypeSecp256k1:
		return generate(crypto.NameSecp256k1, secret)
	case nano.TypeLedgerEd25519:
		return generate(nano.NameLedgerEd25519, secret)
	default:
		err := errors.Errorf("Cannot generate keys for algorithm: %X", typ)
		return crypto.PrivKey{}, err
	}
}

func encrypt(key crypto.PrivKey, pass string) ([]byte, error) {
	if pass == "" {
		return key.Bytes(), nil
	}
	s := secret(pass)
	cipher := crypto.EncryptSymmetric(key.Bytes(), s)
	return cipher, nil
}

func decrypt(data []byte, pass string) (key crypto.PrivKey, err error) {
	private := data
	if pass != "" {
		s := secret(pass)
		private, err = crypto.DecryptSymmetric(data, s)
		if err != nil {
			return crypto.PrivKey{}, errors.Wrap(err, "Invalid Passphrase")
		}
	}
	key, err = crypto.PrivKeyFromBytes(private)
	return key, errors.Wrap(err, "Invalid Passphrase")
}
