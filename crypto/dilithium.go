package crypto

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/theQRL/go-libp2p-qrl/protos"
	"github.com/theQRL/go-qrllib-crypto/dilithium"
	"io"
)

type DilithiumPrivateKey struct {
	pb *protos.DilithiumKeys
}

type DilithiumPublicKey struct {
	pb *protos.DilithiumPublicKey
}

func GenerateDilithiumKey(src io.Reader) (crypto.PrivKey, crypto.PubKey, error) {
	if !isLoaded {
		return nil, nil, errors.New("LoadAllExtendedKeyTypes before using")
	}
	d := dilithium.NewDilithium()

	return &DilithiumPrivateKey{
			pb: &protos.DilithiumKeys{
				Sk: d.SK(),
				Pk: d.PK(),
			},
		},
		&DilithiumPublicKey{
			pb: &protos.DilithiumPublicKey{
				Pk: d.PK(),
			},
		},
		nil
}

func (sk *DilithiumPrivateKey) Type() pb.KeyType {
	return mapExtendedKeyTypes[LATTICE]
}

func (sk *DilithiumPrivateKey) Bytes() ([]byte, error) {
	return crypto.MarshalPrivateKey(sk)
}

func (sk *DilithiumPrivateKey) Raw() ([]byte, error) {
	return proto.Marshal(sk.pb)
}

func (sk *DilithiumPrivateKey) Equals(k crypto.Key) bool {
	return basicEquals(sk, k)
}

func (sk *DilithiumPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	hash := h.Sum(nil)
	d := dilithium.RecoverDilithium(sk.pb.Pk, sk.pb.Sk)
	signature := d.Sign(hash)

	return signature, nil
}

func (sk *DilithiumPrivateKey) GetPublic() crypto.PubKey {
	return &DilithiumPublicKey{
		pb: &protos.DilithiumPublicKey{
			Pk: sk.pb.Pk,
		},
	}
}

func (pk *DilithiumPublicKey) Bytes() ([]byte, error) {
	return crypto.MarshalPublicKey(pk)
}

func (pk *DilithiumPublicKey) Type() pb.KeyType {
	return mapExtendedKeyTypes[LATTICE]
}

func (pk *DilithiumPublicKey) Raw() ([]byte, error) {
	return proto.Marshal(pk.pb)
}

func (pk *DilithiumPublicKey) Equals(k crypto.Key) bool {
	return basicEquals(pk, k)
}

func (pk *DilithiumPublicKey) Verify(data, sigBytes []byte) (bool, error) {
	h := sha256.New()
	h.Write(data)
	hashedMessage := h.Sum(nil)
	result := dilithium.DilithiumVerify(sigBytes, pk.pb.Pk, hashedMessage)
	return result, nil
}

func basicEquals(k1, k2 crypto.Key) bool {
	if k1.Type() != k2.Type() {
		return false
	}

	a, err := k1.Raw()
	if err != nil {
		return false
	}
	b, err := k2.Raw()
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(a, b) == 1
}

func UnmarshalDilithiumPublicKey(b []byte) (crypto.PubKey, error) {
	pbData := &protos.DilithiumPublicKey{}
	err := proto.Unmarshal(b, pbData)
	if err != nil {
		return nil, err
	}
	d := &DilithiumPublicKey{
		pb: pbData,
	}
	return d, err
}

func UnmarshalDilithiumPrivateKey(b []byte) (crypto.PrivKey, error) {
	pbData := &protos.DilithiumKeys{}
	err := proto.Unmarshal(b, pbData)
	if err != nil {
		return nil, err
	}
	d := &DilithiumPrivateKey{
		pb: pbData,
	}
	return d, err
}
