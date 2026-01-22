package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// GenerateSlug génère un slug URL-friendly à partir d'un nom
func GenerateSlug(name string) string {
	// Convertir en minuscules
	slug := strings.ToLower(name)

	// Remplacer les caractères accentués
	slug = replaceAccents(slug)

	// Remplacer les espaces et caractères spéciaux par des tirets
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Supprimer les tirets en début et fin
	slug = strings.Trim(slug, "-")

	// Limiter la longueur (optionnel, max 100 caractères)
	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.Trim(slug, "-")
	}

	return slug
}

// replaceAccents remplace les caractères accentués par leurs équivalents non accentués
func replaceAccents(s string) string {
	replacements := map[rune]string{
		'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ä': "a", 'å': "a",
		'è': "e", 'é': "e", 'ê': "e", 'ë': "e",
		'ì': "i", 'í': "i", 'î': "i", 'ï': "i",
		'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ö': "o",
		'ù': "u", 'ú': "u", 'û': "u", 'ü': "u",
		'ý': "y", 'ÿ': "y",
		'ç': "c",
		'ñ': "n",
	}

	var result strings.Builder
	for _, r := range s {
		if replacement, ok := replacements[r]; ok {
			result.WriteString(replacement)
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}
