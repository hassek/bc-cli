package templates

import (
	"strings"
	"testing"
)

func TestStyleFunctions(t *testing.T) {
	tests := []struct {
		name     string
		template string
		wantErr  bool
		contains string
	}{
		{
			name:     "highlight function",
			template: `{{highlight "test text"}}`,
			wantErr:  false,
			contains: "test text",
		},
		{
			name:     "emphasis function",
			template: `{{emphasis "emphasized text"}}`,
			wantErr:  false,
			contains: "emphasized text",
		},
		{
			name:     "section function",
			template: `{{section "Section Title"}}`,
			wantErr:  false,
			contains: "Section Title",
		},
		{
			name:     "bold function",
			template: `{{bold "bold text"}}`,
			wantErr:  false,
			contains: "bold text",
		},
		{
			name:     "faint function",
			template: `{{faint "faint text"}}`,
			wantErr:  false,
			contains: "faint text",
		},
		{
			name:     "cyan function",
			template: `{{cyan "cyan text"}}`,
			wantErr:  false,
			contains: "cyan text",
		},
		{
			name:     "green function",
			template: `{{green "green text"}}`,
			wantErr:  false,
			contains: "green text",
		},
		{
			name:     "yellow function",
			template: `{{yellow "yellow text"}}`,
			wantErr:  false,
			contains: "yellow text",
		},
		{
			name:     "red function",
			template: `{{red "red text"}}`,
			wantErr:  false,
			contains: "red text",
		},
		{
			name:     "paragraph function",
			template: `{{paragraph "test paragraph" 50}}`,
			wantErr:  false,
			contains: "test paragraph",
		},
		{
			name:     "paragraphAuto function",
			template: `{{paragraphAuto "auto paragraph"}}`,
			wantErr:  false,
			contains: "auto paragraph",
		},
		{
			name:     "combined styles with printf",
			template: `{{printf "This is %s text" (highlight "highlighted")}}`,
			wantErr:  false,
			contains: "This is",
		},
		{
			name:     "nested template functions",
			template: `{{paragraphAuto (printf "We only offer %s — and we mean it!" (highlight "what we like"))}}`,
			wantErr:  false,
			contains: "We only offer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderToString(tt.template, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.Contains(result, tt.contains) {
				t.Errorf("RenderToString() result does not contain %q, got: %q", tt.contains, result)
			}
		})
	}
}

func TestAboutTemplateRenders(t *testing.T) {
	// This simulates what the backend CMS would send
	tmpl := `
{{paragraphAuto "At Butler Coffee we keep things simple."}}

{{section "Our Principle"}}

{{paragraphAuto (printf "%s — nothing goes on our stock unless we'd happily drink it ourselves." (highlight "We only offer what we like"))}}

{{faint "Made with ☕ and love"}}
`

	result, err := RenderToString(tmpl, nil)
	if err != nil {
		t.Fatalf("Failed to render about template: %v", err)
	}

	expectedTexts := []string{
		"At Butler Coffee we keep things simple",
		"Our Principle",
		"We only offer what we like",
		"nothing goes on our stock",
		"Made with ☕ and love",
	}

	for _, expected := range expectedTexts {
		if !strings.Contains(result, expected) {
			t.Errorf("Template output missing expected text: %q", expected)
		}
	}
}

func TestRenderDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		contains    []string
	}{
		{
			name:        "plain text without template syntax",
			description: "This is a simple description",
			contains:    []string{"This is a simple description"},
		},
		{
			name:        "description with highlight",
			description: "We offer {{highlight \"premium coffee\"}} from around the world",
			contains:    []string{"We offer", "premium coffee", "from around the world"},
		},
		{
			name:        "description with multiple styles",
			description: "{{bold \"Premium Coffee\"}} - {{emphasis \"sourced with care\"}}",
			contains:    []string{"Premium Coffee", "sourced with care"},
		},
		{
			name:        "description with section and highlight",
			description: "{{section \"Quality\"}}\\n\\nWe only offer {{highlight \"the best\"}}",
			contains:    []string{"Quality", "We only offer", "the best"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderDescription(tt.description)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("RenderDescription() result does not contain %q, got: %q", expected, result)
				}
			}
		})
	}
}

func TestProductDescriptionWithTemplates(t *testing.T) {
	// Simulate a product description from the backend CMS
	description := "Our {{highlight \"signature blend\"}} combines beans from Ethiopia and Colombia. " +
		"Perfect for {{emphasis \"espresso\"}} or {{emphasis \"pour over\"}}."

	result := RenderDescription(description)

	expectedTexts := []string{
		"Our",
		"signature blend",
		"combines beans",
		"Perfect for",
		"espresso",
		"pour over",
	}

	for _, expected := range expectedTexts {
		if !strings.Contains(result, expected) {
			t.Errorf("Product description missing expected text: %q", expected)
		}
	}
}

func TestSubscriptionDescriptionWithTemplates(t *testing.T) {
	// Simulate a subscription description from the backend CMS
	description := "Get {{highlight \"fresh coffee\"}} delivered monthly. " +
		"{{bold \"Free shipping\"}} on all orders. " +
		"{{faint \"Cancel anytime.\"}}"

	result := RenderDescription(description)

	expectedTexts := []string{
		"Get",
		"fresh coffee",
		"delivered monthly",
		"Free shipping",
		"on all orders",
		"Cancel anytime",
	}

	for _, expected := range expectedTexts {
		if !strings.Contains(result, expected) {
			t.Errorf("Subscription description missing expected text: %q", expected)
		}
	}
}
