package models

type SchemaColumnsInfo struct {
	ColumnName string
	DataType   string
	Tags       string
	Seq        byte
}
type ColumnsSortModel []*SchemaColumnsInfo

func (c ColumnsSortModel) Len() int {
	return len(c)
}

func (c ColumnsSortModel) Less(i, j int) bool {
	return c[i].Seq < c[j].Seq
}
func (c ColumnsSortModel) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

