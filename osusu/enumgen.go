// Code generated by "enumgen -sql"; DO NOT EDIT.

package osusu

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"goki.dev/enums"
)

var _SourcesValues = []Sources{0, 1, 2, 3}

// SourcesN is the highest valid value
// for type Sources, plus one.
const SourcesN Sources = 4

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the enumgen command to generate them again.
func _SourcesNoOp() {
	var x [1]struct{}
	_ = x[Cooking-(0)]
	_ = x[DineIn-(1)]
	_ = x[Takeout-(2)]
	_ = x[Delivery-(3)]
}

var _SourcesNameToValueMap = map[string]Sources{
	`Cooking`:  0,
	`cooking`:  0,
	`DineIn`:   1,
	`dinein`:   1,
	`Takeout`:  2,
	`takeout`:  2,
	`Delivery`: 3,
	`delivery`: 3,
}

var _SourcesDescMap = map[Sources]string{
	0: ``,
	1: ``,
	2: ``,
	3: ``,
}

var _SourcesMap = map[Sources]string{
	0: `Cooking`,
	1: `DineIn`,
	2: `Takeout`,
	3: `Delivery`,
}

// String returns the string representation
// of this Sources value.
func (i Sources) String() string {
	str := ""
	for _, ie := range _SourcesValues {
		if i.HasFlag(ie) {
			ies := ie.BitIndexString()
			if str == "" {
				str = ies
			} else {
				str += "|" + ies
			}
		}
	}
	return str
}

// BitIndexString returns the string
// representation of this Sources value
// if it is a bit index value
// (typically an enum constant), and
// not an actual bit flag value.
func (i Sources) BitIndexString() string {
	if str, ok := _SourcesMap[i]; ok {
		return str
	}
	return strconv.FormatInt(int64(i), 10)
}

// SetString sets the Sources value from its
// string representation, and returns an
// error if the string is invalid.
func (i *Sources) SetString(s string) error {
	*i = 0
	return i.SetStringOr(s)
}

// SetStringOr sets the Sources value from its
// string representation while preserving any
// bit flags already set, and returns an
// error if the string is invalid.
func (i *Sources) SetStringOr(s string) error {
	flgs := strings.Split(s, "|")
	for _, flg := range flgs {
		if val, ok := _SourcesNameToValueMap[flg]; ok {
			i.SetFlag(true, &val)
		} else if val, ok := _SourcesNameToValueMap[strings.ToLower(flg)]; ok {
			i.SetFlag(true, &val)
		} else if flg == "" {
			continue
		} else {
			return fmt.Errorf("%q is not a valid value for type Sources", flg)
		}
	}
	return nil
}

// Int64 returns the Sources value as an int64.
func (i Sources) Int64() int64 {
	return int64(i)
}

// SetInt64 sets the Sources value from an int64.
func (i *Sources) SetInt64(in int64) {
	*i = Sources(in)
}

// Desc returns the description of the Sources value.
func (i Sources) Desc() string {
	if str, ok := _SourcesDescMap[i]; ok {
		return str
	}
	return i.String()
}

// SourcesValues returns all possible values
// for the type Sources.
func SourcesValues() []Sources {
	return _SourcesValues
}

// Values returns all possible values
// for the type Sources.
func (i Sources) Values() []enums.Enum {
	res := make([]enums.Enum, len(_SourcesValues))
	for i, d := range _SourcesValues {
		res[i] = d
	}
	return res
}

// IsValid returns whether the value is a
// valid option for type Sources.
func (i Sources) IsValid() bool {
	_, ok := _SourcesMap[i]
	return ok
}

// HasFlag returns whether these
// bit flags have the given bit flag set.
func (i Sources) HasFlag(f enums.BitFlag) bool {
	return atomic.LoadInt64((*int64)(&i))&(1<<uint32(f.Int64())) != 0
}

// SetFlag sets the value of the given
// flags in these flags to the given value.
func (i *Sources) SetFlag(on bool, f ...enums.BitFlag) {
	var mask int64
	for _, v := range f {
		mask |= 1 << v.Int64()
	}
	in := int64(*i)
	if on {
		in |= mask
		atomic.StoreInt64((*int64)(i), in)
	} else {
		in &^= mask
		atomic.StoreInt64((*int64)(i), in)
	}
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i Sources) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *Sources) UnmarshalText(text []byte) error {
	return i.SetString(string(text))
}

// Scan implements the [driver.Valuer] interface.
func (i Sources) Value() (driver.Value, error) {
	return i.String(), nil
}

// Value implements the [sql.Scanner] interface.
func (i *Sources) Scan(value any) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value for type Sources: %[1]T(%[1]v)", value)
	}

	return i.SetString(str)
}

var _CategoriesValues = []Categories{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

// CategoriesN is the highest valid value
// for type Categories, plus one.
const CategoriesN Categories = 10

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the enumgen command to generate them again.
func _CategoriesNoOp() {
	var x [1]struct{}
	_ = x[Breakfast-(0)]
	_ = x[Brunch-(1)]
	_ = x[Lunch-(2)]
	_ = x[Dinner-(3)]
	_ = x[Dessert-(4)]
	_ = x[Snack-(5)]
	_ = x[Appetizer-(6)]
	_ = x[Side-(7)]
	_ = x[Drink-(8)]
	_ = x[Ingredient-(9)]
}

var _CategoriesNameToValueMap = map[string]Categories{
	`Breakfast`:  0,
	`breakfast`:  0,
	`Brunch`:     1,
	`brunch`:     1,
	`Lunch`:      2,
	`lunch`:      2,
	`Dinner`:     3,
	`dinner`:     3,
	`Dessert`:    4,
	`dessert`:    4,
	`Snack`:      5,
	`snack`:      5,
	`Appetizer`:  6,
	`appetizer`:  6,
	`Side`:       7,
	`side`:       7,
	`Drink`:      8,
	`drink`:      8,
	`Ingredient`: 9,
	`ingredient`: 9,
}

var _CategoriesDescMap = map[Categories]string{
	0: ``,
	1: ``,
	2: ``,
	3: ``,
	4: ``,
	5: ``,
	6: ``,
	7: ``,
	8: ``,
	9: ``,
}

var _CategoriesMap = map[Categories]string{
	0: `Breakfast`,
	1: `Brunch`,
	2: `Lunch`,
	3: `Dinner`,
	4: `Dessert`,
	5: `Snack`,
	6: `Appetizer`,
	7: `Side`,
	8: `Drink`,
	9: `Ingredient`,
}

// String returns the string representation
// of this Categories value.
func (i Categories) String() string {
	str := ""
	for _, ie := range _CategoriesValues {
		if i.HasFlag(ie) {
			ies := ie.BitIndexString()
			if str == "" {
				str = ies
			} else {
				str += "|" + ies
			}
		}
	}
	return str
}

// BitIndexString returns the string
// representation of this Categories value
// if it is a bit index value
// (typically an enum constant), and
// not an actual bit flag value.
func (i Categories) BitIndexString() string {
	if str, ok := _CategoriesMap[i]; ok {
		return str
	}
	return strconv.FormatInt(int64(i), 10)
}

// SetString sets the Categories value from its
// string representation, and returns an
// error if the string is invalid.
func (i *Categories) SetString(s string) error {
	*i = 0
	return i.SetStringOr(s)
}

// SetStringOr sets the Categories value from its
// string representation while preserving any
// bit flags already set, and returns an
// error if the string is invalid.
func (i *Categories) SetStringOr(s string) error {
	flgs := strings.Split(s, "|")
	for _, flg := range flgs {
		if val, ok := _CategoriesNameToValueMap[flg]; ok {
			i.SetFlag(true, &val)
		} else if val, ok := _CategoriesNameToValueMap[strings.ToLower(flg)]; ok {
			i.SetFlag(true, &val)
		} else if flg == "" {
			continue
		} else {
			return fmt.Errorf("%q is not a valid value for type Categories", flg)
		}
	}
	return nil
}

// Int64 returns the Categories value as an int64.
func (i Categories) Int64() int64 {
	return int64(i)
}

// SetInt64 sets the Categories value from an int64.
func (i *Categories) SetInt64(in int64) {
	*i = Categories(in)
}

// Desc returns the description of the Categories value.
func (i Categories) Desc() string {
	if str, ok := _CategoriesDescMap[i]; ok {
		return str
	}
	return i.String()
}

// CategoriesValues returns all possible values
// for the type Categories.
func CategoriesValues() []Categories {
	return _CategoriesValues
}

// Values returns all possible values
// for the type Categories.
func (i Categories) Values() []enums.Enum {
	res := make([]enums.Enum, len(_CategoriesValues))
	for i, d := range _CategoriesValues {
		res[i] = d
	}
	return res
}

// IsValid returns whether the value is a
// valid option for type Categories.
func (i Categories) IsValid() bool {
	_, ok := _CategoriesMap[i]
	return ok
}

// HasFlag returns whether these
// bit flags have the given bit flag set.
func (i Categories) HasFlag(f enums.BitFlag) bool {
	return atomic.LoadInt64((*int64)(&i))&(1<<uint32(f.Int64())) != 0
}

// SetFlag sets the value of the given
// flags in these flags to the given value.
func (i *Categories) SetFlag(on bool, f ...enums.BitFlag) {
	var mask int64
	for _, v := range f {
		mask |= 1 << v.Int64()
	}
	in := int64(*i)
	if on {
		in |= mask
		atomic.StoreInt64((*int64)(i), in)
	} else {
		in &^= mask
		atomic.StoreInt64((*int64)(i), in)
	}
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i Categories) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *Categories) UnmarshalText(text []byte) error {
	return i.SetString(string(text))
}

// Scan implements the [driver.Valuer] interface.
func (i Categories) Value() (driver.Value, error) {
	return i.String(), nil
}

// Value implements the [sql.Scanner] interface.
func (i *Categories) Scan(value any) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value for type Categories: %[1]T(%[1]v)", value)
	}

	return i.SetString(str)
}

var _CuisinesValues = []Cuisines{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

// CuisinesN is the highest valid value
// for type Cuisines, plus one.
const CuisinesN Cuisines = 17

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the enumgen command to generate them again.
func _CuisinesNoOp() {
	var x [1]struct{}
	_ = x[African-(0)]
	_ = x[American-(1)]
	_ = x[Asian-(2)]
	_ = x[British-(3)]
	_ = x[Chinese-(4)]
	_ = x[European-(5)]
	_ = x[French-(6)]
	_ = x[Greek-(7)]
	_ = x[Indian-(8)]
	_ = x[Italian-(9)]
	_ = x[Japanese-(10)]
	_ = x[Jewish-(11)]
	_ = x[Korean-(12)]
	_ = x[LatinAmerican-(13)]
	_ = x[Mexican-(14)]
	_ = x[MiddleEastern-(15)]
	_ = x[Thai-(16)]
}

var _CuisinesNameToValueMap = map[string]Cuisines{
	`African`:       0,
	`african`:       0,
	`American`:      1,
	`american`:      1,
	`Asian`:         2,
	`asian`:         2,
	`British`:       3,
	`british`:       3,
	`Chinese`:       4,
	`chinese`:       4,
	`European`:      5,
	`european`:      5,
	`French`:        6,
	`french`:        6,
	`Greek`:         7,
	`greek`:         7,
	`Indian`:        8,
	`indian`:        8,
	`Italian`:       9,
	`italian`:       9,
	`Japanese`:      10,
	`japanese`:      10,
	`Jewish`:        11,
	`jewish`:        11,
	`Korean`:        12,
	`korean`:        12,
	`LatinAmerican`: 13,
	`latinamerican`: 13,
	`Mexican`:       14,
	`mexican`:       14,
	`MiddleEastern`: 15,
	`middleeastern`: 15,
	`Thai`:          16,
	`thai`:          16,
}

var _CuisinesDescMap = map[Cuisines]string{
	0:  ``,
	1:  ``,
	2:  ``,
	3:  ``,
	4:  ``,
	5:  ``,
	6:  ``,
	7:  ``,
	8:  ``,
	9:  ``,
	10: ``,
	11: ``,
	12: ``,
	13: ``,
	14: ``,
	15: ``,
	16: ``,
}

var _CuisinesMap = map[Cuisines]string{
	0:  `African`,
	1:  `American`,
	2:  `Asian`,
	3:  `British`,
	4:  `Chinese`,
	5:  `European`,
	6:  `French`,
	7:  `Greek`,
	8:  `Indian`,
	9:  `Italian`,
	10: `Japanese`,
	11: `Jewish`,
	12: `Korean`,
	13: `LatinAmerican`,
	14: `Mexican`,
	15: `MiddleEastern`,
	16: `Thai`,
}

// String returns the string representation
// of this Cuisines value.
func (i Cuisines) String() string {
	str := ""
	for _, ie := range _CuisinesValues {
		if i.HasFlag(ie) {
			ies := ie.BitIndexString()
			if str == "" {
				str = ies
			} else {
				str += "|" + ies
			}
		}
	}
	return str
}

// BitIndexString returns the string
// representation of this Cuisines value
// if it is a bit index value
// (typically an enum constant), and
// not an actual bit flag value.
func (i Cuisines) BitIndexString() string {
	if str, ok := _CuisinesMap[i]; ok {
		return str
	}
	return strconv.FormatInt(int64(i), 10)
}

// SetString sets the Cuisines value from its
// string representation, and returns an
// error if the string is invalid.
func (i *Cuisines) SetString(s string) error {
	*i = 0
	return i.SetStringOr(s)
}

// SetStringOr sets the Cuisines value from its
// string representation while preserving any
// bit flags already set, and returns an
// error if the string is invalid.
func (i *Cuisines) SetStringOr(s string) error {
	flgs := strings.Split(s, "|")
	for _, flg := range flgs {
		if val, ok := _CuisinesNameToValueMap[flg]; ok {
			i.SetFlag(true, &val)
		} else if val, ok := _CuisinesNameToValueMap[strings.ToLower(flg)]; ok {
			i.SetFlag(true, &val)
		} else if flg == "" {
			continue
		} else {
			return fmt.Errorf("%q is not a valid value for type Cuisines", flg)
		}
	}
	return nil
}

// Int64 returns the Cuisines value as an int64.
func (i Cuisines) Int64() int64 {
	return int64(i)
}

// SetInt64 sets the Cuisines value from an int64.
func (i *Cuisines) SetInt64(in int64) {
	*i = Cuisines(in)
}

// Desc returns the description of the Cuisines value.
func (i Cuisines) Desc() string {
	if str, ok := _CuisinesDescMap[i]; ok {
		return str
	}
	return i.String()
}

// CuisinesValues returns all possible values
// for the type Cuisines.
func CuisinesValues() []Cuisines {
	return _CuisinesValues
}

// Values returns all possible values
// for the type Cuisines.
func (i Cuisines) Values() []enums.Enum {
	res := make([]enums.Enum, len(_CuisinesValues))
	for i, d := range _CuisinesValues {
		res[i] = d
	}
	return res
}

// IsValid returns whether the value is a
// valid option for type Cuisines.
func (i Cuisines) IsValid() bool {
	_, ok := _CuisinesMap[i]
	return ok
}

// HasFlag returns whether these
// bit flags have the given bit flag set.
func (i Cuisines) HasFlag(f enums.BitFlag) bool {
	return atomic.LoadInt64((*int64)(&i))&(1<<uint32(f.Int64())) != 0
}

// SetFlag sets the value of the given
// flags in these flags to the given value.
func (i *Cuisines) SetFlag(on bool, f ...enums.BitFlag) {
	var mask int64
	for _, v := range f {
		mask |= 1 << v.Int64()
	}
	in := int64(*i)
	if on {
		in |= mask
		atomic.StoreInt64((*int64)(i), in)
	} else {
		in &^= mask
		atomic.StoreInt64((*int64)(i), in)
	}
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i Cuisines) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *Cuisines) UnmarshalText(text []byte) error {
	return i.SetString(string(text))
}

// Scan implements the [driver.Valuer] interface.
func (i Cuisines) Value() (driver.Value, error) {
	return i.String(), nil
}

// Value implements the [sql.Scanner] interface.
func (i *Cuisines) Scan(value any) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value for type Cuisines: %[1]T(%[1]v)", value)
	}

	return i.SetString(str)
}
