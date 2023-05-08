package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Ptr:
**************************************/
type Ptr struct {
	tCommon
	Base ValueType
}

func (m *Module) GenValueType_Ptr(base ValueType) *Ptr {
	ptr_t := Ptr{Base: base}
	t, ok := m.findValueType(ptr_t.Name())
	if ok {
		return t.(*Ptr)
	}

	m.addValueType(&ptr_t)
	return &ptr_t
}

func (t *Ptr) Name() string         { return t.Base.Name() + ".$$ptr" }
func (t *Ptr) Size() int            { return 4 }
func (t *Ptr) align() int           { return 4 }
func (t *Ptr) Kind() TypeKind       { return kPtr }
func (t *Ptr) onFree() int          { return 0 }
func (t *Ptr) Raw() []wat.ValueType { return []wat.ValueType{toWatType(t)} }
func (t *Ptr) Equal(u ValueType) bool {
	if ut, ok := u.(*Ptr); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}
func (t *Ptr) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
aPtr:
**************************************/
type aPtr struct {
	aBasic
}

func newValue_Ptr(name string, kind ValueKind, typ *Ptr) *aPtr {
	var v aPtr
	v.aValue = aValue{name: name, kind: kind, typ: typ}
	return &v
}

func (v *aPtr) emitGetValue() []wat.Inst {
	t := v.Type().(*Ptr).Base
	return t.EmitLoadFromAddr(v, 0)
}

func (v *aPtr) emitSetValue(d Value) []wat.Inst {
	if !d.Type().Equal(v.Type().(*Ptr).Base) {
		logger.Fatal("Type not match")
		return nil
	}
	return d.emitStoreToAddr(v, 0)
}

func (v *aPtr) Bin() (b []byte) {
	if v.Kind() != ValueKindConst {
		panic("Value.bin(): const only!")
	}

	b = make([]byte, 4)
	i, _ := strconv.Atoi(v.Name())
	b[0] = byte(i & 0xFF)
	b[1] = byte((i >> 8) & 0xFF)
	b[2] = byte((i >> 16) & 0xFF)
	b[3] = byte((i >> 24) & 0xFF)

	return
}
