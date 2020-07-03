package logic

type Sender interface {
	Send(text string)
	SendPhoto(name string, b []byte) error
}
