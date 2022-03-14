package shorten

import (
	"crypto/md5"
	"fmt"
	"github.com/justdomepaul/toolbox/base58"
	"math/big"
	"math/bits"
	"strings"
	"time"
)

var (
	now = time.Now
)

const (
	beginTimestampUnixMicro = int64(1633046400000000) // 2021-10-01T00:00:00Z
)

func MaxTimestamp(bitSize int, unit time.Duration) *time.Time {
	max := (1<<bitSize - 1) * unit
	t := time.UnixMicro(beginTimestampUnixMicro).Add(max)
	return &t
}

// Timestamp42 generates 42 bits strings in 7~8 char encoded by Base58.
// The resolution is 100 microseconds and max timestamp is 2035-09-08T07:57:31.1103Z.
func Timestamp42() string {
	us := now().UnixMicro() - beginTimestampUnixMicro
	normalized := us / 100 // precision 100 microseconds
	return timestamp(40, uint64(normalized))
}

func timestamp(bitSize int, normalized uint64) string {
	b64 := bits.Reverse64(normalized)
	bx := b64 >> (64 - bitSize)
	bbx := big.NewInt(int64(bx))
	return base58.Encode(bbx.Bytes())
}

func Timestamp() string {
	reverse := bits.Reverse64(uint64(now().UnixNano()))
	return base58.Encode(big.NewInt(int64(reverse)).Bytes())
}

func String(inputs ...string) string {
	hashBytes := md5f(strings.Join(inputs, ""))
	generatedNumber := new(big.Int).SetBytes(hashBytes).Uint64()
	finalString := base58.Encode([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString
}

func md5f(input string) []byte {
	algorithm := md5.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}
