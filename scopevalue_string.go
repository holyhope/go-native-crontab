// Code generated by "stringer -type=ScopeValue -linecomment"; DO NOT EDIT.

package god

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ScopeUser-0]
	_ = x[ScopeSystem-1]
}

const _ScopeValue_name = "usersystem"

var _ScopeValue_index = [...]uint8{0, 4, 10}

func (i ScopeValue) String() string {
	if i >= ScopeValue(len(_ScopeValue_index)-1) {
		return "ScopeValue(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ScopeValue_name[_ScopeValue_index[i]:_ScopeValue_index[i+1]]
}
