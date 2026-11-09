package seatbelt

// BuildProfile compiles the policy into a complete SBPL sandbox profile.
func BuildProfile(policy Policy) (*Profile, error) {

	builder := newProfileBuilder()

	builder.addImport("bsd.sb")
	builder.addRule("(allow process-fork)")
	builder.addRule("(allow signal)")
	builder.addRule("(allow sysctl-read)")
	builder.addRule("(allow mach-lookup (global-name \"com.apple.trustd.agent\") (global-name \"com.apple.SystemConfiguration.configd\"))")
	builder.addRule("(allow ipc-posix-sem)")

	steps := []func() error{
		func() error {
			return builder.addFileRules(
				policy.ReadWrite,
				"file-read* file-write* file-read-metadata",
				[]string{
					"com.apple.app-sandbox.read",
					"com.apple.app-sandbox.read-write",
				},
			)
		},
		func() error {
			return builder.addProcessExecRules(policy.ReadWrite)
		},
		func() error {
			return builder.addFileRules(
				policy.ReadOnly,
				"file-read* file-read-metadata",
				[]string{
					"com.apple.app-sandbox.read",
				},
			)
		},
		func() error {
			return builder.addProcessExecRules(policy.ReadOnly)
		},
		func() error {
			return builder.addFileRules(
				policy.Devices,
				"file-read* file-write* file-read-metadata",
				[]string{},
			)
		},
	}

	for _, step := range steps {

		err := step()

		if err == nil {
			continue
		} else {
			return nil, err
		}

	}

	return builder.build()

}
