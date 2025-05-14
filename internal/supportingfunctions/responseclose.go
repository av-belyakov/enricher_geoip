package supportingfunctions

import "net/http"

func ResponseClose(res *http.Response) {
	if res == nil || res.Body == nil {
		return
	}

	res.Body.Close()
}
