package uuidx

import (
	"fmt"

	"github.com/google/uuid"
)

func Parse(uuids []string) (uuid.UUIDs, error) {
	if uuids == nil {
		return nil, fmt.Errorf("uuids are nil")
	}

	out := make(uuid.UUIDs, len(uuids))

	for i, v := range uuids {
		uuid, err := uuid.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("parse uuid: %w", err)
		}

		out[i] = uuid
	}

	return out, nil
}

func MustParse(uuids []string) uuid.UUIDs {
	out := make(uuid.UUIDs, len(uuids))

	for i, v := range uuids {
		uuid := uuid.MustParse(v)

		out[i] = uuid
	}

	return out
}
