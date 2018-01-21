package fakes

type Leftovers struct {
	DeleteCall struct {
		CallCount int
		Receives  struct {
			Filter string
		}
		Returns struct {
			Error error
		}
	}
}

func (l *Leftovers) Delete(filter string) error {
	l.DeleteCall.CallCount++
	l.DeleteCall.Receives.Filter = filter

	return l.DeleteCall.Returns.Error
}
