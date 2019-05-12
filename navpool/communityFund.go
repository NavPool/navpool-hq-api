package navpool

import (
	"fmt"
)

func (e *PoolApi) ProposalVote(spendingAddress string, hash string, vote string) (err error) {
	method := "/community-fund/proposal/vote"

	data := fmt.Sprintf(`{"spending_address":"%s", "hash":"%s", "vote":"%s"}`, spendingAddress, hash, vote)

	_, err = e.client.call(method, "POST", data)

	return
}

func (e *PoolApi) PaymentRequestVote(spendingAddress string, hash string, vote string) (err error) {
	method := "/community-fund/payment-request/vote"

	data := fmt.Sprintf(`{"spending_address":"%s", "hash":"%s", "vote":"%s"}`, spendingAddress, hash, vote)

	_, err = e.client.call(method, "POST", data)

	return
}
