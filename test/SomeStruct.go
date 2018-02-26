package somestruct

//SomeStruct used as an example
type SomeStruct struct {
	SChar byte
	UChar byte
	Char  byte
	SInt  int16
	USInt int16
	Int   int32
	UInt  int32
	Mess  [10]byte `mg:"omit=true"`
	Str   string   `mg:"len=10"`
}

//go:generate methodgen -tmpl=$GOPATH/src/github.com/liamCDI/MethodGen/tmpls/Equal.tmpl -struct=SomeStruct

//go:generate methodgen -tmpl=$GOPATH/src/github.com/liamCDI/MethodGen/tmpls/PrintTag.tmpl -struct=SomeStruct
