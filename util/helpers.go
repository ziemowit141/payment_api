package util

import uuid "github.com/nu7hatch/gouuid"

func GenerateUniqeId() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return id.String()
}
