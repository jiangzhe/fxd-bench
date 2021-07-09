package fxd

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/jiangzhe/fxd"
	"github.com/pingcap/tidb/types"
)

var r []int = make([]int, 1024)
var r2 [][2]int8 = make([][2]int8, 1024)

func BenchmarkDecimalFromString(b *testing.B) {
	const n = 1024
	inputBytes := make([][]byte, n)
	inputStrings := make([]string, n)
	output1 := make([]fxd.FixedDecimal, n)
	output2 := make([]types.MyDecimal, n)
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
	b.ResetTimer()
	b.Run("mydecimal-fromString", func(b *testing.B) {
		input := inputBytes
		output2 := output2
		for i := 0; i < b.N; i++ {
			_ = input[n-1]
			_ = output2[n-1]
			for j := 0; j < n; j++ {
				if err := output2[j].FromString(input[j]); err != nil {
					b.Fatal("failed")
				}
			}
		}
	})
}

func BenchmarkDecimalToString(b *testing.B) {
	const n = 1024
	input1 := make([]fxd.FixedDecimal, n)
	input2 := make([]types.MyDecimal, n)
	buf := make([]byte, 0, 128)
	for i := 0; i < n; i++ {
		intg := rand.Int31n(1 << 30)
		frac := rand.Int31n(1 << 20)
		buf = strconv.AppendInt(nil, int64(intg), 10)
		if frac > 0 {
			buf = append(buf, '.')
			buf = strconv.AppendInt(buf, int64(frac), 10)
		}
		s := string(buf)
		input1[i].FromBytesString(buf, true)
		input2[i] = *types.NewDecFromStringForTest(s)
	}
	b.ResetTimer()
	b.Run("fd-ToStringBuffer", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				buf = input1[j].AppendStringBuffer(buf[:0], -1)
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-ToString", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = input1[j].ToString(-1)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-ToString", func(b *testing.B) {
		input2 := input2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			for j := 0; j < n; j++ {
				buf = input2[j].ToString()
			}
		}
	})
}

func BenchmarkDecimalArithPos(b *testing.B) {
	const n = 1024
	input1 := make([][2]fxd.FixedDecimal, n)
	input2 := make([][2]types.MyDecimal, n)
	buf := make([]byte, 0, 128)
	for i := 0; i < n; i++ {
		for j := 0; j < 2; j++ {
			intg := rand.Int31n(1 << 30)
			frac := rand.Int31n(1 << 20)
			buf = strconv.AppendInt(nil, int64(intg), 10)
			if frac > 0 {
				buf = append(buf, '.')
				buf = strconv.AppendInt(buf, int64(frac), 10)
			}
			s := string(buf)
			input1[i][j].FromBytesString(buf, true)
			input2[i][j] = *types.NewDecFromStringForTest(s)
		}
	}
	var fdRes fxd.FixedDecimal
	var mydRes types.MyDecimal

	// addition
	b.ResetTimer()
	b.Run("fd-AddAnyPos", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalAddAny(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-AddPos", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalAdd(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-AddPos", func(b *testing.B) {
		input2 := input2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			for j := 0; j < n; j++ {
				_ = types.DecimalAdd(&input2[j][0], &input2[j][1], &mydRes)
			}
		}
	})

	// subtraction
	b.ResetTimer()
	b.Run("fd-SubAnyPos", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalSubAny(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-SubPos", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalSub(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-SubPos", func(b *testing.B) {
		input2 := input2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			for j := 0; j < n; j++ {
				_ = types.DecimalSub(&input2[j][0], &input2[j][1], &mydRes)
			}
		}
	})

	// multiplication
	b.ResetTimer()
	b.Run("fd-MulAnyPos", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalMulAny(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-MulPos", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalMul(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-MulPos", func(b *testing.B) {
		input2 := input2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			for j := 0; j < n; j++ {
				_ = types.DecimalMul(&input2[j][0], &input2[j][1], &mydRes)
			}
		}
	})

	b.ResetTimer()
	b.Run("fd-DivAnyPos", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalDivAny(&input1[j][0], &input1[j][1], &fdRes, fxd.DivIncrFrac)
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-DivPos", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalDiv(&input1[j][0], &input1[j][1], &fdRes, fxd.DivIncrFrac)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-DivPos", func(b *testing.B) {
		input2 := input2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			for j := 0; j < n; j++ {
				_ = types.DecimalDiv(&input2[j][0], &input2[j][1], &mydRes, 4)
			}
		}
	})
}

func BenchmarkDecimalArithNeg(b *testing.B) {
	const n = 1024
	input1 := make([][2]fxd.FixedDecimal, n)
	input2 := make([][2]types.MyDecimal, n)
	buf := make([]byte, 0, 128)
	for i := 0; i < n; i++ {
		for j := 0; j < 2; j++ {
			intg := rand.Int31n(1 << 30)
			frac := rand.Int31n(1 << 20)
			buf = buf[:0]
			if j == 1 {
				buf = append(buf, '-')
			}
			buf = strconv.AppendInt(buf, int64(intg), 10)
			if frac > 0 {
				buf = append(buf, '.')
				buf = strconv.AppendInt(buf, int64(frac), 10)
			}
			s := string(buf)
			input1[i][j].FromBytesString(buf, true)
			input2[i][j] = *types.NewDecFromStringForTest(s)
		}
	}
	var fdRes fxd.FixedDecimal
	var mydRes types.MyDecimal

	// addition
	b.ResetTimer()
	b.Run("fd-AddAnyNeg", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalAddAny(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-AddNeg", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalAdd(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-AddNeg", func(b *testing.B) {
		input2 := input2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			for j := 0; j < n; j++ {
				_ = types.DecimalAdd(&input2[j][0], &input2[j][1], &mydRes)
			}
		}
	})

	// subtraction
	b.ResetTimer()
	b.Run("fd-SubAnyNeg", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalSubAny(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-SubNeg", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalSub(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-SubNeg", func(b *testing.B) {
		input2 := input2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			for j := 0; j < n; j++ {
				_ = types.DecimalSub(&input2[j][0], &input2[j][1], &mydRes)
			}
		}
	})

	// multiplication
	b.ResetTimer()
	b.Run("fd-MulAnyNeg", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalMulAny(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-MulNeg", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalMul(&input1[j][0], &input1[j][1], &fdRes)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-MulNeg", func(b *testing.B) {
		input2 := input2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			for j := 0; j < n; j++ {
				_ = types.DecimalMul(&input2[j][0], &input2[j][1], &mydRes)
			}
		}
	})

	// division
	b.ResetTimer()
	b.Run("fd-DivAnyNeg", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalDivAny(&input1[j][0], &input1[j][1], &fdRes, fxd.DivIncrFrac)
			}
		}
	})
	b.ResetTimer()
	b.Run("fd-DivNeg", func(b *testing.B) {
		input1 := input1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			for j := 0; j < n; j++ {
				_ = fxd.DecimalDiv(&input1[j][0], &input1[j][1], &fdRes, fxd.DivIncrFrac)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-DivNeg", func(b *testing.B) {
		input2 := input2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			for j := 0; j < n; j++ {
				_ = types.DecimalDiv(&input2[j][0], &input2[j][1], &mydRes, fxd.DivIncrFrac)
			}
		}
	})
}

func BenchmarkDecimalCompare(b *testing.B) {
	const n = 1024
	input1 := make([][2]fxd.FixedDecimal, n)
	input2 := make([][2]types.MyDecimal, n)
	expected := make([]int, n)
	res1 := make([]int, n)
	res2 := make([]int, n)
	buf := make([]byte, 0, 128)
	for i := 0; i < n; i++ {
		for j := 0; j < 2; j++ {
			intg := rand.Int31n(1 << 30)
			frac := rand.Int31n(1 << 20)
			buf = strconv.AppendInt(nil, int64(intg), 10)
			if frac > 0 {
				buf = append(buf, '.')
				buf = strconv.AppendInt(buf, int64(frac), 10)
			}
			s := string(buf)
			input1[i][j].FromBytesString(buf, true)
			input2[i][j] = *types.NewDecFromStringForTest(s)
		}
		expected[i] = input1[i][0].Compare(&input1[i][1])
		myRes := input2[i][0].Compare(&input2[i][1])
		if expected[i] != myRes {
			b.Fatalf("inconsistent result: (%v)compare(%v) fd=%v, myd=%v", input1[i][0], input1[i][1], expected[i], myRes)
		}
	}
	b.ResetTimer()
	b.Run("fd-Compare", func(b *testing.B) {
		input1 := input1
		res1 := res1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			_ = res1[n-1]
			for j := 0; j < n; j++ {
				res1[j] = input1[j][0].Compare(&input1[j][1])
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-Compare", func(b *testing.B) {
		input2 := input2
		res2 := res2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			_ = res2[n-1]
			for j := 0; j < n; j++ {
				res2[j] = input2[j][0].Compare(&input2[j][1])
			}
		}
	})
	for i := 0; i < n; i++ {
		if res1[i] != expected[i] {
			b.Fatalf("result mismatch")
		}
		if res2[i] != expected[i] {
			b.Fatalf("result mismatch")
		}
	}
}

func BenchmarkDecimalRound(b *testing.B) {
	const n = 1024
	input1 := make([]fxd.FixedDecimal, n)
	input2 := make([]types.MyDecimal, n)
	// expected := make([]int, n)
	res1 := make([]fxd.FixedDecimal, n)
	res2 := make([]types.MyDecimal, n)
	buf := make([]byte, 0, 128)
	for i := 0; i < n; i++ {
		intg := rand.Int31n(1 << 30)
		frac := rand.Int31n(1 << 20)
		buf = strconv.AppendInt(nil, int64(intg), 10)
		if frac > 0 {
			buf = append(buf, '.')
			buf = strconv.AppendInt(buf, int64(frac), 10)
		}
		s := string(buf)
		input1[i].FromBytesString(buf, true)
		input2[i] = *types.NewDecFromStringForTest(s)
	}
	b.ResetTimer()
	b.Run("fd-RoundTo", func(b *testing.B) {
		input1 := input1
		res1 := res1
		for i := 0; i < b.N; i++ {
			_ = input1[n-1]
			_ = res1[n-1]
			for j := 0; j < n; j++ {
				input1[j].RoundTo(&res1[j], 4)
			}
		}
	})
	b.ResetTimer()
	b.Run("myd-Round", func(b *testing.B) {
		input2 := input2
		res2 := res2
		for i := 0; i < b.N; i++ {
			_ = input2[n-1]
			_ = res2[n-1]
			for j := 0; j < n; j++ {
				input2[j].Round(&res2[j], 4, types.ModeHalfEven)
			}
		}
	})
}
