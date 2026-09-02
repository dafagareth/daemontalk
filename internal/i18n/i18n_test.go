package i18n

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		lang     string
		expected string
	}{
		{"en", "en"},
		{"id", "id"},
		{"es", "en"}, // Fallback to en
		{"", "en"},   // Fallback to en
		{"de", "en"}, // Fallback to en
	}

	for _, tc := range tests {
		t.Run("lang_"+tc.lang, func(t *testing.T) {
			ui := Get(tc.lang)
			if tc.expected == "id" {
				if ui.Nav_Home != "Beranda" {
					t.Errorf("expected Nav_Home 'Beranda' for id, got %q", ui.Nav_Home)
				}
			} else {
				if ui.Nav_Home != "Home" {
					t.Errorf("expected Nav_Home 'Home' for en/fallback, got %q", ui.Nav_Home)
				}
			}
		})
	}
}

func TestLocalesCompleteness(t *testing.T) {
	locales := []string{"en", "id"}

	for _, lang := range locales {
		t.Run("completeness_"+lang, func(t *testing.T) {
			ui := Get(lang)
			val := reflect.ValueOf(ui)
			typ := val.Type()

			emptyFields := 0
			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)
				fieldName := typ.Field(i).Name
				if field.Kind() == reflect.String && field.String() == "" {
					t.Errorf("locale %s has empty translation for field: %s", lang, fieldName)
					emptyFields++
				}
			}
			if emptyFields > 0 {
				t.Errorf("locale %s has %d empty fields", lang, emptyFields)
			}
		})
	}
}
