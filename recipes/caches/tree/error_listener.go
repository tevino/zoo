package tree

// ErrorListener is a handler of error.
type ErrorListener struct {
	fn func(error)
}

// Handle calls the function with e.
func (l *ErrorListener) Handle(e error) {
	l.fn(e)
}

// NewErrorListener creates ErrorListener from fn.
func NewErrorListener(fn func(error)) *ErrorListener {
	return &ErrorListener{fn}
}
