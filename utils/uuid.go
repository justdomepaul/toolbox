package utils

import (
	"github.com/google/uuid"
	"github.com/justdomepaul/toolbox/errorhandler"
)

func ParseUUID(input string) []byte {
	if input != "" {
		uid, err := uuid.Parse(input)
		if err != nil {
			panic(errorhandler.NewErrExecute(err))
		}
		return uid[:]
	}
	return []byte("")
}

func FromUUID(input []byte) string {
	if len(input) > 0 {
		uid, err := uuid.FromBytes(input)
		if err != nil {
			panic(errorhandler.NewErrExecute(err))
		}
		return uid.String()
	}
	return ""
}
