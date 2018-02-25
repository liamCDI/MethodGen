package somestruct

type SomeStruct struct {
	SChar byte
	UChar byte
	Char  byte
	SInt  int16
	USInt int16
	Int   int32
	UInt  int32
	Mess  [10]byte
	Str   string `shp:"len=10"`
}

//go:generate methodgen -tmpl=$GOPATH/src/github.com/liamCDI/MethodGen/test/Equal.tmpl -struct=SomeStruct
