package termtable

import (
	"bytes"
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

	hasRendered     bool
	numRenderedRows int

	isRenderDyn bool
	numDynRows  int
	dynamicRows map[int]bool
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
		Options:     options,
		dynamicRows: make(map[int]bool),
	}
	if t.Options == nil {
		t.Options = defaultTableOptions
	}
	if rows != nil {
		t.Rows = rows
		t.computeProperties()
	}
	return t
}

func (t *Table) SetHeader(header []string) {
	t.HasHeader = true
	// There is a better way to do this
	t.Rows = append([][]string{header}, t.Rows...)
	t.computeProperties()
}

func (t *Table) AddRow(row []string) {
	t.Rows = append(t.Rows, row)
	if t.hasRendered {
		t.numDynRows++
	} else {
		t.computeProperties()
	}
}

func (t *Table) computeProperties() {
	if len(t.Rows) > 0 {
		t.numColumns = len(t.Rows[0])
		t.columnsWidth = make([]int, t.numColumns)
		t.recalculate()
	}
}

func (t *Table) recalculate() {
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

func (t *Table) Render() string {
	t.isRenderDyn = false
	// allocate a 1k byte buffer
	bb := make([]byte, 0, 1024)
	buf := bytes.NewBuffer(bb)

	i := 0

	if t.HasHeader {
		if t.Options.UseSeparator {
			buf.WriteString(t.separatorLine())
			buf.WriteRune('\n')
		}
		for j := range t.Rows[0] {
			buf.WriteString(t.getCell(i, j))
		}
		i = 1
		buf.WriteRune('\n')
	}

	if t.Options.UseSeparator {
		buf.WriteString(t.separatorLine())
		buf.WriteRune('\n')
	}

	for i < len(t.Rows) {
		row := t.Rows[i]
		for j := range row {
			buf.WriteString(t.getCell(i, j))
		}
		if i < len(t.Rows)-t.numDynRows-1 {
			buf.WriteRune('\n')
		}
		i++
		t.numRenderedRows++
	}

	if t.Options.UseSeparator {
		buf.WriteRune('\n')
		buf.WriteString(t.separatorLine())
	}

	t.hasRendered = true

	return buf.String()
}

func (t *Table) RenderDynamic() string {
	if t.numDynRows == 0 {
		return t.Render()
	}
	t.isRenderDyn = true

	// allocate a 1k byte buffer
	bb := make([]byte, 0, 1024)
	buf := bytes.NewBuffer(bb)

	i := t.numRenderedRows
	if t.Options.UseSeparator {
		buf.WriteString("\033[1A")
	}

	for i < len(t.Rows) {
		row := t.Rows[i]
		for j := range row {
			buf.WriteString(t.getCell(i, j))
		}
		if i < len(t.Rows)-1 {
			buf.WriteRune('\n')
		}
		i++
		t.numRenderedRows++
	}

	if t.Options.UseSeparator {
		buf.WriteRune('\n')
		buf.WriteString(t.separatorLine())
	}

	return buf.String()
}

func (t *Table) separatorLine() string {
	sep := "+"
	for _, w := range t.columnsWidth {
		sep += strings.Repeat("-", w+2*t.Options.Padding)
		sep += "+"
	}
	return sep
}

func (t *Table) getCell(row, col int) string {
	cellContent := t.Rows[row][col]
	//log.Println(cellContent)
	if t.isRenderDyn {
		colWidth := t.columnsWidth[col]
		if len(cellContent) > colWidth {
			cellContent = t.handleCellOverflow(row, col)
		}
	}
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

// If a dynamic row has a greater width then the current computed length for
// the given column, trim and wrapped to the next row down
func (t *Table) handleCellOverflow(row, col int) string {
	origCellContent := t.Rows[row][col]

	index := t.columnsWidth[col]
	if tindex := strings.LastIndex(origCellContent[:index], " "); tindex != -1 {
		_, isDynRow := t.dynamicRows[row]
		if !isDynRow || (isDynRow && tindex != 1) {
			index = tindex
		}
	}

	// trim content, create new row
	trimCellContent := origCellContent[:index]
	newCellContent := "> " + strings.Trim(origCellContent[index:], " ")
	_, isNextDynRow := t.dynamicRows[row+1]
	if isNextDynRow {
		t.Rows[row+1][col] = newCellContent
	} else {
		newRow := make([]string, t.numColumns)
		newRow[col] = newCellContent
		i := row + 1
		// insert new row into it's appropriate spotgit st
		t.Rows = append(t.Rows, []string{})
		copy(t.Rows[i+1:], t.Rows[i:])
		t.Rows[i] = newRow

		t.dynamicRows[row+1] = true
	}

	return trimCellContent
}
