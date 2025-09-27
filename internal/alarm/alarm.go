package alarm

const maxRetryCount = 3

type Alarmer interface {
	Alarm(message string)
}
