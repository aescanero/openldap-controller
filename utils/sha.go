/*Copyright [2023] [Alejandro Escanero Blanco <aescanero@disasterproject.com>]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.*/

package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"math/rand"
)

type Encode struct{}

func (enc Encode) MakeSSHAEncode(pass []byte) (passEncode []byte) {
	salt := make([]byte, 4)
	rand.Read(salt)
	sha := sha1.New()
	sha.Write(pass)
	sha.Write(salt)
	h := sha.Sum(nil)
	return append(h, salt...)
}

func (enc Encode) MakeSSHA256Encode(pass []byte) (passEncode []byte) {
	salt := make([]byte, 4)
	rand.Read(salt)
	sha := sha256.New()
	sha.Write(pass)
	sha.Write(salt)
	h := sha.Sum(nil)
	return append(h, salt...)
}
