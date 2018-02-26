package somestruct

import "testing"

func TestEqual(t *testing.T) {
	str1 := &SomeStruct{1, 1, 1, 1, 1, 1, 1, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, "test"}
	str2 := &SomeStruct{1, 1, 1, 1, 1, 1, 1, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, "test"}
	str3 := &SomeStruct{2, 1, 1, 1, 1, 1, 1, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, "test"}

	if str1.Equals(str2) == false {
		t.Errorf("Str1 and Str2 are equal")
	}

	if str2.Equals(str3) == true {
		t.Errorf("Str1 and Str3 are not equal")
	}
}

func TestString(t *testing.T) {
	str1 := &SomeStruct{1, 1, 1, 1, 1, 1, 1, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, "test"}
	exp := "SChar:1 UChar:1 Char:1 SInt:1 USInt:1 Int:1 UInt:1 Str:test "
	if str1.String() != exp {
		t.Errorf("Expected %s Got %s \n", exp, str1)
	}
}
