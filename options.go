package god

import "time"

type DarwinOptions interface {
	darwinLimitLoadToSessionTypeOption
}
type Options interface {
	DarwinOptions

	nameOption
	scopeOption
	stateOption
	programOption
	intervalOption
	argumentsOption
	runAtLoadOption
	userOwnerOption
	groupOwnerOption
	descriptionOption
	errorOutputOption
	standardOutputOption
	workingDirectoryOption
	startLimitIntervalOption
	environmentVariablesOption
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
	scope                *Scope
	state                *State
	program              *string
	interval             *time.Duration
	arguments            *[]string
	runAtLoad            *bool
	userOwner            *int
	groupOwner           *int
	description          *string
	errorOutput          *string
	standardOutput       *string
	workingDirectory     *string
	startLimitInterval   *time.Duration
	environmentVariables *map[string]string
}

func (opts *options) copy() *options {
	return &options{
		name:                 opts.name,
		scope:                opts.scope,
		state:                opts.state,
		darwin:               *opts.darwin.copy(),
		program:              opts.program,
		interval:             opts.interval,
		arguments:            opts.arguments,
		runAtLoad:            opts.runAtLoad,
		userOwner:            opts.userOwner,
		groupOwner:           opts.groupOwner,
		description:          opts.description,
		errorOutput:          opts.errorOutput,
		standardOutput:       opts.standardOutput,
		workingDirectory:     opts.workingDirectory,
		startLimitInterval:   opts.startLimitInterval,
		environmentVariables: opts.environmentVariables,
	}
}

func (opts *darwinOptions) copy() *darwinOptions {
	return &darwinOptions{
		limitLoadToSessionTypeOption: opts.limitLoadToSessionTypeOption,
	}
}
