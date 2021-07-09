package fxd

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/jiangzhe/fxd"
)

func TestDecimalRandString(t *testing.T) {
	const n = 1024
	inputStrings := make([]string, n)
	buf := make([]byte, 0, 128)
	for i := 0; i < n; i++ {
		intg := rand.Int31n(1 << 30)
		frac := rand.Int31n(1 << 20)
		buf = strconv.AppendInt(nil, int64(intg), 10)
		if frac > 0 {
			buf = append(buf, '.')
			buf = strconv.AppendInt(buf, int64(frac), 10)
		}
		inputStrings[i] = string(buf)
		fmt.Printf("s=%v\n", inputStrings[i])
		buf = buf[:0]
	}
}

func BenchmarkDecimalFromInfString(b *testing.B) {
	const n = 1024
	inputBytes := make([][]byte, n)
	inputStrings := make([]string, n)
	output1 := make([]fxd.FixedDecimal, n)
	buf := make([]byte, 0, 128)
	for i := 0; i < n; i++ {
		intg := rand.Int31n(1 << 30)
		frac := rand.Int31n(1 << 20)
		buf = strconv.AppendInt(nil, int64(intg), 10)
		if frac > 0 {
			buf = append(buf, '.')
			buf = strconv.AppendInt(buf, int64(frac), 10)
		}
		inputStrings[i] = string(buf)
		inputStrings[i] = "Infinity"
		inputBytes[i] = []byte(inputStrings[i])
		buf = buf[:0]
	}

	b.ResetTimer()
	b.Run("fd-FromString", func(b *testing.B) {
		input := inputStrings
		output1 := output1
		for i := 0; i < b.N; i++ {
			_ = input[n-1]
			_ = output1[n-1]
			for j := 0; j < n; j++ {
				if err := output1[j].FromAsciiString(input[j], true); err != nil {
					b.Fatal("failed")
				}
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-FromBytesString", func(b *testing.B) {
		input := inputBytes
		output1 := output1
		for i := 0; i < b.N; i++ {
			_ = input[n-1]
			_ = output1[n-1]
			for j := 0; j < n; j++ {
				if err := output1[j].FromBytesString(input[j], true); err != nil {
					b.Fatal("failed")
				}
			}
		}
	})
}
