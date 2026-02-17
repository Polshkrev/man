package man

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/fayl"
)

// Representation of a linux manual page with its content and metadata.
type Page struct {
	Name    string  `json:"name"`
	Section Section `json:"section"`
	Content string  `json:"content"`
}

// Construct a new [Page] from its given parts.
// Returns a new [Page] constructed from its given parts.
func NewPage(name string, section Section, content string) *Page {
	var page *Page = new(Page)
	page.Name = name
	page.Section = section
	page.Content = content
	return page
}

// Normalize the given name.
// If the given name can not be cut from the token, a [gopolutils.ValueError] is returned with an empty string, else the name cut from after the given token is returned with a nil exception pointer.
func normalizeSection(name, token string) (string, *gopolutils.Exception) {
	var strip string
	var after string
	var found bool
	strip, after, found = strings.Cut(name, token)
	if !found {
		return "", gopolutils.NewNamedException(gopolutils.ValueError, fmt.Sprintf("Can not find token '%s' in '%s'; Before: %s, After %s, Found %t.", token, name, strip, after, found))
	}
	return after, nil
}

// Cut the name of the file from its given [fayl.Path].
// Returns the name of the file cut from its given path.
// If the given path can not be cut, a [gopolutils.ValueError] is returned with an empty string.
func getNameFromPath(file *fayl.Path) (string, *gopolutils.Exception) {
	var after string
	var found bool
	after, found = strings.CutPrefix(file.ToString(), string(filepath.Separator))
	if !found {
		return "", gopolutils.NewNamedException(gopolutils.ValueError, fmt.Sprintf("Can not find token '%s' in '%s'; After %s, Found %t.", string(filepath.Separator), file.ToString(), after, found))
	}
	var strip string
	strip, found = strings.CutSuffix(file.ToString(), "(")
	if !found {
		return "", gopolutils.NewNamedException(gopolutils.ValueError, fmt.Sprintf("Can not find token '%s' in '%s'; Before %s, Found %t.", "(", file.ToString(), strip, found))
	}
	return strip, nil
}

// Obtain the string of the section from the given filename.
// Returns the string of the section cut from the given filename.
// If the given filename can not be cut, a [gopolutils.ValueError] is returned.
func getSection(filename string) (string, *gopolutils.Exception) {
	var initialCut string
	var except *gopolutils.Exception
	initialCut, except = cutNameFromFile(filename, ")")
	if except != nil {
		return "", except
	}
	return normalizeSection(initialCut, "(")
}

// Construct a new [Page] from its given [fayl.Path].
// Returns a new [Page] from its given [fayl.Path].
// If the [Page] properties can not be cut, the constructor panics.
func PageFromFile(file *fayl.Path) *Page {
	return NewPage(gopolutils.Must(getNameFromPath(file)), gopolutils.Must(getSection(file.ToString())), string(gopolutils.Must(fayl.Read(file))))
}
