/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"bytes"
	"github.com/hedzr/cmdr"
	"gopkg.in/hedzr/errors.v3"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

type User struct {
	Name     string
	Birthday *time.Time
	Nickname string
	Role     string
	Age      int32
	Retry    int8
	Times    int16
	RetryU   uint8
	TimesU   uint16
	FxReal   uint32
	FxTime   int64
	FxTimeU  uint64
	UxA      uint
	UxB      int
	FakeAge  *int32
	Notes    []string
	flags    []byte
	Born     *int
	BornU    *uint
	Ro       []int
	F11      float32
	F12      float64
	C11      complex64
	C12      complex128
	Sptr     *string
	Bool1    bool
	Bool2    bool
	// Feat     []byte
}

func (user User) DoubleAge() int32 {
	return 2 * user.Age
}

type Employee struct {
	Name      string
	Birthday  *time.Time
	F11       float32
	F12       float64
	C11       complex64
	C12       complex128
	Feat      []byte
	Sptr      *string
	Nickname  *string
	Age       int64
	FakeAge   int
	EmployeID int64
	DoubleAge int32
	SuperRule string
	Notes     []string
	RetryU    uint8
	TimesU    uint16
	FxReal    uint32
	FxTime    int64
	FxTimeU   uint64
	UxA       uint
	UxB       int
	Retry     int8
	Times     int16
	Born      *int
	BornU     *uint
	flags     []byte //nolint:structcheck,unused
	Bool1     bool
	Bool2     bool
	Ro        []int
}

type X0 struct{}

type X1 struct {
	A uintptr
	B map[string]interface{}
	C bytes.Buffer
	D []string
	E []*X0
	F chan struct{}
	G chan bool
	H chan int
	I func()
	J interface{}
	K *X0
	L unsafe.Pointer
	M unsafe.Pointer
	N []int
	O [2]string
	P [2]string
	Q [2]string
}

type X2 struct {
	A uintptr
	B map[string]interface{}
	C bytes.Buffer
	D []string
	E []*X0
	F chan struct{}
	G chan bool
	H chan int
	I func()
	J interface{}
	K *X0
	L unsafe.Pointer
	M unsafe.Pointer
	N []int
	O [2]string
	P [2]string
	Q [3]string
}

func TestCopyCov(t *testing.T) {
	nn := []int{2, 9, 77, 111, 23, 29}
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	x0 := X0{}
	x1 := X1{
		A: uintptr(unsafe.Pointer(&x0)),
		H: make(chan int, 5),
		M: unsafe.Pointer(&x0),
		// E: []*X0{&x0},
		N: nn[1:3],
		O: a,
		Q: a,
	}
	x2 := &X2{N: nn[1:3]}
	_ = cmdr.GormDefaultCopier.Copy(&x2, &x1, "Shit", "Memo", "Name")
}

func (employee *Employee) Role(role string) {
	employee.SuperRule = "Super " + role
}

func checkEmployee(employee Employee, user User, t *testing.T, testCase string) {
	if employee.Name != user.Name {
		t.Errorf("%v: Name haven't been copied correctly.", testCase)
	}
	if employee.Nickname == nil || *employee.Nickname != user.Nickname {
		t.Errorf("%v: NickName haven't been copied correctly.", testCase)
	}
	if employee.Birthday == nil && user.Birthday != nil {
		t.Errorf("%v: Birthday haven't been copied correctly.", testCase)
	}
	if employee.Birthday != nil && user.Birthday == nil {
		t.Errorf("%v: Birthday haven't been copied correctly.", testCase)
	}
	if employee.Age != int64(user.Age) {
		t.Errorf("%v: Age haven't been copied correctly.", testCase)
	}
	if user.FakeAge != nil && employee.FakeAge != int(*user.FakeAge) {
		t.Errorf("%v: FakeAge haven't been copied correctly.", testCase)
	}
	if employee.DoubleAge != user.DoubleAge() {
		t.Errorf("%v: Copy from method doesn't work", testCase)
	}
	if employee.SuperRule != "Super "+user.Role {
		t.Errorf("%v: Copy to method doesn't work", testCase)
	}
	if !reflect.DeepEqual(employee.Notes, user.Notes) {
		t.Errorf("%v: Copy from slice doen't work", testCase)
	}
}

func TestCopyTwoStruct(t *testing.T) {
	user := User{Name: "Real Faked"}
	userTo := User{Name: "Faked", Role: "NN"}

	_ = cmdr.GormDefaultCopier.Copy(&userTo, &user, "Shit", "Memo")

	if userTo.Name != user.Name || userTo.Role != "NN" {
		t.Fatal("wrong")
	}
}

func TestCopyStruct(t *testing.T) {
	var fakeAge int32 = 12
	var born = 7
	var bornU uint = 7
	var sz = "dablo"
	user := User{Name: "Faked"}
	employee := Employee{}

	if err := cmdr.StandardCopier.Copy(employee, &user); err == nil {
		t.Error("Copy to unaddressable value should get error")
	}

	_ = cmdr.GormDefaultCopier.Copy(&employee, &user, "Shit", "Memo", "Name")
	// cmdr.StandardCopier.Copy(&employee, &user, "Shit", "Memo", "Name")

	user = User{Name: "Faked", Nickname: "user", Age: 18, FakeAge: &fakeAge,
		Role: "User", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'},
		Retry: 3, Times: 17, RetryU: 23, TimesU: 21, FxReal: 31, FxTime: 37,
		FxTimeU: 13, UxA: 11, UxB: 0, Born: &born, BornU: &bornU,
		Ro: []int{1, 2, 3}, Sptr: &sz, Bool1: true, // Feat: []byte(sz),
	}
	employee = Employee{}
	err := cmdr.StandardCopier.Copy(&employee, &user)
	if err != nil {
		t.Errorf("%v", err)
	}
	checkEmployee(employee, user, t, "Copy From Ptr To Ptr")

	employee2 := Employee{}
	_ = cmdr.StandardCopier.Copy(&employee2, user)
	checkEmployee(employee2, user, t, "Copy From Struct To Ptr")

	employee3 := Employee{}
	ptrToUser := &user
	_ = cmdr.StandardCopier.Copy(&employee3, &ptrToUser)
	checkEmployee(employee3, user, t, "Copy From Double Ptr To Ptr")

	employee4 := &Employee{}
	_ = cmdr.StandardCopier.Copy(&employee4, user)
	checkEmployee(*employee4, user, t, "Copy From Ptr To Double Ptr")
}

func TestCopyFromStructToSlice(t *testing.T) {
	user := User{Name: "Faked", Age: 18, Role: "User", Notes: []string{"hello world"}}
	var employees []Employee

	if err := cmdr.StandardCopier.Copy(employees, &user); err != nil && len(employees) != 0 {
		t.Errorf("Copy to unaddressable value should get error")
	}

	if _ = cmdr.StandardCopier.Copy(&employees, &user); len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(employees[0], user, t, "Copy From Struct To Slice Ptr")
	}

	employees2 := &[]Employee{}
	if _ = cmdr.StandardCopier.Copy(&employees2, user); len(*employees2) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee((*employees2)[0], user, t, "Copy From Struct To Double Slice Ptr")
	}

	employees3 := []*Employee{}
	if _ = cmdr.StandardCopier.Copy(&employees3, user); len(employees3) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(*(employees3[0]), user, t, "Copy From Struct To Ptr Slice Ptr")
	}

	employees4 := &[]*Employee{}
	if _ = cmdr.StandardCopier.Copy(&employees4, user); len(*employees4) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(*((*employees4)[0]), user, t, "Copy From Struct To Double Ptr Slice Ptr")
	}
}

func TestCopyFromSliceToSlice(t *testing.T) {
	users := []User{{Name: "Faked", Age: 18, Role: "User", Notes: []string{"hello world", "chaos"}}, {Name: "Real", Age: 22, Role: "World", Notes: []string{"hello world", "hello", "winner"}}}
	var employees []Employee

	if _ = cmdr.StandardCopier.Copy(&employees, users); len(employees) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(employees[0], users[0], t, "Copy From Slice To Slice Ptr @ 1")
		checkEmployee(employees[1], users[1], t, "Copy From Slice To Slice Ptr @ 2")
	}

	employees2 := &[]Employee{}
	if _ = cmdr.StandardCopier.Copy(&employees2, &users); len(*employees2) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee((*employees2)[0], users[0], t, "Copy From Slice Ptr To Double Slice Ptr @ 1")
		checkEmployee((*employees2)[1], users[1], t, "Copy From Slice Ptr To Double Slice Ptr @ 2")
	}

	var employees3 []*Employee
	if _ = cmdr.StandardCopier.Copy(&employees3, users); len(employees3) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(*(employees3[0]), users[0], t, "Copy From Slice To Ptr Slice Ptr @ 1")
		checkEmployee(*(employees3[1]), users[1], t, "Copy From Slice To Ptr Slice Ptr @ 2")
	}

	employees4 := &[]*Employee{}
	if _ = cmdr.StandardCopier.Copy(&employees4, users); len(*employees4) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(*((*employees4)[0]), users[0], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 1")
		checkEmployee(*((*employees4)[1]), users[1], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 2")
	}
}

func TestEmbedded(t *testing.T) {
	type Base struct {
		BaseField1 int
		BaseField2 int
	}

	type Embed struct {
		EmbedField1 int
		EmbedField2 int
		Base
	}

	base := Base{}
	embedded := Embed{}
	embedded.BaseField1 = 1
	embedded.BaseField2 = 2
	embedded.EmbedField1 = 3
	embedded.EmbedField2 = 4

	_ = cmdr.StandardCopier.Copy(&base, &embedded)

	if base.BaseField1 != 1 {
		t.Error("Embedded fields not copied")
	}

	if err := cmdr.GormDefaultCopier.Copy(&base, &embedded); err != nil {
		t.Error(err)
	}
}

type structSameName1 struct {
	A string
	B int64
	C time.Time
}

type structSameName2 struct {
	A string
	B time.Time
	C int64
}

func TestCopyFieldsWithSameNameButDifferentTypes(t *testing.T) {
	obj1 := structSameName1{A: "123", B: 2, C: time.Now()}
	obj2 := &structSameName2{}
	err := cmdr.StandardCopier.Copy(obj2, &obj1)
	if err != nil {
		t.Error("Should not raise error")
	}

	if obj2.A != obj1.A {
		t.Errorf("Field A should be copied")
	}

	err = cmdr.GormDefaultCopier.Copy(obj2, &obj1)
	if err != nil {
		t.Error(err)
	}
}

type ScannerValue struct {
	V int
}

func (s *ScannerValue) Scan(src interface{}) error {
	return errors.New("I failed")
}

type ScannerStruct struct {
	V *ScannerValue
}

type ScannerStructTo struct {
	V *ScannerValue
}

func TestScanner(t *testing.T) {
	s := &ScannerStruct{
		V: &ScannerValue{
			V: 12,
		},
	}

	s2 := &ScannerStructTo{}

	err := cmdr.StandardCopier.Copy(s2, s)
	if err != nil {
		t.Error("Should not raise error")
	}

	if s.V.V != s2.V.V {
		t.Errorf("Field V should be copied")
	}
}

func TestClone(t *testing.T) {
	var aa = "dsajkld"
	var b int

	// cmdr.Clone(b, aa)

	cmdr.Clone(b, &aa)

	cmdr.Clone(&b, &aa)

	cmdr.Clone(&b, nil)

	var c, d bool
	cmdr.Clone(&c, &d)

	var e, f int
	cmdr.Clone(&e, &f)
	var e1, f1 int8
	f1 = 1
	cmdr.Clone(&e1, &f1)
	if e1 != 1 {
		t.Failed()
	}
	var e2, f2 int16
	cmdr.Clone(&e2, &f2)
	var e3, f3 int32
	cmdr.Clone(&e3, &f3)
	var e4, f4 int64
	f4 = 9
	cmdr.Clone(&e4, &f4)
	if e1 != 9 {
		t.Failed()
	}

	var g, h string
	cmdr.Clone(&g, &h)
}
