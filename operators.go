package athena

type Operator struct {
	Column    string    `json:"column"`
	Operation operation `json:"operation"`
	Value     OpValue   `json:"value"`
}

type OpValue interface{}

type operation string

const (
	EqualOp              operation = "="
	NotEqualOp           operation = "!="
	GreaterThanOp        operation = ">"
	GreaterThanEqualToOp operation = ">="
	LessThanOp           operation = "<"
	LessThanEqualToOp    operation = "<="
	InOp                 operation = "in"
	NotInOp              operation = "not in"
	LikeOp               operation = "like"

	LimitOp  operation = "limit"
	OrderOp  operation = "order"
	SkipOp   operation = "skip"
	OrOp     operation = "or"
	AndOp    operation = "and"
	ExistsOp operation = "exists"
)

var AllOperations = []operation{
	EqualOp,
	NotEqualOp,
	GreaterThanOp,
	GreaterThanEqualToOp,
	LessThanOp,
	LessThanEqualToOp,
	InOp,
	NotInOp,
	LikeOp,
	LimitOp,
	OrderOp,
	SkipOp,
	OrOp,
	AndOp,
	ExistsOp,
}

func (o operation) IsValid() bool {
	switch o {
	case EqualOp, NotEqualOp,
		GreaterThanOp, LessThanOp, GreaterThanEqualToOp, LessThanEqualToOp,
		InOp, NotInOp, LikeOp,
		LimitOp, OrderOp, SkipOp, OrOp, AndOp, ExistsOp:
		return true
	}
	return false
}

func (o operation) Value() string {
	return string(o)
}

func NewOperators(operators ...*Operator) []*Operator {
	return operators
}

func NewLikeOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: LikeOp,
		Value:     value,
	}
}

func NewEqualOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: EqualOp,
		Value:     value,
	}
}

func NewNotEqualOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: NotEqualOp,
		Value:     value,
	}
}

func NewGreaterThanOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: GreaterThanOp,
		Value:     value,
	}
}

func NewGreaterThanEqualToOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: GreaterThanEqualToOp,
		Value:     value,
	}
}

func NewLessThanOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: LessThanOp,
		Value:     value,
	}
}

func NewLessThanEqualToOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: LessThanEqualToOp,
		Value:     value,
	}
}

type Sort int

const (
	SortAsc  Sort = 1
	SortDesc Sort = -1
)

var AllSort = []Sort{
	SortAsc,
	SortDesc,
}

func (e Sort) IsValid() bool {
	switch e {
	case SortAsc, SortDesc:
		return true
	}
	return false
}

func (e Sort) Value() int {
	return int(e)
}

func NewOrderOperator(column string, sort Sort) *Operator {

	if !sort.IsValid() {
		return nil
	}

	return &Operator{
		Column:    column,
		Operation: OrderOp,
		Value:     sort.Value(),
	}

}

func NewInOperator(column string, values ...interface{}) *Operator {

	return &Operator{
		Column:    column,
		Operation: InOp,
		Value:     values,
	}

}

func NewNotInOperator(column string, value interface{}) *Operator {

	return &Operator{
		Column:    column,
		Operation: NotInOp,
		Value:     value,
	}

}

func NewLimitOperator(value int64) *Operator {
	return &Operator{
		Column:    "",
		Operation: LimitOp,
		Value:     value,
	}
}

func NewSkipOperator(value int64) *Operator {
	return &Operator{
		Column:    "",
		Operation: SkipOp,
		Value:     value,
	}
}

func NewOrOperator(value ...*Operator) *Operator {
	return &Operator{
		Column:    "",
		Operation: OrOp,
		Value:     value,
	}
}

func NewAndOperator(value ...*Operator) *Operator {
	return &Operator{
		Column:    "",
		Operation: AndOp,
		Value:     value,
	}
}

func NewExistsOperator(column string, value bool) *Operator {
	return &Operator{
		Column:    column,
		Operation: ExistsOp,
		Value:     value,
	}
}
