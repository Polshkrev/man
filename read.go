package man

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/collections/safe"
	"github.com/Polshkrev/gopolutils/fayl"
)

// Obtain all the files in a given folder.
// Returns a [collections.View] of [fayl.Path] containing all the files of the given directory.
// If the given directory is empty, an [gopolutils.IOError] is returned with a nil data pointer.
func GetFiles(folder string) (collections.View[*fayl.Path], *gopolutils.Exception) {
	var result safe.Collection[*fayl.Path] = safe.NewArray[*fayl.Path]()
	var root *fayl.Path = getRoot()
	var directory *fayl.Path = findFolder(root, folder)
	var folders collections.View[*fayl.Path]
	var except *gopolutils.Exception
	folders, except = directoryAsView(directory)
	if except != nil {
		return nil, except
	}
	var i int
	for i = range folders.Collect() {
		var file *fayl.Path = folders.Collect()[i]
		var view collections.View[*fayl.Path]
		view, except = directoryAsView(file)
		if except != nil {
			return nil, except
		}
		result.Extend(view)
	}
	if result.IsEmpty() {
		return nil, gopolutils.NewException(fmt.Sprintf("Directory '%s' seems to be empty.", folder))
	}
	return result, nil
}

// Obtain the root of the filesystem.
// Returns the root of the current directory.
func getRoot() *fayl.Path {
	var current *fayl.Path = fayl.NewPath()
	return gopolutils.Must(current.Root())
}

// Append the a child folder to the root of the filesystem.
// This needs to be defined due to a bug in the implementation of [fayl.Path.AppendAs]
// Returns the given child folder to the destination path.
func appendRoot(destination *fayl.Path, child string) *fayl.Path {
	return fayl.PathFrom(filepath.Join(destination.ToString(), string(filepath.Separator), child))
}

// Convert a collection of [os.DirEntry] into a collection of [fayl.Path].
// Returns a collection of [fayl.Path] containing the names of the given [os.DirEntry]
func entriesAsPaths(root *fayl.Path, entries []os.DirEntry) collections.View[*fayl.Path] {
	var result safe.Collection[*fayl.Path] = safe.NewArray[*fayl.Path]()
	var i int
	for i = range entries {
		var file os.DirEntry = entries[i]
		var absolute *fayl.Path = gopolutils.Must(root.Absolute())
		result.Append(fayl.PathFrom(filepath.Join(absolute.ToString(), file.Name())))
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
	var full *fayl.Path = appendRoot(root, child)
	var except *gopolutils.Exception = makeDirectoryIfNotExists(full)
	if except != nil {
		panic(except)
	}
	return full
}

// Concurrently read a directory of a given path.
func readDirectory(path *fayl.Path, resultChannel chan<- []os.DirEntry, errorChannel chan<- error) {
	var files []os.DirEntry
	var directoryError error
	files, directoryError = os.ReadDir(path.ToString())
	resultChannel <- files
	errorChannel <- directoryError
	defer close(resultChannel)
	defer close(errorChannel)
}

// Obtain the files in a directory of a given path.
// Returns a [collections.View] of [fayl.Path] from the given path.
// If the directory can not be read, an [gopolutils.IOError] is returned with a nil data pointer.
func directoryAsView(path *fayl.Path) (collections.View[*fayl.Path], *gopolutils.Exception) {
	var fileChannel chan []os.DirEntry = make(chan []os.DirEntry, 1)
	var errorChannel chan error = make(chan error, 1)
	go readDirectory(path, fileChannel, errorChannel)
	var files []os.DirEntry = <-fileChannel
	var except error = <-errorChannel
	if except != nil {
		return nil, gopolutils.NewNamedException(gopolutils.IOError, except.Error())
	}
	var result collections.View[*fayl.Path] = entriesAsPaths(path, files)
	if result.IsEmpty() {
		return nil, gopolutils.NewNamedException(gopolutils.IOError, fmt.Sprintf("Directory '%s' seems to be empty.", path.ToString()))
	}
	return result, nil
}
