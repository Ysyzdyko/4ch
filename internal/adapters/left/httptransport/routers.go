package http

import (
	lports "1337b04rd/internal/ports/left"
	"net/http"
)

func RegisterRoutes(svc lports.APIPort, router *http.ServeMux) {
	h, err := NewHandler(svc)
	if err != nil {
		panic("failed to create handler: " + err.Error())
	}

	router.HandleFunc("GET /catalog", h.HandleCatalog)
	router.HandleFunc("GET /post/{id}", h.HandleGetPost)
	router.HandleFunc("GET /archive.html", h.HandleArchiveList)
	router.HandleFunc("GET /archive/post/{id}", h.HandleGetArchivedPost)
	router.HandleFunc("GET /create-post.html", h.HandleCreatePostForm)
	router.HandleFunc("GET /create-post", h.HandleCreatePostForm) // форма создания
	router.HandleFunc("POST /submit-post", h.HandleSubmitPost)    // отправка формы
	router.HandleFunc("POST /post/submit-comment", h.HandleAddComment)
}
