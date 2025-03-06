package routes

import (
	"encoding/json"
	"net/http"
)

const InternalServerError = "Internal Server Error"

// jsonResponse 函数用于生成 JSON 响应。
// 它接受一个 http.ResponseWriter 用于写入响应，
// 一个状态码 status 表示 HTTP 响应状态，
// 以及一个 data 接口类型的数据，将被转换为 JSON 格式。
// 函数返回一个错误，如果在处理过程中遇到任何问题。
func jsonResponse(w http.ResponseWriter, status int, data interface{}) error { //nolint:gofmt
	// 设置响应的 Content-Type 为 application/json。
	w.Header().Set("Content-Type", "application/json")

	// 将 data 参数转换为 JSON 格式。
	jsonData, err := json.Marshal(data)
	if err != nil {
		// 如果转换过程中出现错误，返回 500 Internal Server Error。
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return err
	}

	// 发送 HTTP 响应状态码。
	w.WriteHeader(status)

	// 写入转换后的 JSON 数据到响应体中。
	_, err = w.Write(jsonData)
	if err != nil {
		// 如果写入过程中出现错误，返回 500 Internal Server Error。
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return err
	}

	// 如果一切正常，返回 nil 表示没有发生错误。
	return nil
}
