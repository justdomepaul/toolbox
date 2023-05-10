package authorizer

type UintSeq interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}
