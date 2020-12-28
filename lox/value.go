package lox

type Type int

const (
	TypeNil Type = iota
	TypeBool
	TypeNumber
	TypeString
)

type Nil struct{}
type Bool struct{ value bool }
type Number struct{ value float64 }
type String struct{ value string }

var nilInstance = Nil{}

type Value interface {
	Type() Type
	IsNil() bool
	AsBool() bool
	AsNumber() *float64
	AsString() *string
	Equal(other Value) bool
}

// nil

func NewNil() Nil {
	return nilInstance
}

func (x Nil) Type() Type {
	return TypeNil
}

func (x Nil) IsNil() bool {
	return true
}

func (x Nil) AsBool() bool {
	return false
}

func (x Nil) AsNumber() *float64 {
	return nil
}

func (x Nil) AsString() *string {
	return nil
}

func (x Nil) Equal(other Value) bool {
	return x == other
}

// bool

func NewBool(value bool) Bool {
	return Bool{value: value}
}

func (x Bool) Type() Type {
	return TypeBool
}

func (x Bool) IsNil() bool {
	return false
}

func (x Bool) AsBool() bool {
	return x.value
}

func (x Bool) AsNumber() *float64 {
	return nil
}

func (x Bool) AsString() *string {
	return nil
}

func (x Bool) Equal(other Value) bool {
	return other.Type() == TypeBool && x.value == other.(Bool).value
}

// number

func NewNumber(value float64) Number {
	return Number{value: value}
}

func (x Number) Type() Type {
	return TypeNumber
}

func (x Number) IsNil() bool {
	return false
}

func (x Number) AsBool() bool {
	return true
}

func (x Number) AsNumber() *float64 {
	return &x.value
}

func (x Number) AsString() *string {
	return nil
}

func (x Number) Equal(other Value) bool {
	return other.Type() == TypeNumber && x.value == other.(Number).value
}

// string

func NewString(value string) String {
	return String{value: value}
}

func (x String) Type() Type {
	return TypeString
}

func (x String) IsNil() bool {
	return false
}

func (x String) AsBool() bool {
	return true
}

func (x String) AsNumber() *float64 {
	return nil
}

func (x String) AsString() *string {
	return &x.value
}

func (x String) Equal(other Value) bool {
	return other.Type() == TypeString && x.value == other.(String).value
}
