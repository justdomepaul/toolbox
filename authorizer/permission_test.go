package authorizer

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type PermissionCodeSuite struct {
	suite.Suite
}

func (suite *PermissionCodeSuite) TestRelationMethod() {
	var authPermission PermissionCode = 1
	suite.Equal([]uint64{2}, authPermission.Code())
	suite.Equal([]uint8{2}, authPermission.Uint8Code())
	suite.Equal([]uint16{2}, authPermission.Uint16Code())
	suite.Equal([]uint32{2}, authPermission.Uint32Code())
	suite.Equal([]uint64{2}, authPermission.Uint64Code())
}

func TestPermissionCodeSuite(t *testing.T) {
	suite.Run(t, new(PermissionCodeSuite))
}

type PermissionSuite struct {
	suite.Suite
}

func (suite *PermissionSuite) TestGeneratePermissionByte() {
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, GeneratePermission[byte](77))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}, GeneratePermission[byte](233))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4}, GeneratePermission[byte](234))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8}, GeneratePermission[byte](235))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16}, GeneratePermission[byte](236))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, GeneratePermission[byte](237))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64}, GeneratePermission[byte](238))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, GeneratePermission[byte](479))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, GeneratePermission[byte](773))
}

func (suite *PermissionSuite) TestGeneratePermissionUint8() {
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, GeneratePermission[uint8](77))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}, GeneratePermission[uint8](233))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4}, GeneratePermission[uint8](234))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8}, GeneratePermission[uint8](235))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16}, GeneratePermission[uint8](236))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, GeneratePermission[uint8](237))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64}, GeneratePermission[uint8](238))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, GeneratePermission[uint8](479))
	suite.Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, GeneratePermission[uint8](773))
}

func (suite *PermissionSuite) TestGeneratePermissionUint16() {
	suite.Equal([]uint16{0, 0, 0, 0, 8192}, GeneratePermission[uint16](77))
	suite.Equal([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 512}, GeneratePermission[uint16](233))
	suite.Equal([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1024}, GeneratePermission[uint16](234))
	suite.Equal([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2048}, GeneratePermission[uint16](235))
	suite.Equal([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4096}, GeneratePermission[uint16](236))
	suite.Equal([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8192}, GeneratePermission[uint16](237))
	suite.Equal([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16384}, GeneratePermission[uint16](238))
	suite.Equal([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32768}, GeneratePermission[uint16](479))
	suite.Equal([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, GeneratePermission[uint16](773))
}

func (suite *PermissionSuite) TestGeneratePermissionUint32() {
	suite.Equal([]uint32{0, 0, 8192}, GeneratePermission[uint32](77))
	suite.Equal([]uint32{0, 0, 0, 0, 0, 0, 0, 512}, GeneratePermission[uint32](233))
	suite.Equal([]uint32{0, 0, 0, 0, 0, 0, 0, 1024}, GeneratePermission[uint32](234))
	suite.Equal([]uint32{0, 0, 0, 0, 0, 0, 0, 2048}, GeneratePermission[uint32](235))
	suite.Equal([]uint32{0, 0, 0, 0, 0, 0, 0, 4096}, GeneratePermission[uint32](236))
	suite.Equal([]uint32{0, 0, 0, 0, 0, 0, 0, 8192}, GeneratePermission[uint32](237))
	suite.Equal([]uint32{0, 0, 0, 0, 0, 0, 0, 16384}, GeneratePermission[uint32](238))
	suite.Equal([]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2147483648}, GeneratePermission[uint32](479))
	suite.Equal([]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, GeneratePermission[uint32](773))
}

func (suite *PermissionSuite) TestGeneratePermissionUint64() {
	suite.Equal([]uint64{0, 8192}, GeneratePermission[uint64](77))
	suite.Equal([]uint64{0, 0, 0, 2199023255552}, GeneratePermission[uint64](233))
	suite.Equal([]uint64{0, 0, 0, 4398046511104}, GeneratePermission[uint64](234))
	suite.Equal([]uint64{0, 0, 0, 8796093022208}, GeneratePermission[uint64](235))
	suite.Equal([]uint64{0, 0, 0, 17592186044416}, GeneratePermission[uint64](236))
	suite.Equal([]uint64{0, 0, 0, 35184372088832}, GeneratePermission[uint64](237))
	suite.Equal([]uint64{0, 0, 0, 70368744177664}, GeneratePermission[uint64](238))
	suite.Equal([]uint64{0, 0, 0, 0, 0, 0, 0, 2147483648}, GeneratePermission[uint64](479))
	suite.Equal([]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, GeneratePermission[uint64](773))
}

func (suite *PermissionSuite) TestSumPermissionUint8() {
	suite.Equal([]uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 34}, SumPermission(GeneratePermission[uint8](0), GeneratePermission[uint8](201), GeneratePermission[uint8](205)))
}

func (suite *PermissionSuite) TestSumPermissionUint16() {
	suite.Equal([]uint16{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8704}, SumPermission(GeneratePermission[uint16](0), GeneratePermission[uint16](201), GeneratePermission[uint16](205)))
}

func (suite *PermissionSuite) TestSumPermissionUint32() {
	suite.Equal([]uint32{1, 0, 0, 0, 0, 0, 8704}, SumPermission(GeneratePermission[uint32](0), GeneratePermission[uint32](201), GeneratePermission[uint32](205)))
}

func (suite *PermissionSuite) TestSumPermissionUint64() {
	suite.Equal([]uint64{1, 0, 0, 8704}, SumPermission(GeneratePermission[uint64](0), GeneratePermission[uint64](201), GeneratePermission[uint64](205)))
}

func (suite *PermissionSuite) TestRemovePermissionUint8() {
	total := SumPermission(GeneratePermission[uint8](0), GeneratePermission[uint8](201), GeneratePermission[uint8](205))
	suite.Equal([]uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 34}, total)
	suite.Equal([]uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, RemovePermission(total, GeneratePermission[uint8](201)))
	suite.Equal([]uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}, RemovePermission(total, GeneratePermission[uint8](205)))
	suite.Equal([]uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 34}, RemovePermission(total, GeneratePermission[uint8](110)))
	suite.Equal([]uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, RemovePermission(total, GeneratePermission[uint8](201), GeneratePermission[uint8](205), GeneratePermission[uint8](110)))
}

func (suite *PermissionSuite) TestRemovePermissionUint16() {
	total := SumPermission(GeneratePermission[uint16](0), GeneratePermission[uint16](201), GeneratePermission[uint16](205))
	suite.Equal([]uint16{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8704}, total)
	suite.Equal([]uint16{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8192}, RemovePermission(total, GeneratePermission[uint16](201)))
	suite.Equal([]uint16{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 512}, RemovePermission(total, GeneratePermission[uint16](205)))
	suite.Equal([]uint16{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8704}, RemovePermission(total, GeneratePermission[uint16](110)))
	suite.Equal([]uint16{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, RemovePermission(total, GeneratePermission[uint16](201), GeneratePermission[uint16](205), GeneratePermission[uint16](110)))
}

func (suite *PermissionSuite) TestRemovePermissionUint32() {
	total := SumPermission(GeneratePermission[uint32](0), GeneratePermission[uint32](201), GeneratePermission[uint32](205))
	suite.Equal([]uint32{1, 0, 0, 0, 0, 0, 8704}, total)
	suite.Equal([]uint32{1, 0, 0, 0, 0, 0, 8192}, RemovePermission(total, GeneratePermission[uint32](201)))
	suite.Equal([]uint32{1, 0, 0, 0, 0, 0, 512}, RemovePermission(total, GeneratePermission[uint32](205)))
	suite.Equal([]uint32{1, 0, 0, 0, 0, 0, 8704}, RemovePermission(total, GeneratePermission[uint32](110)))
	suite.Equal([]uint32{1, 0, 0, 0, 0, 0, 0}, RemovePermission(total, GeneratePermission[uint32](201), GeneratePermission[uint32](205), GeneratePermission[uint32](110)))
}

func (suite *PermissionSuite) TestRemovePermissionUint64() {
	total := SumPermission(GeneratePermission[uint64](0), GeneratePermission[uint64](201), GeneratePermission[uint64](205))
	suite.Equal([]uint64{1, 0, 0, 8704}, total)
	suite.Equal([]uint64{1, 0, 0, 8192}, RemovePermission(total, GeneratePermission[uint64](201)))
	suite.Equal([]uint64{1, 0, 0, 512}, RemovePermission(total, GeneratePermission[uint64](205)))
	suite.Equal([]uint64{1, 0, 0, 8704}, RemovePermission(total, GeneratePermission[uint64](110)))
	suite.Equal([]uint64{1, 0, 0, 0}, RemovePermission(total, GeneratePermission[uint64](201), GeneratePermission[uint64](205), GeneratePermission[uint64](110)))
}

func (suite *PermissionSuite) TestValidPermissionUint8() {
	total := SumPermission(GeneratePermission[uint8](0), GeneratePermission[uint8](201), GeneratePermission[uint8](205))
	suite.Equal([]uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 34}, total)
	suite.True(ValidPermission(total, GeneratePermission[uint8](0)))
	suite.True(ValidPermission(total, GeneratePermission[uint8](201)))
	suite.True(ValidPermission(total, GeneratePermission[uint8](205)))
	suite.False(ValidPermission(total, GeneratePermission[uint8](70)))
	suite.False(ValidPermission(total, GeneratePermission[uint8](110)))
}

func (suite *PermissionSuite) TestValidPermissionUint16() {
	total := SumPermission(GeneratePermission[uint16](0), GeneratePermission[uint16](201), GeneratePermission[uint16](205))
	suite.Equal([]uint16{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8704}, total)
	suite.True(ValidPermission(total, GeneratePermission[uint16](0)))
	suite.True(ValidPermission(total, GeneratePermission[uint16](201)))
	suite.True(ValidPermission(total, GeneratePermission[uint16](205)))
	suite.False(ValidPermission(total, GeneratePermission[uint16](70)))
	suite.False(ValidPermission(total, GeneratePermission[uint16](110)))
}

func (suite *PermissionSuite) TestValidPermissionUint32() {
	total := SumPermission(GeneratePermission[uint32](0), GeneratePermission[uint32](201), GeneratePermission[uint32](205))
	suite.Equal([]uint32{1, 0, 0, 0, 0, 0, 8704}, total)
	suite.True(ValidPermission(total, GeneratePermission[uint32](0)))
	suite.True(ValidPermission(total, GeneratePermission[uint32](201)))
	suite.True(ValidPermission(total, GeneratePermission[uint32](205)))
	suite.False(ValidPermission(total, GeneratePermission[uint32](70)))
	suite.False(ValidPermission(total, GeneratePermission[uint32](110)))
}

func (suite *PermissionSuite) TestValidPermissionUint64() {
	total := SumPermission(GeneratePermission[uint64](0), GeneratePermission[uint64](201), GeneratePermission[uint64](205))
	suite.Equal([]uint64{1, 0, 0, 8704}, total)
	suite.True(ValidPermission(total, GeneratePermission[uint64](0)))
	suite.True(ValidPermission(total, GeneratePermission[uint64](201)))
	suite.True(ValidPermission(total, GeneratePermission[uint64](205)))
	suite.False(ValidPermission(total, GeneratePermission[uint64](70)))
	suite.False(ValidPermission(total, GeneratePermission[uint64](110)))
}

func TestPermissionSuite(t *testing.T) {
	suite.Run(t, new(PermissionSuite))
}

func BenchmarkGeneratePermissionUint8(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GeneratePermission[uint8](4096)
	}
	b.StopTimer()
}

func BenchmarkGeneratePermissionUint16(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GeneratePermission[uint8](4096)
	}
	b.StopTimer()
}

func BenchmarkGeneratePermissionUint32(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GeneratePermission[uint8](4096)
	}
	b.StopTimer()
}

func BenchmarkGeneratePermissionUint64(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GeneratePermission[uint8](4096)
	}
	b.StopTimer()
}

func BenchmarkSumPermissionUint8(b *testing.B) {
	var result [][]byte
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[byte](i))
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SumPermission(result...)
	}
	b.StopTimer()
}

func BenchmarkSumPermissionUint16(b *testing.B) {
	var result [][]uint16
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint16](i))
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SumPermission(result...)
	}
	b.StopTimer()
}

func BenchmarkSumPermissionUint32(b *testing.B) {
	var result [][]uint32
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint32](i))
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SumPermission(result...)
	}
	b.StopTimer()
}

func BenchmarkSumPermissionUint64(b *testing.B) {
	var result [][]uint64
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint64](i))
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SumPermission(result...)
	}
	b.StopTimer()
}

func BenchmarkRemovePermissionUint8(b *testing.B) {
	var result [][]byte
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint8](i))
	}
	total := SumPermission(result...)
	r3 := GeneratePermission[uint8](3)
	r77 := GeneratePermission[uint8](77)
	r1058 := GeneratePermission[uint8](1058)
	r2048 := GeneratePermission[uint8](2048)
	r3077 := GeneratePermission[uint8](3077)
	r4066 := GeneratePermission[uint8](4066)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		RemovePermission(total, r3, r77, r1058, r2048, r3077, r4066)
	}
	b.StopTimer()
}

func BenchmarkRemovePermissionUint16(b *testing.B) {
	var result [][]uint16
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint16](i))
	}
	total := SumPermission(result...)
	r3 := GeneratePermission[uint16](3)
	r77 := GeneratePermission[uint16](77)
	r1058 := GeneratePermission[uint16](1058)
	r2048 := GeneratePermission[uint16](2048)
	r3077 := GeneratePermission[uint16](3077)
	r4066 := GeneratePermission[uint16](4066)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		RemovePermission(total, r3, r77, r1058, r2048, r3077, r4066)
	}
	b.StopTimer()
}

func BenchmarkRemovePermissionUint32(b *testing.B) {
	var result [][]uint32
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint32](i))
	}
	total := SumPermission(result...)
	r3 := GeneratePermission[uint32](3)
	r77 := GeneratePermission[uint32](77)
	r1058 := GeneratePermission[uint32](1058)
	r2048 := GeneratePermission[uint32](2048)
	r3077 := GeneratePermission[uint32](3077)
	r4066 := GeneratePermission[uint32](4066)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		RemovePermission(total, r3, r77, r1058, r2048, r3077, r4066)
	}
	b.StopTimer()
}

func BenchmarkRemovePermissionUint64(b *testing.B) {
	var result [][]uint64
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint64](i))
	}
	total := SumPermission(result...)
	r3 := GeneratePermission[uint64](3)
	r77 := GeneratePermission[uint64](77)
	r1058 := GeneratePermission[uint64](1058)
	r2048 := GeneratePermission[uint64](2048)
	r3077 := GeneratePermission[uint64](3077)
	r4066 := GeneratePermission[uint64](4066)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		RemovePermission(total, r3, r77, r1058, r2048, r3077, r4066)
	}
	b.StopTimer()
}

func BenchmarkValidPermissionUint8(b *testing.B) {
	var result [][]byte
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint8](i))
	}
	total := SumPermission(result...)
	valid3056 := GeneratePermission[uint8](3056)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ValidPermission(total, valid3056)
	}
	b.StopTimer()
}

func BenchmarkValidPermissionUint16(b *testing.B) {
	var result [][]uint16
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint16](i))
	}
	total := SumPermission(result...)
	valid3056 := GeneratePermission[uint16](3056)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ValidPermission(total, valid3056)
	}
	b.StopTimer()
}

func BenchmarkValidPermissionUint32(b *testing.B) {
	var result [][]uint32
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint32](i))
	}
	total := SumPermission(result...)
	valid3056 := GeneratePermission[uint32](3056)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ValidPermission(total, valid3056)
	}
	b.StopTimer()
}

func BenchmarkValidPermissionUint64(b *testing.B) {
	var result [][]uint64
	for i := 1; i < 4096; i++ {
		result = append(result, GeneratePermission[uint64](i))
	}
	total := SumPermission(result...)
	valid3066 := GeneratePermission[uint64](3066)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ValidPermission(total, valid3066)
	}
	b.StopTimer()
}
