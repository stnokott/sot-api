package structs

import "fmt"

// Profile contains the gold, doubloon and ancient coin amount
type Profile struct {
	Gold            int
	Doubloons       int
	AncientCoins    int
	ProfileImageURL string `json:"image"`
	Title           string
}

func (p Profile) String() string {
	return fmt.Sprintf("Balance{Gold=%d, Doubloons=%d, AncientCoins=%d}", p.Gold, p.Doubloons, p.AncientCoins)
}
