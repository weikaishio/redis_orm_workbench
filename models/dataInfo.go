package models

type DataConditionInfo struct {
	IdxNameKey      string
	CType           int
	IndividualValue string
	StartTime       uint32
	EndTime         uint32
	StartNumber     int
	EndNumber       int
	CType2           int
	IndividualValue2 string
	StartTime2       uint32
	EndTime2         uint32
	StartNumber2     int
	EndNumber2       int
}

const (
	CType_IndividualValue = 0
	CType_Time            = 1
	CType_Number          = 2
)

