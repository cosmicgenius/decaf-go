package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func productionRuleToString[K comparable](pr ProductionRule[K]) string {
	var outputStringBuilder strings.Builder
	for _, symbol := range pr.Output {
		var outputString string
		switch s := symbol.(type) {
		case Terminal:
			outputString = fmt.Sprintf("Term(%s)", string([]rune{s.Value}))
		case Variable[K]:
			outputString = fmt.Sprintf("Var(%v)", s.Value)
		default:
			panic("Unknown symbol type")
		}
        outputStringBuilder.WriteString(outputString)
    }
    return fmt.Sprintf("%v -> %s", pr.Input, outputStringBuilder.String())
}

func sameProductionRules[K comparable](
	prs1 []ProductionRule[K],
    prs2 []ProductionRule[K],
) bool {
	if len(prs1) != len(prs2) {
        return false
    }

	difference := make(map[string]int)

	for _, pr1 := range prs1 {
		difference[productionRuleToString(pr1)]++
	}

	for _, pr2 := range prs2 {
        difference[productionRuleToString(pr2)]--
    }

	for _, v := range difference {
        if v != 0 {
			return false
		}
	}
	return true
}

var (
	palindromeOrig = NewContextFreeGrammar[string](
		"S", /* start */
        []ProductionRule[string]{
            {"S", append(append(
				StringToTerminals[string]("a"),
				[]Symbol[string]{Variable[string]{Value: "S"}}...),
				StringToTerminals[string]("a")...)},
            {"S", append(append(
				StringToTerminals[string]("b"),
				[]Symbol[string]{Variable[string]{Value: "S"}}...),
				StringToTerminals[string]("b")...)},
            {"S", StringToTerminals[string]("a")},
            {"S", StringToTerminals[string]("b")},
            {"S", []Symbol[string]{}},
        },
	)
)

func TestIsInChomskyNormalForm(t *testing.T) {
	assert.False(t, palindromeOrig.IsInChomskyNormalForm())
}

