package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/collections/safe"
	"github.com/Polshkrev/gopolutils/fayl"
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
func readFiles(read bool, documentationFolder, manualsFolder string, result *collections.View[man.Page]) {
	if !read {
		return
	}
	var root *fayl.Path = getRoot()
	(*result) = gopolutils.Must(man.ReadFiles(root, documentationFolder, manualsFolder))
}

// Write given content to a given path based on a sentinal boolean flag.
func writeFiles(write bool, target *fayl.Path, content collections.View[man.Page]) {
	if !write {
		return
	} else if content == nil {
		return
	}
	var except *gopolutils.Exception = fayl.WriteList(target, content)
	if except != nil {
		panic(except)
	}
}

func appendRoot(root *fayl.Path, child string) string {
	return filepath.Join(root.ToString(), string(filepath.Separator), child)
}

func getTargetFile(name string, fileType fayl.Suffix) string {
	var root *fayl.Path = getRoot()
	var documentationPath *fayl.Path = fayl.PathFrom(appendRoot(root, documentationFolder))
	var manualsPath *fayl.Path = fayl.PathFrom(appendRoot(documentationPath, manualsFolder))
	return fayl.PathFromParts(manualsPath.ToString(), name, fileType).ToString()
}

func find(name string, section man.Section, entries collections.View[man.Page]) man.Page {
	if section == man.None {
		return gopolutils.Must(man.FindByTitle(entries, name))
	}
	return gopolutils.Must(man.FindByNameFromSection(entries, name, section))
}

func main() {
	var write *bool = flag.Bool("write", false, "Write the in-memory cache to a persistant target file.")
	var read *bool = flag.Bool("read", false, "Read files into the in-memory cache")
	var target *string = flag.String("o", getTargetFile("pages", fayl.Json), "Output file to dump the in-memory cache. This will only matter if the 'read' flag is set.")
	var size *bool = flag.Bool("n", false, "Print the total amount of pages.")
	var section *string = flag.String("s", man.None, "Specify the section from which to lookup.")
	flag.Parse()
	var targetPath *fayl.Path = fayl.PathFrom(*target)
	var data collections.View[man.Page] = safe.NewArray[man.Page]()
	readFiles(*read, documentationFolder, manualsFolder, &data)
	writeFiles(*write, targetPath, data)
	data = gopolutils.Must(fayl.ReadList[man.Page](targetPath))
	if *size {
		fmt.Println(len(*data))
		os.Exit(0)
	}
	var name string = gopolutils.Must(getArgument(0, minimumArgumentCount, maximumArgumentCount, flag.Args()...))
	var page man.Page = find(name, man.Section(*section), data)
	fmt.Print(page.Content)
}
