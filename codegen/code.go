package code

import (
	"fmt"
	"sync"
)

type Coder interface {
	// HTTPStatus HTTP status that should be used for the associated error code.
	HTTPStatus() int

	// String External (user) facing error text.
	String() string

	// Reference returns the detail documents for user.
	Reference() string

	// Code returns the code of the coder
	Code() uint32
}

type errCode struct {
	// C refers to the code of the errCode.
	C uint32

	// HTTP status that should be used for the associated error code.
	HTTP int

	// External (user) facing error text.
	Ext string

	// Ref specify the reference document.
	Ref string
}

func (e errCode) HTTPStatus() int {
	return e.HTTP
}

func (e errCode) String() string {
	return e.Ext
}

func (e errCode) Reference() string {
	return e.Ref
}

func (e errCode) Code() uint32 {
	return e.C
}

// codes contains a map of error codes to metadata.
var codes = map[uint32]errCode{}
var codeMux = &sync.Mutex{}

func mustRegister(code errCode) {
	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[code.C]; ok {
		panic(fmt.Sprintf("code: %d already exist", code.C))
	}

	codes[code.C] = code
}

func Register(code uint32, httpStatus int, message string, refs ...string) {
	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := errCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	}

	mustRegister(coder)
}

func GetCoder(code uint32) (Coder, bool) {
	coder, ok := codes[code]
	return coder, ok
}
