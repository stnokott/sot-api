package structs

import "fmt"

// Balance contains the gold, doubloon and ancient coin amount
type Balance struct {
	Gold         int
	Doubloons    int
	AncientCoins int
}

func (b Balance) String() string {
	return fmt.Sprintf("Balance{Gold=%d, Doubloons=%d, AncientCoins=%d}", b.Gold, b.Doubloons, b.AncientCoins)
}
