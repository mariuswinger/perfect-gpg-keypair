package utils

type UserInterrupt struct{}

func (e *UserInterrupt) Error() string {
	return "Exiting ..."
}
