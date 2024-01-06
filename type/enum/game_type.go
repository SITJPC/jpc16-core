package enum

import (
	"encoding/json"
	"fmt"
)

type GameType string

var (
	GameTypeAudit  GameType = "audit"
	GameTypeCredit GameType = "credit"
)

func (r *GameType) UnmarshalJSON(data []byte) error {
	v := new(string)
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	val := GameType(*v)
	if val != GameTypeAudit && val != GameTypeCredit {
		return fmt.Errorf("%s is not a valid GameType value", val)
	}

	*r = val
	return nil
}
