package textlib

import (
	"testing"
)

func TestExtractAdvancedEntities(t *testing.T) {
	text := `John Smith earned $50,000 last year. He lives at 123 Main St, New York, NY 10001. 
	Call him at (555) 123-4567 or email john@example.com. The meeting is at 3:30 PM on January 15, 2024.
	The project is 75% complete and covers 2.5 kilometers.`

	entities := ExtractAdvancedEntities(text)

	if len(entities) == 0 {
		t.Fatal("Expected to find advanced entities")
	}

	// Check for different entity types
	entityTypes := make(map[string]int)
	for _, entity := range entities {
		entityTypes[entity.Type]++
	}

	expectedTypes := []string{
		EntityMoney, EntityPhone, EntityEmail, EntityTime, 
		EntityDate, EntityPercent, EntityMetric,
	}

	for _, expectedType := range expectedTypes {
		if count := entityTypes[expectedType]; count == 0 {
			t.Errorf("Expected to find entity type %s", expectedType)
		}
	}

	// Validate specific entities
	for _, entity := range entities {
		// All entities should have valid positions
		if entity.Position.Start < 0 || entity.Position.End <= entity.Position.Start {
			t.Errorf("Invalid position for entity %s: start=%d, end=%d", 
				entity.Text, entity.Position.Start, entity.Position.End)
		}

		// All entities should have confidence scores
		if entity.Confidence < 0 || entity.Confidence > 1 {
			t.Errorf("Invalid confidence score for entity %s: %f", 
				entity.Text, entity.Confidence)
		}

		// Check specific entity validations
		switch entity.Type {
		case EntityMoney:
			if entity.Attributes["amount"] == "" {
				t.Errorf("Money entity should have amount attribute: %s", entity.Text)
			}
		case EntityPercent:
			if entity.Attributes["value"] == "" {
				t.Errorf("Percent entity should have value attribute: %s", entity.Text)
			}
		case EntityPhone:
			if entity.Attributes["normalized"] == "" {
				t.Errorf("Phone entity should have normalized attribute: %s", entity.Text)
			}
		}
	}
}

func TestMoneyEntityExtraction(t *testing.T) {
	tests := []struct {
		text          string
		expectedCount int
	}{
		{"I paid $100 for this.", 1},
		{"The cost was $1,500.00 and €200.", 2},
		{"No money mentioned here.", 0},
		{"$5.99, $10, and $1,000,000", 3},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			entities := extractMoneyEntities(tt.text)
			if len(entities) != tt.expectedCount {
				t.Errorf("Expected %d money entities, got %d", tt.expectedCount, len(entities))
			}

			// Validate money entities
			for _, entity := range entities {
				if entity.Type != EntityMoney {
					t.Errorf("Expected entity type %s, got %s", EntityMoney, entity.Type)
				}
			}
		})
	}
}

func TestPercentEntityExtraction(t *testing.T) {
	tests := []struct {
		text          string
		expectedCount int
	}{
		{"The completion is 75%.", 1},
		{"We achieved 50% and 25% success rates.", 2},
		{"No percentages here.", 0},
		{"100% and 0% coverage", 2},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			entities := extractPercentEntities(tt.text)
			if len(entities) != tt.expectedCount {
				t.Errorf("Expected %d percent entities, got %d", tt.expectedCount, len(entities))
			}

			// Validate percent entities
			for _, entity := range entities {
				if entity.Type != EntityPercent {
					t.Errorf("Expected entity type %s, got %s", EntityPercent, entity.Type)
				}
			}
		})
	}
}

func TestPhoneEntityExtraction(t *testing.T) {
	tests := []struct {
		text          string
		expectedCount int
	}{
		{"Call me at (555) 123-4567.", 1},
		{"Phone: 555-123-4567 or 555.123.4567", 2},
		{"No phone numbers here.", 0},
		{"International: +1-555-123-4567", 1},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			entities := extractPhoneEntities(tt.text)
			if len(entities) != tt.expectedCount {
				t.Errorf("Expected %d phone entities, got %d", tt.expectedCount, len(entities))
			}

			// Validate phone entities
			for _, entity := range entities {
				if entity.Type != EntityPhone {
					t.Errorf("Expected entity type %s, got %s", EntityPhone, entity.Type)
				}
				
				// Should have normalized attribute
				if entity.Attributes["normalized"] == "" {
					t.Errorf("Phone entity should have normalized attribute: %s", entity.Text)
				}
			}
		})
	}
}

func TestTimeEntityExtraction(t *testing.T) {
	tests := []struct {
		text          string
		expectedCount int
	}{
		{"Meeting at 3:30 PM.", 1},
		{"From 9:00 AM to 5:00 PM", 2},
		{"No time mentioned.", 0},
		{"At 12:00 and 15:30", 2},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			entities := extractTimeEntities(tt.text)
			if len(entities) != tt.expectedCount {
				t.Errorf("Expected %d time entities, got %d", tt.expectedCount, len(entities))
			}

			// Validate time entities
			for _, entity := range entities {
				if entity.Type != EntityTime {
					t.Errorf("Expected entity type %s, got %s", EntityTime, entity.Type)
				}
			}
		})
	}
}

func TestCodeEntityExtraction(t *testing.T) {
	text := `Here's some Python code:
	def hello():
	    print("Hello world")
	
	And some JavaScript:
	function test() { return true; }`

	entities := extractCodeEntities(text)

	if len(entities) == 0 {
		t.Error("Expected to find code entities")
	}

	// Should detect programming languages
	foundLanguages := make(map[string]bool)
	for _, entity := range entities {
		if entity.Type == EntityCode {
			if lang, exists := entity.Attributes["language"]; exists {
				foundLanguages[lang] = true
			}
		}
	}

	if len(foundLanguages) == 0 {
		t.Error("Expected to detect programming languages")
	}
}

func TestNumberEntityExtraction(t *testing.T) {
	text := "I have 5 apples, 3.14 is pi, and this is the 1st attempt."

	entities := extractNumberEntities(text)

	if len(entities) == 0 {
		t.Error("Expected to find number entities")
	}

	// Check for different number types
	numberTypes := make(map[string]bool)
	for _, entity := range entities {
		if entity.Type == EntityNumber {
			if subtype, exists := entity.Attributes["numberType"]; exists {
				numberTypes[subtype] = true
			}
		}
	}

	expectedTypes := []string{"cardinal", "decimal", "ordinal"}
	for _, expectedType := range expectedTypes {
		if !numberTypes[expectedType] {
			t.Errorf("Expected to find number type %s", expectedType)
		}
	}
}

func TestMetricEntityExtraction(t *testing.T) {
	text := "The distance is 5.2 kilometers and the weight is 3.5 pounds."

	entities := extractMetricEntities(text)

	if len(entities) == 0 {
		t.Error("Expected to find metric entities")
	}

	// Check for different metric types
	metricTypes := make(map[string]bool)
	for _, entity := range entities {
		if entity.Type == EntityMetric {
			if metricType, exists := entity.Attributes["metricType"]; exists {
				metricTypes[metricType] = true
			}
		}
	}

	expectedTypes := []string{"distance", "weight"}
	for _, expectedType := range expectedTypes {
		if !metricTypes[expectedType] {
			t.Errorf("Expected to find metric type %s", expectedType)
		}
	}
}

func TestEntityDeduplication(t *testing.T) {
	// Create entities with overlapping positions
	entities := []AdvancedEntity{
		{
			Entity: Entity{
				Type: EntityMoney,
				Text: "$100",
				Position: Position{Start: 0, End: 4},
			},
		},
		{
			Entity: Entity{
				Type: EntityNumber,
				Text: "100",
				Position: Position{Start: 1, End: 4}, // Overlaps with money entity
			},
		},
		{
			Entity: Entity{
				Type: EntityPercent,
				Text: "50%",
				Position: Position{Start: 10, End: 13},
			},
		},
	}

	deduplicated := deduplicateAdvancedEntities(entities)

	// Should remove overlapping entities, keeping the more specific one (money over number)
	if len(deduplicated) != 2 {
		t.Errorf("Expected 2 entities after deduplication, got %d", len(deduplicated))
	}

	// Money entity should be kept over number entity
	hasMoneyEntity := false
	hasNumberEntity := false
	for _, entity := range deduplicated {
		if entity.Type == EntityMoney {
			hasMoneyEntity = true
		}
		if entity.Type == EntityNumber && entity.Text == "100" {
			hasNumberEntity = true
		}
	}

	if !hasMoneyEntity {
		t.Error("Money entity should be kept during deduplication")
	}
	if hasNumberEntity {
		t.Error("Overlapping number entity should be removed during deduplication")
	}
}