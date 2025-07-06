package parser

type Symbol[K comparable] interface {
	isSymbol()
}

type Variable[K comparable] struct {
	Value K
}

func (v Variable[K]) isSymbol() {}

type Terminal struct {
	Value rune
}

func (t Terminal) isSymbol() {}

func StringToTerminals[K comparable](s string) []Symbol[K] {
	terminals := make([]Symbol[K], 0)
	for _, r := range s {
		terminals = append(terminals, Terminal{Value: r})
	}
	return terminals
}

type ProductionRule[K comparable] struct {
	Input  K
	Output []Symbol[K]
}

type ContextFreeGrammar[K comparable] struct {
	start           K
	productionRules []ProductionRule[K]
}

func NewContextFreeGrammar[K comparable](
	start K,
	productionRules []ProductionRule[K],
) *ContextFreeGrammar[K] {
	return &ContextFreeGrammar[K]{
		start:           start,
		productionRules: productionRules,
	}
}

func (g *ContextFreeGrammar[K]) Clone() *ContextFreeGrammar[K] {
	newProductionRules := make([]ProductionRule[K], len(g.productionRules))
	for i, pr := range g.productionRules {
		newProductionRules[i] = ProductionRule[K]{
			Input:  pr.Input,
			Output: make([]Symbol[K], len(pr.Output)),
		}
		copy(newProductionRules[i].Output, pr.Output)
	}

	return &ContextFreeGrammar[K]{
		start:           g.start,
		productionRules: newProductionRules,
	}
}
