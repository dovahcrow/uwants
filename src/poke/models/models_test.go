package models

import (
	"testing"
)

func Test_Delete(t *testing.T) {
	t.Log(DeleteAllFids())
}

func Test_Create(t *testing.T) {
	t.Log(InsertTids([]string{"123", "456"}))
}

func Test_Get(t *testing.T) {
	t.Log(GetAllFids())
}
