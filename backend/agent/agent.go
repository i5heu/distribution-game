package agent

import (
	"encoding/json"
	"time"

	v8 "rogchap.com/v8go"
)

func New(agentID string, jsCode string) {
	iso := v8.NewIsolate()

	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		data, err := interfaceToString(map[string]interface{}{
			"msg":  info.Args()[0].String(),
			"time": time.Now().Format(time.RFC3339),
		})
		if err != nil {
			// helper.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
			return nil
		}

		val, err := v8.NewValue(iso, data)
		if err != nil {
			// helper.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
			return nil
		}
		// helper.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		// 	"msg":  info.Args()[0].String(),
		// 	"time": time.Now().Format(time.RFC3339),
		// })

		return val
	})
	global := v8.NewObjectTemplate(iso)
	global.Set("test", printfn)
	ctx := v8.NewContext(iso, global)

	_, err := ctx.RunScript(jsCode, "main.js") // execute some JS code
	if err != nil {
		// helper.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
}

func interfaceToString(data interface{}) (string, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(val), nil
}
