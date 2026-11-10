module exocomp-app

go 1.27

require (
	exocomp v0.0.0
	github.com/cookiengineer/gooey v0.0.0
)

replace exocomp => ../source

// NOTE: Local sibling checkout at ../gooey relative to the repository root.
replace github.com/cookiengineer/gooey => ../../gooey
