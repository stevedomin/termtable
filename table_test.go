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

func (t *Table) OldRender() string {
	var tableStr string
	i := 0

	if t.HasHeader {
		if t.Options.UseSeparator {
			tableStr += t.separatorLine() + "\n"
		}
		for j := range t.Rows[0] {
			tableStr += t.getCell(i, j)
		}
		i = 1
		tableStr += "\n"
	}

	if t.Options.UseSeparator {
		tableStr += t.separatorLine() + "\n"
	}

	for i < len(t.Rows) {
		row := t.Rows[i]
		for j := range row {
			tableStr += t.getCell(i, j)
		}
		if i < len(t.Rows)-1 {
			tableStr += "\n"
		}
		i++
	}

	if t.Options.UseSeparator {
		tableStr += "\n" + t.separatorLine()
	}

	return tableStr
}

func BenchmarkRenderBuffer(b *testing.B) {
	b.StopTimer()
	ta := NewTable(nil, nil)
	ta.SetHeader([]string{"LOWERCASE", "UPPERCASE", "NUMBERS"})
	ta.AddRow([]string{"abc", "ABCD", "12345"})
	ta.AddRow([]string{"defg", "EFGHI", "678"})
	ta.AddRow([]string{"hijkl", "JKL", "9000"})
	ta.AddRow([]string{"defg", "EFGHI", "678"})
	ta.AddRow([]string{"hijkl", "JKL", "9000"})
	ta.AddRow([]string{"defg", "EFGHI", "678"})
	ta.AddRow([]string{"hijkl", "JKL", "9000"})
	ta.AddRow([]string{"defg", "EFGHI", "678"})
	ta.AddRow([]string{"hijkl", "JKL", "9000"})

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ta.Render()
	}
}

func BenchmarkRenderString(b *testing.B) {
	b.StopTimer()
	ta := NewTable(nil, nil)
	ta.SetHeader([]string{"LOWERCASE", "UPPERCASE", "NUMBERS"})
	ta.AddRow([]string{"abc", "ABCD", "12345"})
	ta.AddRow([]string{"defg", "EFGHI", "678"})
	ta.AddRow([]string{"hijkl", "JKL", "9000"})
	ta.AddRow([]string{"defg", "EFGHI", "678"})
	ta.AddRow([]string{"hijkl", "JKL", "9000"})
	ta.AddRow([]string{"defg", "EFGHI", "678"})
	ta.AddRow([]string{"hijkl", "JKL", "9000"})
	ta.AddRow([]string{"defg", "EFGHI", "678"})
	ta.AddRow([]string{"hijkl", "JKL", "9000"})

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ta.OldRender()
	}
}
