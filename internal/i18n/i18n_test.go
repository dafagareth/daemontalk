package i18n

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []string{"en", "id", "es", "", "de"}

	for _, lang := range tests {
		t.Run("lang_"+lang, func(t *testing.T) {
			ui := Get(lang)
			if ui.Nav_Home != "Home" {
				t.Errorf("expected Nav_Home 'Home' for %q, got %q", lang, ui.Nav_Home)
			}
		})
	}
}

func TestLocalesCompleteness(t *testing.T) {
	locales := []string{"en"}

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
