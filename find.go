package man

import (
	"fmt"
	"strings"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/fayl"
)

// Normalize the given title.
// If the given title can not be cut from the token, a [gopolutils.ValueError] is returned with an empty string, else the name of the title is returned with a nil exception pointer.
func normalizeTitle(title, token string) (string, *gopolutils.Exception) {
	var lower string = strings.ToLower(title)
	var strip string
	var after string
	var found bool
	strip, after, found = strings.Cut(lower, token)
	if !found {
		return "", gopolutils.NewNamedException(gopolutils.ValueError, fmt.Sprintf("Can not find token '%s' in '%s'; Before: %s, After %s, Found %t.", token, lower, strip, after, found))
	}
	return strip, nil
}

// Concurrently search a view of [fayl.Path] based on a given name.
func concurrentSearch(files collections.View[*fayl.Path], name string, resultChannel chan<- []byte, errorChannel chan<- *gopolutils.Exception) {
	var i int
	for i = range files.Collect() {
		var file *fayl.Path = files.Collect()[i]
		var title string
		var normalizeException *gopolutils.Exception
		title, normalizeException = normalizeTitle(file.ToString(), "(")
		if normalizeException != nil {
			resultChannel <- nil
			errorChannel <- normalizeException
			defer close(resultChannel)
			defer close(errorChannel)
			return
		}
		if !strings.Contains(title, strings.ToLower(name)) {
			continue
		}
		var result []byte
		var except *gopolutils.Exception
		result, except = fayl.Read(file)
		resultChannel <- result
		errorChannel <- except
	}
	resultChannel <- nil
	errorChannel <- gopolutils.NewNamedException(gopolutils.ValueError, fmt.Sprintf("File '%s' can not be found.", name))
	defer close(resultChannel)
	defer close(errorChannel)
}

// Find a given name in a collection of [fayl.Path].
// Returns a slice of bytes as the contents of the file with the given name.
// If the given title can not be cut from the token, a [gopolutils.ValueError] is returned with a nil data pointer.
func FindByTitle(files collections.View[*fayl.Path], name string) ([]byte, *gopolutils.Exception) {
	var resultChannel chan []byte = make(chan []byte, 1)
	var exceptChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception, 1)
	go concurrentSearch(files, name, resultChannel, exceptChannel)
	var result []byte = <-resultChannel
	var except *gopolutils.Exception = <-exceptChannel
	return result, except
}
