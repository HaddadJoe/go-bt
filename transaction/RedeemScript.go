package transaction

import (
	"crypto/sha256"
	"errors"
	"log"

	"bitbucket.org/simon_ordish/cryptolib"
	"github.com/btcsuite/btcutil/base58"

	"golang.org/x/crypto/ripemd160"
)

const (
	opZERO          = 0x00
	opBASE          = 0x50
	opHASH160       = 0xa9
	opCHECKMULTISIG = 0xae
	opEQUAL         = 0x87
)

// RedeemScript type
type RedeemScript struct {
	SignaturesRequired int
	PublicKeys         [][]byte
	Signatures         [][]byte
}

// NewRedeemScript comment
func NewRedeemScript(signaturesRequired int) (*RedeemScript, error) {
	if signaturesRequired < 2 {
		return nil, errors.New("Must have 2 or more required signatures for multisig")
	}

	if signaturesRequired > 16 {
		return nil, errors.New("More than 16 signatures is not supported")
	}

	rs := &RedeemScript{
		SignaturesRequired: signaturesRequired,
	}

	return rs, nil
}

func hash160(data []byte) []byte {
	sha := sha256.New()
	ripe := ripemd160.New()
	sha.Write(data)
	ripe.Write(sha.Sum(nil))
	return ripe.Sum(nil)
}

// AddPublicKey comment
func (rs *RedeemScript) AddPublicKey(pkey string) error {

	pk, err := cryptolib.NewPublicKey(pkey)
	if err != nil {
		return err
	}

	result, err := pk.Child(0)
	if err != nil {
		return err
	}

	result2, err := result.Child(0)
	if err != nil {
		return err
	}

	log.Println(result2.PublicKeyStr)

	rs.PublicKeys = append(rs.PublicKeys, result2.PublicKey)

	// b, _, err := base58.CheckDecode(pkey)
	// if err != nil {
	// 	return err
	// }

	// rs.PublicKeys = append(rs.PublicKeys, b)
	return nil
}

func (rs *RedeemScript) getAddress() string {
	script := rs.getRedeemScript()
	hash := hash160(script)
	// hash = append([]byte{0x05}, hash...)
	return base58.CheckEncode(hash, 0x05)
}

func (rs *RedeemScript) getPublicKeys() [][]byte {
	return rs.PublicKeys
}

func (rs *RedeemScript) getRedeemScript() []byte {
	var b []byte

	b = append(b, byte(opBASE+rs.SignaturesRequired))

	for _, pk := range rs.PublicKeys {
		b = append(b, byte(len(pk)))
		b = append(b, pk...)
	}

	b = append(b, byte(len(rs.PublicKeys)))

	b = append(b, opCHECKMULTISIG)

	return b
}

func (rs *RedeemScript) getRedeemScriptHash() []byte {
	return hash160(rs.getRedeemScript())
}

func (rs *RedeemScript) getScriptPubKey() []byte {
	var b []byte

	h := rs.getRedeemScriptHash()

	b = append(b, opHASH160)
	b = append(b, byte(len(h)))
	b = append(b, h...)
	b = append(b, opEQUAL)

	return b
}