package engine

import "exocomp/schemas"
import "exocomp/types"
import "encoding/json"
import "fmt"
import "io"
import net_http "net/http"

func logResponseError(console *types.Console, response *net_http.Response) {

	var api_error schemas.Error

	response_payload, _ := io.ReadAll(response.Body)
	err1 := json.Unmarshal(response_payload, &api_error)

	if err1 == nil {
		console.Error(fmt.Sprintf("> API Error: %s", api_error.Error()))
	} else {
		console.Error(fmt.Sprintf("> Unknown Error: \n%s\n", string(response_payload)))
	}

}
