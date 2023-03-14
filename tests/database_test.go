package tests

import (
	"testing"

	"github.com/n0033/stock-game/services/database"
)

func TestGetDatabase(t *testing.T) {
	db := database.GetDatabase()
	if db == nil {
		t.Fail()
	}
}
