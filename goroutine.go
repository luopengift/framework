package framework

type goroutine struct {
	name     string
	exec     Goroutiner
	min, max int
}
