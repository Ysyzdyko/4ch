package dto

import "time"

// Origin представляет данные о месте происхождения персонажа
type Origin struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Location представляет данные о последнем известном местоположении персонажа
type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Character представляет основную модель данных для персонажа
type Character struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Status   string    `json:"status"` // 'Alive', 'Dead' или 'unknown'
	Species  string    `json:"species"`
	Type     string    `json:"type"`     // подвид или тип персонажа
	Gender   string    `json:"gender"`   // 'Female', 'Male', 'Genderless' или 'unknown'
	Origin   Origin    `json:"origin"`   // место происхождения персонажа
	Location Location  `json:"location"` // последнее известное местоположение
	Image    string    `json:"image"`    // URL изображения персонажа (300x300px)
	Episode  []string  `json:"episode"`  // список URL эпизодов
	URL      string    `json:"url"`      // URL endpoint персонажа
	Created  time.Time `json:"created"`  // время создания персонажа в базе данных
}

func (c *Character) GetName() string {
	return c.Name
}

func (c *Character) GetImageURL() string {
	return c.Image
}
