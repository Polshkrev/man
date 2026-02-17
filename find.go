package man

import (
	"fmt"
	"strings"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/collections/safe"
)

func normalizeNeedle(needle string) string {
	return strings.ToLower(strings.TrimSpace(needle))
}

// Concurrently seach for a need in an [goserialize.Object] haystack.
func concurrentNameSearch(name string, pages collections.View[Page], resultChannel chan<- Page, errorChannel chan<- *gopolutils.Exception) {
	defer close(resultChannel)
	defer close(errorChannel)
	var i int
	for i = range pages.Collect() {
		var page Page = pages.Collect()[i]
		if normalizeNeedle(name) != normalizeNeedle(page.Name) {
			continue
		}
		resultChannel <- page
		errorChannel <- nil
		return
	}
	resultChannel <- *NewPage("", None, "")
	errorChannel <- gopolutils.NewNamedException(gopolutils.LookupError, fmt.Sprintf("Can not find page with name '%s'.", name))
}

// Concurrently seach for a need in an [goserialize.Object] haystack.
func concurrentSectionSearch(section Section, pages collections.View[Page], resultChannel chan<- collections.View[Page], errorChannel chan<- *gopolutils.Exception) {
	defer close(resultChannel)
	defer close(errorChannel)
	var result safe.Collection[Page] = safe.NewArray[Page]()
	var i int
	for i = range pages.Collect() {
		var page Page = pages.Collect()[i]
		if section != normalizeNeedle(page.Section) {
			continue
		}
		result.Append(page)
	}
	if result.IsEmpty() {
		resultChannel <- nil
		errorChannel <- gopolutils.NewNamedException(gopolutils.LookupError, fmt.Sprintf("Can not find pages with section '%s'.", section))
		return
	}
	resultChannel <- result
	errorChannel <- nil
}

// Find a given name in an [goserialize.Object].
// Returns the contents of the file with the given name.
// If the given title can not be cut from the token, a [gopolutils.ValueError] is returned with a nil data pointer.
func FindByName(entries collections.View[Page], name string) (Page, *gopolutils.Exception) {
	var resultChannel chan Page = make(chan Page, 1)
	var exceptChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception, 1)
	go concurrentNameSearch(name, entries, resultChannel, exceptChannel)
	var result Page = <-resultChannel
	var except *gopolutils.Exception = <-exceptChannel
	return result, except
}

// Find a given name in an [goserialize.Object].
// Returns the contents of the file with the given name.
// If the given title can not be cut from the token, a [gopolutils.ValueError] is returned with a nil data pointer.
func FindBySection(entries collections.View[Page], section Section) (collections.View[Page], *gopolutils.Exception) {
	var resultChannel chan collections.View[Page] = make(chan collections.View[Page], 1)
	var exceptChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception, 1)
	go concurrentSectionSearch(section, entries, resultChannel, exceptChannel)
	var result collections.View[Page] = <-resultChannel
	var except *gopolutils.Exception = <-exceptChannel
	return result, except
}

func FindByNameFromSection(entries collections.View[Page], name string, section Section) (Page, *gopolutils.Exception) {
	var sections collections.View[Page]
	var except *gopolutils.Exception
	sections, except = FindBySection(entries, section)
	if except != nil {
		return *NewPage("", None, ""), except
	}
	return FindByTitle(sections, name)
}
