package termtable

import (
	"testing"
)

func TestTermtable(t *testing.T) {
	ta := NewTable(nil, nil)
	ta.SetHeader([]string{"LOWERCASE", "UPPERCASE", "NUMBERS"})
	ta.AddRow([]string{"abc", "ABCD", "12345"})
	ta.AddRow([]string{"defg", "EFGHI", "678"})
	ta.AddRow([]string{"hijkl", "JKL", "9000"})
	s := ta.Render()
	if s != "LOWERCASE UPPERCASE NUMBERS \nabc       ABCD      12345   \ndefg      EFGHI     678     \nhijkl     JKL       9000    " {
		t.Fatalf("Got unexpected result")
	}
}
