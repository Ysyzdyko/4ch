package domain

import "time"

type Comment struct {
	ID         string
	Author     string
	Content    string
	AvaterLink string
	ParentID   string    // ID родительского комментария (пустая строка если это корневой комментарий)
	Replies    []Comment // Вложенные ответы
	CreatedAt  time.Time
}
