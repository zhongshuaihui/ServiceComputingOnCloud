package service

import (
	"github.com/unrolled/render"
	"net/http"
)

type Comment struct {
	Message 	string `json:"message"`
}

var comment Comment

func submitHandler(webRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			req.ParseForm()
			comment.Message = req.FormValue("message")
		}
		http.FileServer(http.Dir(webRoot + "/assets/")).ServeHTTP(w, req)
	}
}

func showHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, comment)
	}
}

func NotImplemented(w http.ResponseWriter, r *http.Request) { http.Error(w, "501 api not implemented", http.StatusNotImplemented) }

func NotImplementedHandler() http.Handler { return http.HandlerFunc(NotImplemented) }