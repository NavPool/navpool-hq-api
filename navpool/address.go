package navpool

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
)

type PoolAddress struct {
	SpendingAddress    string `json:"spendingAddress"`
	StakingAddress     string `json:"stakingAddress"`
	ColdStakingAddress string `json:"coldStakingAddress"`
}

func (e *PoolApi) GetPoolAddress(spendingAddress string, signature string) (address PoolAddress, err error) {
	method := fmt.Sprintf("/address/%s/add/%s", spendingAddress, b64.StdEncoding.EncodeToString([]byte(signature)))

	response, err := e.client.call(method, "GET", nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &address)
	return
}
