package vos

type CurrentStateFingerprints struct {
	code        Fingerprint
	instruction Fingerprint
	environment Environment
	vars        Fingerprint
}

func NewCurrentStateFingerprints(
	code, instruction, vars Fingerprint,
	environment Environment,
) CurrentStateFingerprints {
	return CurrentStateFingerprints{
		code:        code,
		instruction: instruction,
		vars:        vars,
		environment: environment,
	}
}

func (c CurrentStateFingerprints) Code() Fingerprint {
	return c.code
}

func (c CurrentStateFingerprints) Instruction() Fingerprint {
	return c.instruction
}

func (c CurrentStateFingerprints) Environment() Environment {
	return c.environment
}

func (c CurrentStateFingerprints) Vars() Fingerprint {
	return c.vars
}
