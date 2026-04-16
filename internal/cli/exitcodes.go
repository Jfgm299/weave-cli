package cli

const (
	ExitOK                = 0
	ExitInvalidConfig     = 2
	ExitMissingDependency = 3
	ExitRuntimeError      = 4
	ExitDoctorIssues      = 5
)

func exitCodeForError(err error) int {
	if err == nil {
		return ExitOK
	}

	if isInvalidConfigError(err) {
		return ExitInvalidConfig
	}

	if isMissingDependencyError(err) {
		return ExitMissingDependency
	}

	return ExitRuntimeError
}
