package man

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/collections/safe"
	"github.com/Polshkrev/gopolutils/fayl"
)

// Concurrently obtain the entries of a given folder.
func getConcurrentEntries(folder string, resultChannel chan<- []os.DirEntry, errorChannel chan<- error) {
	defer close(resultChannel)
	defer close(errorChannel)
	var entries []os.DirEntry
	var readError error
	entries, readError = os.ReadDir(folder)
	errorChannel <- readError
	resultChannel <- entries
}

// Obtain the entries of a given folder.
// Returns a slice of [os.DirEntry].
// If the entries can not be obtained, an [gopolutils.IOError] is returned with a nil data pointer.
func getEntries(folder string) ([]os.DirEntry, *gopolutils.Exception) {
	var resultChannel chan []os.DirEntry = make(chan []os.DirEntry, 1)
	var errorChannel chan error = make(chan error, 1)
	go getConcurrentEntries(folder, resultChannel, errorChannel)
	var result []os.DirEntry = <-resultChannel
	var except error = <-errorChannel
	if except != nil {
		return nil, gopolutils.NewNamedException(gopolutils.IOError, fmt.Sprintf("Can not read directory '%s': %s", folder, except.Error()))
	}
	return result, nil
}

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

// Obtain the paths of a given folder path.
// Returns a [collections.View] of [fayl.Path] based on the files contained within the given folder.
func getPaths(folder *fayl.Path) collections.View[*fayl.Path] {
	var result safe.Collection[*fayl.Path] = safe.NewArray[*fayl.Path]()
	var entries []os.DirEntry = gopolutils.Must(getEntries(folder.ToString()))
	var i int
	for i = range entries {
		var entry os.DirEntry = entries[i]
		result.Append(fayl.PathFrom(appendRoot(folder, entry.Name())))
	}
	return result
}

// Convert a given [collections.View] of [fayl.Path] into a [collections.View] of [Page].
// Returns a [collections.View] of [Page] based on a given [collections.View] of [fayl.Path].
func pathsToPages(paths collections.View[*fayl.Path]) collections.View[Page] {
	var result safe.Collection[Page] = safe.NewArray[Page]()
	var i int
	for i = range paths.Collect() {
		var path *fayl.Path = paths.Collect()[i]
		result.Append(*PageFromFile(path))
	}
	return result
}

// Make a directory if it doesn't exist on the filesystem.
// If the directory can not be created, an [gopolutils.IOError] is returned, else nil is returned.
func makeDirectoryIfNotExists(path *fayl.Path) *gopolutils.Exception {
	if path.Exists() {
		return nil
	}
	var directoryError error = os.MkdirAll(path.ToString(), os.ModePerm)
	if directoryError != nil {
		return gopolutils.NewNamedException(gopolutils.IOError, directoryError.Error())
	}
	return nil
}

// Find a folder and, if it doesn't exist, create it.
// Returns a full appended path of the given root and child.
func findFolder(root *fayl.Path, child string) *fayl.Path {
	var full string = appendRoot(root, child)
	var fullPath *fayl.Path = fayl.PathFrom(full)
	var except *gopolutils.Exception = makeDirectoryIfNotExists(fullPath)
	if except != nil {
		panic(except)
	}
	return fullPath
}

// Read the files of a given root path concatenated with the given documentation folder and manuals folder.
// Returns a [collections.View] of [Page] based on a [fayl.Path] constructed from its given parts.
func ReadFiles(root *fayl.Path, documentationFolder, manualsFolder string) collections.View[Page] {
	var documentationPath *fayl.Path = fayl.PathFrom(appendRoot(root, documentationFolder))
	var manualsPath *fayl.Path = findFolder(documentationPath, manualsFolder)
	var paths collections.View[*fayl.Path] = getPaths(manualsPath)
	var pages collections.View[Page] = pathsToPages(paths)
	return pages
}
