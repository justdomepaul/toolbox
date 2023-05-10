package authorizer

import "fmt"

type PermissionCode int

func (p PermissionCode) Code() []uint64 {
	return GeneratePermission[uint64](int(p))
}

func (p PermissionCode) Uint8Code() []uint8 {
	return GeneratePermission[uint8](int(p))
}

func (p PermissionCode) Uint16Code() []uint16 {
	return GeneratePermission[uint16](int(p))
}

func (p PermissionCode) Uint32Code() []uint32 {
	return GeneratePermission[uint32](int(p))
}

func (p PermissionCode) Uint64Code() []uint64 {
	return GeneratePermission[uint64](int(p))
}

func GeneratePermission[T UintSeq](bitIdx int) []T {
	var divisor int
	switch fmt.Sprintf("%T", *new(T)) {
	case "uint8":
		divisor = 8
	case "uint16":
		divisor = 16
	case "uint32":
		divisor = 32
	case "uint64":
		divisor = 64
	}
	size := bitIdx/divisor + 1
	offset := bitIdx % divisor
	p := make([]T, size)
	p[size-1] = T(1) << offset
	return p
}

func SumPermission[T UintSeq](permissions ...[]T) []T {
	max := 0
	for _, p := range permissions {
		l := len(p)
		if max < l {
			max = l
		}
	}
	merged := make([]T, max)
	for _, p := range permissions {
		for i, b := range p {
			merged[i] |= b
		}
	}
	return merged
}

func RemovePermission[T UintSeq](origin []T, permissions ...[]T) []T {
	sum := make([]T, len(origin))
	copy(sum, origin)
	for _, p := range permissions {
		maxSize := len(sum)
		minSize := len(p)
		var (
			size, diff int
		)
		if maxSize >= minSize {
			size = maxSize
			diff = maxSize - minSize
		} else {
			continue
		}
		for i := 0; i < size; i++ {
			if i < diff {
				continue
			}
			sum[i] = sum[i] & ^p[i-diff]
		}
	}
	return sum
}

func ValidPermission[T UintSeq](origin, permission []T) bool {
	for _, o := range origin {
		for _, b := range permission {
			if o&b > 0 {
				return true
			}
		}
	}
	return false
}
