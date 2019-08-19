package option

import "encoding/base64"

//Md5 represents md5 value
type Md5 struct {
	Hash []byte
}

//Encode encode base64 hash value
func (m *Md5) Encode() string {
	return base64.StdEncoding.EncodeToString(m.Hash)
}

//Decode base64 decode
func (m *Md5) Decode(encoded string) (err error) {
	m.Hash, err = base64.StdEncoding.DecodeString(encoded)
	return err
}
