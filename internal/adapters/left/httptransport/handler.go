package http

import (
	"1337b04rd/internal/application/domain"
	ports "1337b04rd/internal/ports/left"
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"1337b04rd/pkg/utils"
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	svc       ports.APIPort
	templates *template.Template
}

func NewHandler(svc ports.APIPort) (*Handler, error) {
	// Вариант 1: Использовать относительный путь от рабочей директории
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		// Вариант 2: Попробовать абсолютный путь в Docker-контейнере
		tmpl, err = template.ParseGlob("/app/templates/*.html")
		if err != nil {
			return nil, fmt.Errorf("failed to parse templates: %w", err)
		}
	}

	return &Handler{
		svc:       svc,
		templates: tmpl,
	}, nil
}

func (h *Handler) HandleCatalog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.svc.GetCatalog(ctx)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		slog.Error("GetCatalog error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "catalog.html", data); err != nil {
		slog.Error("Failed to render template", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
	}
}

func (h *Handler) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.svc.GetPostByID(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Используем буфер для безопасного рендеринга шаблона
	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "post.html", data); err != nil {
		slog.Error("Failed to render template", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
		return
	}

	// Пишем рендеренный контент в ответ
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(buf.Bytes())
	if err != nil {
		slog.Error("Failed to send response", "error", err)
	}
}

func (h *Handler) HandleCreatePostForm(w http.ResponseWriter, r *http.Request) {
	// Проверяем существование шаблона
	if h.templates.Lookup("create-post.html") == nil {
		slog.Error("Template not found", "template", "create-post.html")
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// Добавляем данные, если нужно
	data := struct {
		Title string
	}{
		Title: "Create New Post",
	}

	// Рендерим шаблон
	err := h.templates.ExecuteTemplate(w, "create-post.html", data)
	if err != nil {
		slog.Error("Template execution error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) HandleSubmitPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Проверяем тип файла
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		http.Error(w, "Only images are allowed", http.StatusBadRequest)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read file", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ext := filepath.Ext(header.Filename)
	contentType := header.Header.Get("Content-Type")

	session, ok := r.Context().Value(SessionKey).(*domain.Session)
	if !ok || session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	uuid, err := utils.GenerateUUID()
	if err != nil {
		slog.Error(err.Error())
	}

	post := &domain.Post{
		ID:        uuid,
		Title:     title,
		Content:   content,
		Author:    session.UserID,
		CreatedAt: time.Now(),
	}
	objectName := post.ID + ext

	if err := h.svc.UploadAndCreatePost(ctx, post, objectName, fileBytes, contentType); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}

func (h *Handler) HandleAddComment(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	session, ok := r.Context().Value(SessionKey).(*domain.Session)
	if !ok || session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID := r.URL.Query().Get("id")
	parentID := r.FormValue("parent_comment_id")
	content := r.FormValue("content")

	if content == "" {
		http.Error(w, "Content are required", http.StatusBadRequest)
		return
	}

	uuid, err := utils.GenerateUUID()
	if err != nil {
		slog.Error(err.Error())
	}

	comment := &domain.Comment{
		ID:        uuid,
		Author:    session.UserID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if parentID != "" {
		if err := h.svc.ReplyToComment(ctx, parentID, comment); err != nil {
			http.Error(w, "Failed to add reply: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.svc.AddComment(ctx, postID, comment); err != nil {
			http.Error(w, "Failed to add comment: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

func (h *Handler) HandleArchiveList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.svc.GetArchiveList(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "archive.html", data); err != nil {
		slog.Error("Failed to render template", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
	}

}

func (h *Handler) HandleGetArchivedPost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.svc.GetArchivedPostByID(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := h.templates.ExecuteTemplate(w, "archive-post.html", data); err != nil {
		slog.Error("Failed to render template", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
	}

}
