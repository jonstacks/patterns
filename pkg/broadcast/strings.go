package broadcast

// Strings creates a new broadcast between the in string channel and the out
// string channels. The last parameter, propagateClose, determines if the output
// channels are closed when the input channel is closed.
func Strings(in <-chan string, out []chan string, propagateClose bool) {
	StringsWithOptions(in, out, StringOptions{
		NonBlocking:    false,
		PropagateClose: propagateClose,
	})
}

// StringOptions control the behavior of how strings are broadcast by the
// StringsWithOptions function.
type StringOptions struct {
	// If true, messages will be discarded rather than blocking the channel
	NonBlocking bool
	// If propagateClose is true, then a close message on the input
	// will be broadcast to all output channels
	PropagateClose bool
}

// StringsWithOptions creates a new broadcast between the in string channel and
// the out string channels. If the nonBlocking option is true, for channels which
// are already at capacity, the message will be dropped rather than block. If the
// propagateClose option is true, then closing the input channel will close all
// channels that it is broadcasting to.
func StringsWithOptions(in <-chan string, out []chan string, opts StringOptions) {
	go func() {
		for str := range in {
			for _, c := range out {
				if opts.NonBlocking {
					select {
					case c <- str:
					default:
					}
				} else {
					c <- str
				}
			}
		}

		if opts.PropagateClose {
			for _, c := range out {
				close(c)
			}
		}
	}()
}
