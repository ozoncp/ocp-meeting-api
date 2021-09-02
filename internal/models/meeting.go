package models

import (
	"encoding/json"
	"time"
)

type Meeting struct {
	Id        uint64
	UserId    uint64
	Link      string
	Start     time.Time
	End       time.Time
	IsDeleted bool
}

func (m Meeting) String() string {
	result, _ := json.MarshalIndent(&m, "", "\t")
	return string(result)
}
