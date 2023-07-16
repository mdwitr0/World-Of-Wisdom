package quote

import "math/rand"

type IRepository interface {
	GetRandomQuote() Quote
}

type Repository struct {
	quotes []Quote
}

func NewRepository() IRepository {
	var quotes = []Quote{
		{
			Text:   "I like to listen. I have learned a great deal from listening carefully. Most people never listen.",
			Author: "Ernest Hemingway",
		},
		{
			Text:   "I think, that if the world were a bit more like ComicCon, it would be a better place.",
			Author: "Matt Smith",
		},
		{
			Text:   "Voice is not just the sound that comes from your throat, but the feelings that come from your words.",
			Author: "Jennifer Donnelly",
		},
		{
			Text:   "The worst part of being okay is that okay is far from happy.",
			Author: "Anna Todd",
		},
		{
			Text:   "The truest wisdom is a resolute determination.",
			Author: "Napoleon Bonaparte",
		},
		{
			Text:   "Science is organized knowledge. Wisdom is organized life.",
			Author: "Immanuel Kant",
		},
		{
			Text:   "Wisdom is the power to put our time and our knowledge to the proper use.",
			Author: "Thomas J. Watson",
		},
		{
			Text:   "A wise man never loses anything, if he has himself.",
			Author: "Michel de Montaigne",
		},
	}

	return &Repository{
		quotes: quotes,
	}
}

func (r *Repository) GetRandomQuote() Quote {
	randomIndex := rand.Int() % len(r.quotes)

	return r.quotes[randomIndex]
}
