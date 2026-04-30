package response

//统一JSON返回格式
import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	resp := Response{
		Code: status,
		Msg:  message,
		Data: data,
	}
	//json.NewEncoder(w) 得到一个Encoder
	//Encoder.Encode(resp) 把resp编码成JSON格式，写入w
	_ = json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, "success", data)
}
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, message, nil)
}
