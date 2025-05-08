package dto

type Info struct {
	Count int     `json:"count"` // Общее количество результатов
	Pages int     `json:"pages"` // Общее количество страниц
	Next  *string `json:"next"`  // URL следующей страницы (может быть nil)
	Prev  *string `json:"prev"`  // URL предыдущей страницы (может быть nil)
}

// CharacterResponse представляет ответ API с пагинацией для персонажей
type CharacterResponse struct {
	Info    Info        `json:"info"`
	Results []Character `json:"results"` // 20 персонажей на странице
}

func (c *CharacterResponse) GetCharactersCount() int {
	return c.Info.Count
}
