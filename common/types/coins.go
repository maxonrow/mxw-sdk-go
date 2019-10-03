package types

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
)

type Coin struct {
	Denom string `json:"denom"`

	// To allow the use of unsigned integers (see: #1273) a larger refactor will
	// need to be made. So we use signed integers for now with safety measures in
	// place preventing negative values being used.
	Amount Int `json:"amount"`
}

// NewCoin returns a new coin with a denomination and amount. It will panic if
// the amount is negative.
func NewCoin(denom string, amount Int) Coin {
	if err := validate(denom, amount); err != nil {
		panic(err)
	}

	return Coin{
		Denom:  denom,
		Amount: amount,
	}
}

// Int wraps integer with 256 bit range bound
// Checks overflow, underflow and division by zero
// Exists in range from -(2^maxBitLen-1) to 2^maxBitLen-1
type Int struct {
	i *big.Int
}

// validate returns an error if the Coin has a negative amount or if
// the denom is invalid.
func validate(denom string, amount Int) error {
	if err := validateDenom(denom); err != nil {
		return err
	}

	if amount.LT(ZeroInt()) {
		return fmt.Errorf("negative coin amount: %v", amount)
	}

	return nil
}

func (coin Coin) IsZero() bool {
	return coin.Amount == 0
}

func (coin Coin) IsPositive() bool {
	return coin.Amount > 0
}

func (coin Coin) IsNotNegative() bool {
	return coin.Amount >= 0
}

func (coin Coin) SameDenomAs(other Coin) bool {
	return (coin.Denom == other.Denom)
}

func (coin Coin) Plus(coinB Coin) Coin {
	if !coin.SameDenomAs(coinB) {
		return coin
	}
	return Coin{coin.Denom, coin.Amount + coinB.Amount}
}

// Coins def
type Coins []Coin

func (coins Coins) IsValid() bool {
	switch len(coins) {
	case 0:
		return true
	case 1:
		return !coins[0].IsZero()
	default:
		lowDenom := coins[0].Denom
		for _, coin := range coins[1:] {
			if coin.Denom <= lowDenom {
				return false
			}
			if coin.IsZero() {
				return false
			}
			lowDenom = coin.Denom
		}
		return true
	}
}

func (coins Coins) IsPositive() bool {
	if len(coins) == 0 {
		return false
	}
	for _, coin := range coins {
		if !coin.IsPositive() {
			return false
		}
	}
	return true
}

func (coins Coins) Plus(coinsB Coins) Coins {
	sum := ([]Coin)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(coins), len(coinsB)
	for {
		if indexA == lenA {
			if indexB == lenB {
				return sum
			}
			return append(sum, coinsB[indexB:]...)
		} else if indexB == lenB {
			return append(sum, coins[indexA:]...)
		}
		coinA, coinB := coins[indexA], coinsB[indexB]
		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1:
			sum = append(sum, coinA)
			indexA++
		case 0:
			if coinA.Amount+coinB.Amount == 0 {
				// ignore 0 sum coin type
			} else {
				sum = append(sum, coinA.Plus(coinB))
			}
			indexA++
			indexB++
		case 1:
			sum = append(sum, coinB)
			indexB++
		}
	}
}

// IsEqual returns true if the two sets of Coins have the same value
func (coins Coins) IsEqual(coinsB Coins) bool {
	if len(coins) != len(coinsB) {
		return false
	}
	for i := 0; i < len(coins); i++ {
		if coins[i].Denom != coinsB[i].Denom || !(coins[i].Amount == coinsB[i].Amount) {
			return false
		}
	}
	return true
}

func (coins Coins) IsZero() bool {
	for _, coin := range coins {
		if !coin.IsZero() {
			return false
		}
	}
	return true
}

func (coins Coins) IsNotNegative() bool {
	if len(coins) == 0 {
		return true
	}
	for _, coin := range coins {
		if !coin.IsNotNegative() {
			return false
		}
	}
	return true
}

func (coins Coins) AmountOf(denom string) int64 {
	switch len(coins) {
	case 0:
		return 0
	case 1:
		coin := coins[0]
		if coin.Denom == denom {
			return coin.Amount
		}
		return 0
	default:
		midIdx := len(coins) / 2 // 2:1, 3:1, 4:2
		coin := coins[midIdx]
		if denom < coin.Denom {
			return coins[:midIdx].AmountOf(denom)
		} else if denom == coin.Denom {
			return coin.Amount
		} else {
			return coins[midIdx+1:].AmountOf(denom)
		}
	}
}

// Sort interface

//nolint
func (coins Coins) Len() int           { return len(coins) }
func (coins Coins) Less(i, j int) bool { return coins[i].Denom < coins[j].Denom }
func (coins Coins) Swap(i, j int)      { coins[i], coins[j] = coins[j], coins[i] }

// Sort is a helper function to sort the set of coins inplace
func (coins Coins) Sort() Coins {
	sort.Sort(coins)
	return coins
}
