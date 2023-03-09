package util

func SimpleChanGuard[T any](source <-chan T) <-chan T {
	returnChan := make(chan T)
	defer close(returnChan)
	value, ok := <-source
	if ok {
		go func() {
			returnChan <- value
		}()
	}
	return returnChan
}

func ReferencedChanGuard[T any](source <-chan T, ref *chan T) <-chan T {
	returnChan := make(chan T)
	*ref = returnChan
	value, ok := <-source
	if ok {
		go func() {
			returnChan <- value
		}()
	}
	return returnChan
}
