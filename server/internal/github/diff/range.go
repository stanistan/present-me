package diff

type RangeFrom [2]*int

func (r *RangeFrom) extract() (int, bool) {
	if r[0] != nil {
		return *r[0], true
	} else if r[1] != nil {
		return *r[1], true
	}
	return 0, false
}
