package log

type WriterMock struct {
	WriteFunc func(p []byte) (n int, err error)
}

func (w *WriterMock) Write(p []byte) (n int, err error) {
	return w.WriteFunc(p)
}
