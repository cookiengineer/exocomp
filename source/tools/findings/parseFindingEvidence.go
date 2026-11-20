package findings

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func parseFindingEvidence(arguments utils_chat.Arguments) (types.FindingEvidence, bool, error) {

	if len(arguments) == 0 {
		return types.FindingEvidence{}, false, nil
	}

	kind_value, ok0 := arguments.String("kind")

	if ok0 == false {
		return types.FindingEvidence{}, false, fmt.Errorf("findings.BindEvidence: Invalid parameter \"kind\" is not a string.")
	}

	kind := types.FindingEvidenceKind(strings.ToLower(strings.TrimSpace(kind_value)))
	content, has_content := arguments.String("content")
	artifact, has_artifact := arguments.String("artifact")
	title, has_title := arguments.String("title")
	source, has_source := arguments.String("source")
	command, has_command := arguments.String("command")
	content_type, has_content_type := arguments.String("content_type")
	environment, has_environment := arguments.Map("environment")
	exit_code, has_exit_code := arguments.Int("exit_code")

	if kind.IsValid() == false {
		return types.FindingEvidence{}, false, fmt.Errorf("findings.BindEvidence: Invalid evidence kind \"%s\".", kind_value)
	} else if has_content == has_artifact {
		return types.FindingEvidence{}, false, fmt.Errorf("findings.BindEvidence: Exactly one of \"content\" or \"artifact\" is required.")
	}

	evidence := types.FindingEvidence{
		Kind:    kind,
		Title:   string(kind),
		Content: content,
	}

	if has_title == true {
		evidence.Title = utils_findings.LimitText(title, findingsMaxTitle)
	}

	if has_source == true {
		evidence.Source = utils_findings.LimitText(source, findingsMaxTitle)
	}

	if has_command == true {
		evidence.Command = utils_findings.LimitText(command, findingsMaxText)
	}

	if has_content_type == true {
		evidence.ContentType = utils_findings.LimitText(content_type, findingsMaxTitle)
	}

	if has_environment == true {
		evidence.Environment = environment
	}

	if has_exit_code == true {
		evidence.ExitCode = &exit_code
	}

	if has_artifact == true {
		evidence.Artifact = artifact
	}

	return evidence, true, nil

}
