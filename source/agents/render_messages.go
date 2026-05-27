package agents

import "exocomp/schemas"

func render_messages(prompt string) []*schemas.Message {

	result := make([]*schemas.Message, 0)

	result = append(result, &schemas.Message{
		Role:    "system",
		Content: prompt,
		Created: schemas.NewDatetime(),
	})

	return result

}
