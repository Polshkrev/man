package main

import (
	"os"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/fayl"
)

const (
	documentationFolder  string = "documentation"
	minimumArgumentCount uint8  = 2
	maximumArgumentCount uint8  = 3
)

// Obtain an argument at an index if it is between the given minimum and maximum values.
// Returns the argument between the given minimum and maximum values.
// If the slice of arguments provided have a length less than the given minimum, an [gopolutils.UnderflowError] is returned with an empty string.
// If the slice of arguments provided have a length more than the given maximum, an [gopolutils.OverflowError] is returned with an empty string.
func getArgument(minimum, maximum uint8, args ...string) (string, *gopolutils.Exception) {
	if len(args) < int(minimum) {
		return "", gopolutils.NewNamedException(gopolutils.UnderflowError, "Can not provide less than one argument.")
	} else if len(args) > int(maximum) {
		return "", gopolutils.NewNamedException(gopolutils.UnderflowError, "Can not provide more than one argument.")
	}
	return args[1], nil
}

func main() {
	var name = gopolutils.Must(getArgument(minimumArgumentCount, maximumArgumentCount, os.Args...))
	var files collections.View[*fayl.Path] = gopolutils.Must(GetFiles(documentationFolder))
	var content []byte = gopolutils.Must(FindByTitle(files, name))
}
