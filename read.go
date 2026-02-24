package man

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/collections/safe"
	"github.com/Polshkrev/gopolutils/fayl"
	"github.com/Polshkrev/goserialize"
)

// Cut the base from the given file name.
// If the given base can not be cut from the token, a [gopolutils.ValueError] is returned with an empty string, else the name of the title is returned with a nil exception pointer.
func cutNameFromFile(filename, token string) (string, *gopolutils.Exception) {
	var lower string = strings.ToLower(filename)
	var strip string
	var after string
	var found bool
	strip, after, found = strings.Cut(lower, token)
	if !found {
		return "", gopolutils.NewNamedException(gopolutils.ValueError, fmt.Sprintf("Can not find token '%s' in '%s'; Before: %s, After %s, Found %t.", token, lower, strip, after, found))
	}
	return strip, nil
}

// Append the a child folder to the root of the filesystem.
// This needs to be defined due to a bug in the implementation of [fayl.Path.AppendAs]
// Returns the given child folder to the destination path.
func appendRoot(root *fayl.Path, child string) string {
	return filepath.Join(root.ToString(), string(filepath.Separator), child)
}

// Convert a given [collections.View] of [fayl.Path] into a [collections.View] of [Page].
// Returns a [collections.View] of [Page] based on a given [collections.View] of [fayl.Path].
func pathsToPages(entries collections.View[*fayl.Entry]) collections.View[Page] {
	var result safe.Collection[Page] = safe.NewArray[Page]()
	var i int
	for i = range entries.Collect() {
		var entry *fayl.Entry = entries.Collect()[i]
		if entry.Is(fayl.DirectoryType) {
			continue
		}
		var path *fayl.Path = entry.Path()
		result.Append(*PageFromFile(path))
	}
	return result
}

// Read the files of a given root path concatenated with the given documentation folder and manuals folder.
// Returns a [collections.View] of [Page] based on a [fayl.Path] constructed from its given parts.
func ReadFiles(root *fayl.Path, documentationFolder, manualsFolder string) collections.View[Page] {
	var documentationPath *fayl.Path = fayl.PathFrom(appendRoot(root, documentationFolder))
	var manualsPath *fayl.Path = fayl.PathFrom(appendRoot(documentationPath, manualsFolder))
	var directory *fayl.Directory = fayl.NewDirectory(manualsPath)
	var except *gopolutils.Exception = directory.Read()
	if except != nil {
		panic(except)
	}
	return pathsToPages(goserialize.SliceToView(directory.Collect()))
}
