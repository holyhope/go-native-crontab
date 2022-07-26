package god

import "time"

type DarwinOptions interface {
	darwinLimitLoadToSessionTypeOption
}
type Options interface {
	DarwinOptions

	nameOption
	stateOption
	descriptionOption
	programOption
	argumentsOption
	runAtLoadOption
	environmentVariablesOption
	scopeOption
	intervalOption
	userOwnerOption
	groupOwnerOption
	workingDirectoryOption
	standardOutputOption
	errorOutputOption
	startLimitIntervalOption
}

func Opts() Options {
	return &options{}
}

type darwinOptions struct {
	limitLoadToSessionTypeOption *DarwinLimitLoadToSessionType
}

type options struct {
	darwin darwinOptions

	name                 *string
	state                *State
	description          *string
	program              *string
	arguments            *[]string
	runAtLoad            *bool
	environmentVariables *map[string]string
	scope                *Scope
	interval             *time.Duration
	userOwner            *int
	groupOwner           *int
	workingDirectory     *string
	standardOutput       *string
	errorOutput          *string
	startLimitInterval   *time.Duration
}

func (opts *options) copy() *options {
	return &options{
		darwin:               *opts.darwin.copy(),
		name:                 opts.name,
		state:                opts.state,
		description:          opts.description,
		program:              opts.program,
		arguments:            opts.arguments,
		runAtLoad:            opts.runAtLoad,
		environmentVariables: opts.environmentVariables,
		scope:                opts.scope,
		interval:             opts.interval,
		userOwner:            opts.userOwner,
		groupOwner:           opts.groupOwner,
		workingDirectory:     opts.workingDirectory,
		standardOutput:       opts.standardOutput,
		errorOutput:          opts.errorOutput,
	}
}

func (opts *darwinOptions) copy() *darwinOptions {
	return &darwinOptions{
		limitLoadToSessionTypeOption: opts.limitLoadToSessionTypeOption,
	}
}
