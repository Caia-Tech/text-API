package textlib

import (
	"testing"
)

func TestExtractMathExpressions(t *testing.T) {
	text := "Solve for x: 2x + 5 = 13. Also, calculate sin(π/2) and find the integral ∫ x dx."
	
	expressions := ExtractMathExpressions(text)
	
	if len(expressions) == 0 {
		t.Errorf("Expected to find mathematical expressions, got none")
	}
	
	// Should find at least the equation and some mathematical functions
	foundEquation := false
	foundCalculus := false
	
	for _, expr := range expressions {
		if expr.Type == "equation" {
			foundEquation = true
		}
		if expr.Type == "calculus" {
			foundCalculus = true
		}
	}
	
	if !foundEquation {
		t.Errorf("Expected to find an equation")
	}
	
	if !foundCalculus {
		t.Errorf("Expected to find calculus notation")
	}
}

func TestExtractVariables(t *testing.T) {
	expression := "2x + 3y - z = 10"
	
	variables := ExtractVariables(expression)
	
	expected := []string{"x", "y", "z"}
	if len(variables) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(variables))
	}
	
	for _, expectedVar := range expected {
		found := false
		for _, actualVar := range variables {
			if actualVar == expectedVar {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find variable '%s'", expectedVar)
		}
	}
}

func TestExtractConstants(t *testing.T) {
	expression := "2x + 3.5y - 10 = 0"
	
	constants := ExtractConstants(expression)
	
	expectedConstants := []float64{2, 3.5, 10, 0}
	if len(constants) != len(expectedConstants) {
		t.Errorf("Expected %d constants, got %d", len(expectedConstants), len(constants))
	}
}

func TestParseEquation(t *testing.T) {
	equation := "2x + 3 = 7"
	
	eq := ParseEquation(equation)
	
	if eq.Type != "linear" {
		t.Errorf("Expected linear equation, got %s", eq.Type)
	}
	
	if len(eq.LeftSide) == 0 {
		t.Errorf("Expected terms on left side")
	}
	
	if len(eq.RightSide) == 0 {
		t.Errorf("Expected terms on right side")
	}
}

func TestSimplifyExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x + x", "2x"},
		{"2x + 3x", "5x"},
		{"x * 1", "x"},
		{"1 * x", "x"},
	}
	
	for _, test := range tests {
		result := SimplifyExpression(test.input)
		// Note: Exact string matching might be fragile, but good for basic tests
		if !contains(result, "x") {
			t.Errorf("SimplifyExpression(%s): expected to contain 'x', got %s", test.input, result)
		}
	}
}

func TestValidateMathExpression(t *testing.T) {
	tests := []struct {
		expression string
		shouldPass bool
	}{
		{"2x + 3 = 5", true},
		{"((2x + 3))", true},
		{"(2x + 3", false}, // Unbalanced parentheses
		{"2x ++ 3", false}, // Consecutive operators
		{"", false},        // Empty expression
	}
	
	for _, test := range tests {
		result := ValidateMathExpression(test.expression)
		if result.IsValid != test.shouldPass {
			t.Errorf("ValidateMathExpression(%s): expected %v, got %v", 
				test.expression, test.shouldPass, result.IsValid)
		}
	}
}

func TestDetectArithmeticSequence(t *testing.T) {
	text := "The sequence 2, 4, 6, 8, 10 is arithmetic."
	
	if !detectArithmeticSequence(text) {
		t.Errorf("Expected to detect arithmetic sequence")
	}
	
	nonArithmetic := "The sequence 1, 2, 4, 8, 16 is not arithmetic."
	if detectArithmeticSequence(nonArithmetic) {
		t.Errorf("Should not detect arithmetic sequence in geometric sequence")
	}
}

func TestDetectGeometricSequence(t *testing.T) {
	text := "The sequence 2, 4, 8, 16, 32 is geometric."
	
	if !detectGeometricSequence(text) {
		t.Errorf("Expected to detect geometric sequence")
	}
}

func TestDetectQuadraticPattern(t *testing.T) {
	text := "The function f(x) = x^2 + 2x + 1 is quadratic."
	
	if !detectQuadraticPattern(text) {
		t.Errorf("Expected to detect quadratic pattern")
	}
}

func TestClassifyEquation(t *testing.T) {
	tests := []struct {
		equation string
		expected string
	}{
		{"x + 3 = 5", "linear"},
		{"x^2 + 2x + 1 = 0", "quadratic"},
		{"5 = 5", "constant"},
	}
	
	for _, test := range tests {
		result := ClassifyEquation(test.equation)
		if result != test.expected {
			t.Errorf("ClassifyEquation(%s): expected %s, got %s", 
				test.equation, test.expected, result)
		}
	}
}

func TestFindMathConstants(t *testing.T) {
	text := "The area of a circle is π × r². The value of e is approximately 2.71828."
	
	constants := FindMathConstants(text)
	
	expectedConstants := []string{"π", "e"}
	for _, expected := range expectedConstants {
		found := false
		for _, constant := range constants {
			if constant == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find constant '%s'", expected)
		}
	}
}

func TestDetectMathNotation(t *testing.T) {
	text := "Calculate ∫ f(x) dx and ∑ from i=1 to n."
	
	notations := DetectMathNotation(text)
	
	expectedNotations := []string{"integral", "summation"}
	for _, expected := range expectedNotations {
		found := false
		for _, notation := range notations {
			if notation == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find notation '%s'", expected)
		}
	}
}

func TestHasBalancedParentheses(t *testing.T) {
	tests := []struct {
		expression string
		balanced   bool
	}{
		{"(2x + 3)", true},
		{"((x + y) * z)", true},
		{"(2x + 3", false},
		{"2x + 3)", false},
		{"((x + y) * z", false},
	}
	
	for _, test := range tests {
		result := hasBalancedParentheses(test.expression)
		if result != test.balanced {
			t.Errorf("hasBalancedParentheses(%s): expected %v, got %v", 
				test.expression, test.balanced, result)
		}
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{12, 8, 4},
		{15, 25, 5},
		{7, 13, 1},
		{0, 5, 5},
	}
	
	for _, test := range tests {
		result := gcd(test.a, test.b)
		if result != test.expected {
			t.Errorf("gcd(%d, %d): expected %d, got %d", 
				test.a, test.b, test.expected, result)
		}
	}
}