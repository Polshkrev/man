package man

import (
	"fmt"
	"strings"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/goserialize"
)

func normalizeNeedle(needle string) string {
	return strings.ToLower(strings.TrimSpace(needle))
}

// Concurrently seach for a need in an [goserialize.Object] haystack.
func concurrentSearch(needle string, haystack goserialize.Object, resultChannel chan<- string, errorChannel chan<- *gopolutils.Exception) {
	defer close(resultChannel)
	defer close(errorChannel)
	var value any
	var ok bool
	value, ok = haystack[normalizeNeedle(needle)]
	if !ok {
		resultChannel <- ""
		errorChannel <- gopolutils.NewNamedException(gopolutils.KeyError, fmt.Sprintf("Can not find '%s' in mapping.", needle))
		return
	}
	resultChannel <- value.(string)
	errorChannel <- nil
}

// Find a given name in an [goserialize.Object].
// Returns the contents of the file with the given name.
// If the given title can not be cut from the token, a [gopolutils.ValueError] is returned with a nil data pointer.
func FindByTitle(entries goserialize.Object, name string) (string, *gopolutils.Exception) {
	var resultChannel chan string = make(chan string, 1)
	var exceptChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception, 1)
	go concurrentSearch(name, entries, resultChannel, exceptChannel)
	var result string = <-resultChannel
	var except *gopolutils.Exception = <-exceptChannel
	return result, except
}
