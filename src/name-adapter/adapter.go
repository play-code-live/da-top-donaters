package name_adapter

import (
	"strings"
)

type NameAdapter struct {
	names map[string]string
	skip  []string
}

func NewNameAdapter() *NameAdapter {
	return &NameAdapter{names: map[string]string{
		"Обэма": "Молодой Христос",
		"Аноним (точно не Гера)":              "Твоё позитивное отражение в зеркале",
		"Анонимное отражение в зеркале":       "Твоё позитивное отражение в зеркале",
		"Твое позитивное отражение в зеркале": "Твоё позитивное отражение в зеркале",
		"ДжоршШ":         "Твоё позитивное отражение в зеркале",
		"Херсон":         "Fatcock",
		"Gay":            "Fatcock",
		"":               "Аноним",
		"Sergey_tsurkan": "Fatcock",
		"carrie":         "Carrie",
	}, skip: []string{
		"asdsadsa",
		"play_code",
		"don",
		"he110_todd",
		"Отладка",
		"фывавы а",
		"Отладчик",
		"Аноним",
	}}
}

func (a *NameAdapter) Perform(name string) string {
	name = strings.TrimSpace(name)
	if changed, found := a.names[name]; found {
		return changed
	}
	return name
}

func (a *NameAdapter) ShouldBeSkipped(name string) bool {
	for _, s := range a.skip {
		if name == s {
			return true
		}
	}

	return false
}
