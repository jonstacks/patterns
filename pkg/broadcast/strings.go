package broadcast

// Strings creates a new broadcast between the in string channel and the out
// string channels. The last parameter, propagateClose, determines if the output
// channels are closed when the input channel is closed.
func Strings(in <-chan string, out []chan string, propagateClose bool) {
	go func() {
		for str := range in {
			for _, c := range out {
				c <- str
			}
		}

		if propagateClose {
			for _, c := range out {
				close(c)
			}
		}
	}()
}
