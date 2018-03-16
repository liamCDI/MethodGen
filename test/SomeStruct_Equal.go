package somestruct

func (a *SomeStruct) Equals(b *SomeStruct) bool {
	if a.SChar != b.SChar {
		return false
	}
	if a.UChar != b.UChar {
		return false
	}
	if a.Char != b.Char {
		return false
	}
	if a.SInt != b.SInt {
		return false
	}
	if a.USInt != b.USInt {
		return false
	}
	if a.Int != b.Int {
		return false
	}
	if a.UInt != b.UInt {
		return false
	}
	if a.Mess != b.Mess {
		return false
	}
	if a.Str != b.Str {
		return false
	}

	return true
}
