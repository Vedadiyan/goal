package util

func GuardAgainstClosedChan[T any](c <-chan T) <-chan T {
	returnChan := make(chan T)
	value, ok := <-c
	if ok {
		go func() {
			returnChan <- value
		}()
	}
	return returnChan
}
