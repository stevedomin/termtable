package termtable

import (
	"fmt"
	"math"
	"strings"
)

type Table struct {
	Rows    [][]string
	Columns [][]string
	Options *TableOptions

	HasHeader bool

	numColumns   int
	columnsWidth []int
}

type TableOptions struct {
	Padding      int
	UseSeparator bool
}

var defaultTableOptions = &TableOptions{
	Padding:      1,
	UseSeparator: false,
}

func NewTable(rows [][]string, options *TableOptions) *Table {
	t := &Table{
		Options: options,
	}
	if t.Options == nil {
		t.Options = defaultTableOptions
	}
	if rows != nil {
		t.Rows = rows
		t.ComputeProperties()
	}
	return t
}

func (t *Table) SetHeader(header []string) {
	t.HasHeader = true
	// There is a better way to do this
	t.Rows = append([][]string{header}, t.Rows...)
	t.ComputeProperties()
}

func (t *Table) AddRow(row []string) {
	t.Rows = append(t.Rows, row)
	t.ComputeProperties()
}

func (t *Table) ComputeProperties() {
	if len(t.Rows) > 0 {
		t.numColumns = len(t.Rows[0])
		t.columnsWidth = make([]int, t.numColumns)
		t.Recalculate()
	}
}

func (t *Table) Recalculate() {
	t.Columns = [][]string{}
	for i := 0; i < t.numColumns; i++ {
		t.Columns = append(t.Columns, []string{})
	}
	for _, row := range t.Rows {
		for j, cellContent := range row {
			t.Columns[j] = append(t.Columns[j], cellContent)
			t.columnsWidth[j] = int(math.Max(float64(len(cellContent)), float64(t.columnsWidth[j])))
		}
	}
}

func (t *Table) Render() {
	var tableStr string
	i := 0

	if t.HasHeader {
		if t.Options.UseSeparator {
			tableStr += t.SeparatorLine() + "\n"
		}
		for j := range t.Rows[0] {
			tableStr += t.getCell(i, j)
		}
		i = 1
		tableStr += "\n"
	}

	if t.Options.UseSeparator {
		tableStr += t.SeparatorLine() + "\n"
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
		tableStr += "\n" + t.SeparatorLine()
	}

	fmt.Println(tableStr)
}

func (t *Table) SeparatorLine() string {
	sep := "+"
	for _, w := range t.columnsWidth {
		sep += strings.Repeat("-", w+2*t.Options.Padding)
		sep += "+"
	}
	return sep
}

func (t *Table) getCell(row, col int) string {
	cellContent := t.Rows[row][col]
	spacePadding := strings.Repeat(" ", t.Options.Padding)

	var cellStr string

	if t.Options.UseSeparator {
		cellStr += "|"
		cellStr += spacePadding
	}

	cellStr += cellContent
	cellStr += strings.Repeat(" ", t.columnsWidth[col]-len(cellContent))
	cellStr += spacePadding

	if t.Options.UseSeparator {
		if col == t.numColumns-1 {
			cellStr += "|"
		}
	}

	return cellStr
}
