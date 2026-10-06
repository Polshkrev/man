package man

import "github.com/Polshkrev/gopolutils"

// Representation of a section within the linux manual.
//
// Deprecated: Due to a move to c++, this will be deleted.
type Section gopolutils.StringEnum

const (
	// Describes a catch-all section.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	None Section = ""
	// Describes header files within the POSIX standard.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	Header Section = "0p"
	// Describes user commands.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	Command Section = "1"
	// Describes user commands within the POSIX standard.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	PosixCommand Section = "1p"
	// Describes system calls.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	SystemCall Section = "2"
	// Describes structures that are used with system calls.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	SystemType Section = "2type"
	// Describes constants that are used with system calls.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	SystemConstant Section = "2const"
	// Describes functions and subroutines within the standard library.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	LibraryCall Section = "3"
	// Describes POSIX-compliant functions and subroutines within the standard library.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	PosixLibraryCall Section = "3p"
	// Describes specialized, extended, or non-standard library functions.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	ExtendedLibraryCall Section = "3x"
	// Describes constants, macros, and defined types within the standard library.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	LibraryConstant Section = "3const"
	// Describes structures that are used with standard library.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	LibraryType Section = "3type"
	// Describes specific header files within the standard library.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	LibraryHeader Section = "3head"
	// Describes special files and device drivers.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	SpecialFiles Section = "4"
	// Describes file formats, conventions, and configuration files. Stands for "Formats and Conventions".
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	FnC Section = "5"
	// Describes games, jokes, and amusement programmes.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	Miscellaneous Section = "6"
	// Describes overviews, conventions, protocols, character sets, and miscellaneous topics.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	AdministrationCommands Section = "7"
	// Describes system administration and maintenance commands, typically used by the root user or system administrators.
	//
	// Deprecated: Due to a move to c++, this will be deleted.
	KernalRoutines Section = "8"
)

func (section Section) String() string {
	return string(section)
}
