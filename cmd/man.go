package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/fayl"
	"github.com/Polshkrev/goserialize"
	"github.com/Polshkrev/man"
)

const (
	documentationFolder  string = "documentation"
	targetFile           string = "pages"
	manualsFolder        string = "man"
	minimumArgumentCount uint8  = 1
	maximumArgumentCount uint8  = 2
)

func getRoot() *fayl.Path {
	var current *fayl.Path = fayl.NewPath()
	return gopolutils.Must(current.Root())
}

// Obtain an argument at an index if it is between the given minimum and maximum values.
// Returns the argument between the given minimum and maximum values.
// If the slice of arguments provided have a length less than the given minimum, an [gopolutils.UnderflowError] is returned with an empty string.
// If the slice of arguments provided have a length more than the given maximum, an [gopolutils.OverflowError] is returned with an empty string.
func getArgument(index, minimum, maximum uint8, args ...string) (string, *gopolutils.Exception) {
	if len(args) < int(minimum) {
		return "", gopolutils.NewNamedException(gopolutils.UnderflowError, "Can not provide less than one argument.")
	} else if len(args) > int(maximum) {
		return "", gopolutils.NewNamedException(gopolutils.OverflowError, "Can not provide more than one argument.")
	}
	return args[index], nil
}

// Read content into a result parametre from a given path constructed from its root and two isometric children based on a sentinal boolean flag.
func readFiles(read bool, documentationFolder, manualsFolder string, result *goserialize.Object) {
	if !read {
		return
	}
	var root *fayl.Path = getRoot()
	(*result) = gopolutils.Must(man.ReadFiles(root, documentationFolder, manualsFolder))
}

// Write given content to a given path based on a sentinal boolean flag.
func writeFiles(write bool, target *fayl.Path, content *goserialize.Object) {
	if !write {
		return
	} else if content == nil {
		return
	}
	var except *gopolutils.Exception = fayl.WriteObject(target, content)
	if except != nil {
		panic(except)
	}
}

func appendRoot(root *fayl.Path, child string) string {
	return filepath.Join(root.ToString(), string(filepath.Separator), child)
}

// func If[TrueValue, FalseValue any](condition bool, trueValue TrueValue, falseValue FalseValue) {
// 	if !condition {
// 		return falseValue
// 	}
// 	return trueValue
// }

func getTargetFile(name string, fileType fayl.Suffix) string {
	var root *fayl.Path = getRoot()
	var documentationPath *fayl.Path = fayl.PathFrom(appendRoot(root, documentationFolder))
	var manualsPath *fayl.Path = fayl.PathFrom(appendRoot(documentationPath, manualsFolder))
	return appendRoot(manualsPath, fayl.PathFromParts(manualsPath, name, fileType).ToString())
}

func main() {
	var write *bool = flag.Bool("write", false, "Write the in-memory cache to a persistant target file.")
	var read *bool = flag.Bool("read", false, "Read files into the in-memory cache")
	var target *string = flag.String("o", getTargetFile("pages", fayl.Json), "Output file to dump the in-memory cache. This will only matter if the 'read' flag is set.")
	flag.Parse()
	var targetPath *fayl.Path = fayl.PathFrom(*target)
	var data goserialize.Object
	readFiles(*read, documentationFolder, manualsFolder, &data)
	writeFiles(*write, targetPath, &data)
	var results *goserialize.Object = gopolutils.Must(fayl.ReadObject[goserialize.Object](targetPath))
	var name string = gopolutils.Must(getArgument(0, minimumArgumentCount, maximumArgumentCount, flag.Args()...))
	var content string = gopolutils.Must(man.FindByTitle(*results, name))
	fmt.Print(content)
}
