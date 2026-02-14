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
	"github.com/Polshkrev/goserialize"
)

// Read the files of a given root path concatenated with the given documentation folder and manuals folder.
// Returns a [goserialize.Object] of names mapped to their file content.
// If the absolute path of the file can not be obtained, or the file can not be read, an [gopolutils.IOError] is returned with a nil data pointer.
// If the given title can not be cut from the token, a [gopolutils.ValueError] is returned with a nil data pointer.
// If the key is already in the result object, a [gopolutils.KeyError] is returned with a nil data pointer.
func ReadFiles(root *fayl.Path, documentationFolder, manualsFolder string) (goserialize.Object, *gopolutils.Exception) {
	var documentationPath *fayl.Path = fayl.PathFrom(appendRoot(root, documentationFolder))
	var manualsPath *fayl.Path = findFolder(documentationPath, manualsFolder)
	var names collections.View[string]
	var except *gopolutils.Exception
	names, except = getNames(manualsPath)
	if except != nil {
		return nil, except
	}
	var result goserialize.Object
	result, except = namesToObjects(names, manualsPath)
	if except != nil {
		return nil, except
	}
	return result, nil
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

// Concurrently read a directory from its given path.
func readConcurrent(folder *fayl.Path, resultChannel chan<- []os.DirEntry, errorChannel chan<- error) {
	var result []os.DirEntry
	var readError error
	result, readError = os.ReadDir(folder.ToString())
	resultChannel <- result
	errorChannel <- readError
}

// Append the a child folder to the root of the filesystem.
// This needs to be defined due to a bug in the implementation of [fayl.Path.AppendAs]
// Returns the given child folder to the destination path.
func appendRoot(root *fayl.Path, child string) string {
	return filepath.Join(root.ToString(), string(filepath.Separator), child)
}

// Obtain a [collections.View] from a given slice of [os.DirEntry].
// Returns a [collections.View] containing the names of each of the given [os.DirEntry]
func entriesToView(entries []os.DirEntry) collections.View[string] {
	var result safe.Collection[string] = safe.NewArray[string]()
	var i int
	for i = range entries {
		result.Append(entries[i].Name())
	}
	return result
}

// Obtain the names of each of the files in a given directory.
// Returns a [collections.View] of each of the files' names in the given directory.
// If the directory can not be read, an [gopolutils.IOError] is returned with a nil data pointer.
func getNames(folder *fayl.Path) (collections.View[string], *gopolutils.Exception) {
	var resultChannel chan []os.DirEntry = make(chan []os.DirEntry, 1)
	var errorChannel chan error = make(chan error, 1)
	go readConcurrent(folder, resultChannel, errorChannel)
	var rawResult []os.DirEntry = <-resultChannel
	var readError error = <-errorChannel
	if readError != nil {
		return nil, gopolutils.NewNamedException(gopolutils.IOError, readError.Error())
	}
	return entriesToView(rawResult), nil
}

// Pack a page given its name and content.
// This function writes to the given result pointer.
// If the key is already in the result object, a [gopolutils.KeyError] is returned.
func packPage(result *goserialize.Object, name, content string) *gopolutils.Exception {
	var ok bool
	_, ok = (*result)[name]
	if ok {
		return gopolutils.NewNamedException(gopolutils.KeyError, fmt.Sprintf("Can not pack page '%s'.", name))
	}
	(*result)[name] = content
	return nil
}

// Obtain the content of a file constructed from its given parent folder and name.
// Returns the content of the constructed file.
// If the absolute path of the file can not be obtained, or the file can not be read, an [gopolutils.IOError] is returned with a nil data pointer.
func getContent(parentFolder *fayl.Path, name string) (string, *gopolutils.Exception) {
	var content []byte
	var except *gopolutils.Exception
	content, except = fayl.Read(fayl.PathFrom(appendRoot(parentFolder, name)))
	if except != nil {
		return "", except
	}
	return string(content), nil
}

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

// Convert a given [gopolutils.View] of file names into a [goserialize.Object].
// Returns a [goserialize.Object] of filenames mapped to their content based on a given [collections.View] of names.
// If the absolute path of the file can not be obtained, or the file can not be read, an [gopolutils.IOError] is returned with a nil data pointer.
// If the given title can not be cut from the token, a [gopolutils.ValueError] is returned with a nil data pointer.
// If the key is already in the result object, a [gopolutils.KeyError] is returned with a nil data pointer.
func namesToObjects(names collections.View[string], parentPath *fayl.Path) (goserialize.Object, *gopolutils.Exception) {
	var i int
	var packedPage goserialize.Object = make(goserialize.Object, 0)
	for i = range names.Collect() {
		var name string = names.Collect()[i]
		var content string
		var except *gopolutils.Exception
		content, except = getContent(parentPath, name)
		if except != nil {
			return nil, except
		}
		var title string
		title, except = normalizeTitle(name, "(")
		if except != nil {
			return nil, except
		}
		except = packPage(&packedPage, title, content)
		if except != nil {
			return nil, except
		}
	}
	return packedPage, nil
}
