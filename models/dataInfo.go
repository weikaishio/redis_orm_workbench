package models

type DataConditionInfo struct {
	ColumnName      string
	CType           int
	IndividualValue string
	StartTime       uint32
	EndTime         uint32
	StartNumber     int
	EndNumber       int
}

const (
	CType_IndividualValue = 0
	CType_Time            = 1
	CType_Number          = 2
)
