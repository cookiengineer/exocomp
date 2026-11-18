package seatbelt

import "fmt"

// BuildProfile compiles the policy into a complete SBPL sandbox profile.
func BuildProfile(policy Policy) (*Profile, error) {

	builder := newProfileBuilder()

	builder.addImport("bsd.sb")
	builder.addRule("(allow process-fork)")
	builder.addRule("(allow signal)")
	builder.addRule("(allow sysctl-read)")
	builder.addRule("(allow mach-lookup (global-name \"com.apple.trustd.agent\") (global-name \"com.apple.SystemConfiguration.configd\"))")
	builder.addRule("(allow ipc-posix-sem)")

	// NOTE: A network.Connect grant is all-or-nothing at the OS layer. SBPL
	// can filter by port but its host token only accepts "*" or "localhost",
	// so a granted host cannot be matched. Outbound IP traffic is allowed,
	// while unix-socket connects (docker.sock, SSH agent, ...) stay denied.
	// The resolver socket is the one unix socket that scanners need for DNS.
	if policy.Network == true {

		builder.addRule("(allow system-socket (socket-domain AF_INET))")
		builder.addRule("(allow system-socket (socket-domain AF_INET6))")
		builder.addRule("(allow network-outbound (remote ip \"*:*\"))")
		builder.addRule(fmt.Sprintf("(allow network-outbound (remote unix-socket (path-literal %s)))", quoteSBPL(resolverSocket)))
		builder.addRule("(allow network-bind (local ip \"localhost:*\"))")
		builder.addRule("(allow network-inbound (local ip \"localhost:*\"))")

	}

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
		func() error {
			return builder.addDenyRules(
				policy.Deny,
				"file-read* file-write* file-read-metadata",
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
