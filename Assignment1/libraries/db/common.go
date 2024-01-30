package db

const (
	// NullValue To check whether a value is NULL in postgres
	NullValue = "null"
	// NotNullValue To check if a value is not NULL in postgres
	NotNullValue = "not null"
)

type RelationType int

const (
	NONE RelationType = iota
	EQUAL
	NOT_EQUAL
	IN
	NOT_IN
	IS
	LIKE
	ANY
	GT
	LT
)

func (r RelationType) String() string {
	switch r {
	case EQUAL:
		return "="
	case NOT_EQUAL:
		return "!="
	case IN:
		return "in"
	case NOT_IN:
		return "not in"
	case IS:
		return "is"
	case LIKE:
		return "like"
	case ANY:
		return "any"
	case GT:
		return ">"
	case LT:
		return "<"
	default:
		return ""
	}
}

type Cursor struct {
	PageNum      int
	PageSize     int
	TotalPages   uint32
	TotalRecords uint32
	PageToken    string
	OrderBy      string
	Limit        int
	Offset       int
}

type IndexParams struct {
	Name        string
	Type        string
	TableName   string
	ColumnNames []string
}

type WhereClauseType struct {
	ColumnName   string
	RelationType RelationType
	ColumnValue  interface{}
	TableAlias   string
	JsonOperator string
}
