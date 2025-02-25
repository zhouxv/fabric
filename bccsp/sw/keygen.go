/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sw

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"crypto/pqc"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/open-quantum-safe/liboqs-go/oqs"
)

type ecdsaKeyGenerator struct {
	curve elliptic.Curve
}

func (kg *ecdsaKeyGenerator) KeyGen(opts bccsp.KeyGenOpts) (bccsp.Key, error) {
	privKey, err := ecdsa.GenerateKey(kg.curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("Failed generating ECDSA key for [%v]: [%s]", kg.curve, err)
	}

	return &ecdsaPrivateKey{privKey}, nil
}

type aesKeyGenerator struct {
	length int
}

func (kg *aesKeyGenerator) KeyGen(opts bccsp.KeyGenOpts) (bccsp.Key, error) {
	lowLevelKey, err := GetRandomBytes(int(kg.length))
	if err != nil {
		return nil, fmt.Errorf("Failed generating AES %d key [%s]", kg.length, err)
	}

	return &aesPrivateKey{lowLevelKey, false}, nil
}

type pqcKeyGenerator struct{}

func (kg *pqcKeyGenerator) KeyGen(opts bccsp.KeyGenOpts) (bccsp.Key, error) {
	alg := opts.Algorithm()
	if alg == "" {
		alg = "Dilithium3"
	}

	sig := pqc.SigType(alg)

	//初始化
	signer := oqs.Signature{}
	// defer signer.Clean()
	if err := signer.Init(alg, nil); err != nil {
		return nil, fmt.Errorf("Failed PQC init %s key [%s]", alg, err)
	}

	//生成密钥对
	pubKey, err := signer.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("Failed generating PQC KeyPair %s key [%s]", alg, err)
	}
	privKey := signer.ExportSecretKey()

	lowLevelKey := &pqc.SecretKey{}
	lowLevelKey.Sk = privKey
	lowLevelKey.PublicKey = *&pqc.PublicKey{Pk: pubKey, Sig: *&pqc.OQSSigInfo{Algorithm: sig}}

	return &pqcPrivateKey{lowLevelKey}, nil
}
