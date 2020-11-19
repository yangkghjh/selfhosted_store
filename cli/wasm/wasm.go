package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/yankghjh/selfhosted_store/cli/converter"
)

type convertResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Result  string `json:"result"`
}

func convertFunc(this js.Value, args []js.Value) interface{} {
	output, err := converter.Convert(
		args[0].String(), args[1].String(), []byte(args[2].String()),
	)
	var result convertResult
	if err != nil {
		result = convertResult{
			Success: false,
			Error:   err.Error(),
		}
	} else {
		result = convertResult{
			Success: true,
			Result:  string(output),
		}
	}

	r, _ := json.Marshal(&result)

	return js.ValueOf(string(r))
}

func main() {
	done := make(chan int, 0)
	js.Global().Set("ShsConvert", js.FuncOf(convertFunc))
	<-done
}
