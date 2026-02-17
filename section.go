package man

import "github.com/Polshkrev/gopolutils"

// Representation of a section within the linux manual.
type Section = gopolutils.StringEnum

const (
	None                   Section = ""
	Header                 Section = "0p"     // Describes header files within the POSIX standard.
	Command                Section = "1"      // Describes user commands.
	PosixCommand           Section = "1p"     // Describes user commands within the POSIX standard.
	SystemCall             Section = "2"      // Describes system calls.
	SystemType             Section = "2type"  // Describes structures that are used with system calls.
	SystemConstant         Section = "2const" // Describes constants that are used with system calls.
	LibraryCall            Section = "3"      // Describes functions and subroutines within the standard library.
	PosixLibraryCall       Section = "3p"     // Describes POSIX-compliant functions and subroutines within the standard library.
	ExtendedLibraryCall    Section = "3x"     // Describes specialized, extended, or non-standard library functions.
	LibraryConstant        Section = "3const" // Describes constants, macros, and defined types within the standard library.
	LibraryType            Section = "3type"  // Describes structures that are used with standard library.
	LibraryHeader          Section = "3head"  // Describes specific header files within the standard library.
	SpecialFiles           Section = "4"      // Describes special files and device drivers.
	FnC                    Section = "5"      // Describes file formats, conventions, and configuration files. Stands for "Formats and Conventions".
	Miscellaneous          Section = "6"      // Describes games, jokes, and amusement programmes.
	AdministrationCommands Section = "7"      // Describes overviews, conventions, protocols, character sets, and miscellaneous topics.
	KernalRoutines         Section = "8"      // Describes system administration and maintenance commands, typically used by the root user or system administrators.
)
