package somestruct

import "fmt"

func (a *SomeStruct) String() string {
	//SomeStruct
	out := ""

	out += "SChar:" + fmt.Sprintf("%d", a.SChar) + " "

	out += "UChar:" + fmt.Sprintf("%d", a.UChar) + " "

	out += "Char:" + fmt.Sprintf("%d", a.Char) + " "

	out += "SInt:" + fmt.Sprintf("%d", a.SInt) + " "

	out += "USInt:" + fmt.Sprintf("%d", a.USInt) + " "

	out += "Int:" + fmt.Sprintf("%d", a.Int) + " "

	out += "UInt:" + fmt.Sprintf("%d", a.UInt) + " "

	out += "Str:" + a.Str + " "

	return out
}
