package option

import (
	"encoding/base64"
	"fmt"
)

//Crc represents crc value
type Crc struct {
	Hash uint32
}

func (c *Crc) Encode() string {
	b := []byte{byte(c.Hash >> 24), byte(c.Hash >> 16), byte(c.Hash >> 8), byte(c.Hash)}
	return base64.StdEncoding.EncodeToString(b)
}

func (c *Crc) Decode(encoded string) error {
	d, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	if len(d) != 4 {
		return fmt.Errorf("storage: %q does not encode a 32-bit value", d)
	}
	c.Hash = uint32(d[0])<<24 + uint32(d[1])<<16 + uint32(d[2])<<8 + uint32(d[3])
	return nil
}
