package broadcast

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertStringChannelsClosed(t *testing.T, channels []chan string) {
	for _, c := range channels {
		_, ok := <-c
		assert.Equal(t, false, ok, "Expected channel to be closed but it wasn't")
	}
}

func TestStringsPropagatesClose(t *testing.T) {
	cin := make(chan string)
	cout1 := make(chan string, 5)
	cout2 := make(chan string, 5)
	couts := []chan string{cout1, cout2}

	Strings(cin, couts, true)

	close(cin)
	assertStringChannelsClosed(t, couts)
}

func TestStringsAreBroadcasted(t *testing.T) {
	cin := make(chan string)
	cout1 := make(chan string, 5)
	cout2 := make(chan string, 5)
	couts := []chan string{cout1, cout2}

	Strings(cin, couts, true)

	defer func() {
		close(cin)
		assertStringChannelsClosed(t, couts)
	}()

	cin <- "Hello"
	assert.Equal(t, "Hello", <-cout1)
	assert.Equal(t, "Hello", <-cout2)
}

func TestStringsWhenBlocked(t *testing.T) {
	cin := make(chan string)
	cout1 := make(chan string)
	cout2 := make(chan string)
	couts := []chan string{cout1, cout2}

	Strings(cin, couts, true)
	defer func() {
		close(cin)
		assertStringChannelsClosed(t, couts)
	}()

	cin <- "Hello"

	select {
	case cin <- "This should block":
		t.Error("Expected writing to input channel to block, but it didn't")
	default:
	}

	// Now read from output channels
	assert.Equal(t, "Hello", <-cout1)
	assert.Equal(t, "Hello", <-cout2)

	cin <- "Hello2"
	assert.Equal(t, "Hello2", <-cout1)
	assert.Equal(t, "Hello2", <-cout2)
}

func TestStringsWithOptionsNonBlocking(t *testing.T) {
	cin := make(chan string)
	cout1 := make(chan string)
	cout2 := make(chan string)
	couts := []chan string{cout1, cout2}

	StringsWithOptions(cin, couts, StringOptions{
		nonBlocking:    true,
		propagateClose: true,
	})
	defer func() {
		close(cin)
		assertStringChannelsClosed(t, couts)
	}()

	cin <- "This should be dropped since we don't have any readers"
}

func TestStringsWithOptionsNonBlockingWithBufferedChannels(t *testing.T) {
	cin := make(chan string)
	cout1 := make(chan string, 3)
	cout2 := make(chan string, 5)
	couts := []chan string{cout1, cout2}

	StringsWithOptions(cin, couts, StringOptions{
		nonBlocking:    true,
		propagateClose: true,
	})
	defer func() {
		close(cin)
		assertStringChannelsClosed(t, couts)
	}()

	// Writing a number of messages much greater than the size of the
	// output channels should not block at all.
	for i := 1; i <= 100; i++ {
		cin <- fmt.Sprintf("Hello %d", i)
	}

	// Drain messages
CLEAN:
	for {
		select {
		case <-cout1:
		case <-cout2:
		case <-time.After(1 * time.Second):
			break CLEAN
		}
	}
}
