package agents

var Roles map[string]string = map[string]string{

	// Development
	"planner":    "writes with humans and plans projects phases",
	"architect":  "defines software specifications",
	"coder":      "implements specifications into Go code",
	// TODO: "researcher": "reads websites and API documentation",
	"summarizer": "reads long texts and summarizes them",
	"tester":     "implements unit tests, writes bug reports",

	// Pentesting
	"exploiter":    "implements exploits in CGo code",
	// TODO: "reverser":     "translates binaries into Go/C/CGo code",
	// TODO: "threathunter": "discovers weaknesses and vulnerabilities in infrastructure",
	"webscanner":   "discovers weaknesses and vulnerabilities in web applications",

}

