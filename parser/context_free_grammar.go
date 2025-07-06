package parser

import (
	"fmt"
	"io"
	"os"
)

type Symbol[K comparable] interface {
    IsSymbol()
}

type Variable[K comparable] struct {
    Value K
}

func (v Variable[K]) isSymbol() {}

type Terminal struct {
	Value rune
}

func (t Terminal) isSymbol() {}

type ProductionRule[K comparable] struct {
	Input K
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
		newPr[i] = ProductionRule[K]{
			Input:  pr.Input,
			Output: make([]Symbol[K], len(pr.Output)),
		}
		copy(newPr[i].Output, pr.Output)
	}

    return &ContextFreeGrammar[K]{
        start:           g.start,
        productionRules: newProductionRules,
    }
}
