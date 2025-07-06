package parser

func (g *ContextFreeGrammar[K]) IsInChomskyNormalForm() bool {
	for _, pr := range g.productionRules {
		if len(pr.Output) == 0 {
			if pr.Input != g.start {
				// Only allowed production rule of the form A -> \varepsilon
				// is from the start symbol.
				return false
			}
		} else if len(pr.Output) == 1 {
			// Only allowed production rules with output of length 1
			// are of the form A -> a where a is a terminal.
			if _, ok := pr.Output[0].(Terminal); !ok {
				return false
			}
		} else if len(pr.Output) == 2 {
			// Only allowed production rules with output of length 2
			// are of the form A -> B C where B and C are variables that are not the start variable.
			for _, symbol := range pr.Output {
				if variable, ok := symbol.(Variable[K]); !ok || variable.Value == g.start {
					return false
				}
			}
		} else {
			// No production rules with output length > 2 are allowed.
			return false
		}
	}
	return true
}

// Step 1: START. If the start variable appears in a production rule on the RHS,
// produce an additional start variable.
func (g *ContextFreeGrammar[K]) applyCNFStartStep(
	produceNewRandomVariable func() K,
) *ContextFreeGrammar[K] {
	newG := g.Clone()
	rejectionSampleNewVariable := g.rejectionSampleNewVariableFactory(produceNewRandomVariable)

	shouldAppendNewStart := false
	for _, pr := range g.productionRules {
		for _, symbol := range pr.Output {
			if variable, ok := symbol.(Variable[K]); ok && variable.Value == g.start {
				shouldAppendNewStart = true
				break
			}
		}
		if shouldAppendNewStart {
			break
		}
	}

	if !shouldAppendNewStart {
		return newG
	}

	newG.start = rejectionSampleNewVariable()
	newG.productionRules = append(newG.productionRules, ProductionRule[K]{
		Input: newG.start,
		Output: []Symbol[K]{
			Variable[K]{Value: newG.start},
		},
	})
	return newG
}

// Step 2: TERM. For each production rule that with output length > 1,
// alias the terminals with a variable so all such rules produce only variables.
func (g *ContextFreeGrammar[K]) applyCNFTermStep(
	produceNewRandomVariable func() K,
) *ContextFreeGrammar[K] {
	newG := g.Clone()
	rejectionSampleNewVariable := g.rejectionSampleNewVariableFactory(produceNewRandomVariable)

	aliasMap := make(map[Terminal]K)

	for i, pr := range g.productionRules {
		if len(pr.Output) <= 1 {
			continue
		}

		newPr := ProductionRule[K]{Input: pr.Input, Output: make([]Symbol[K], len(pr.Output))}
		for j, symbol := range pr.Output {
			newPr.Output[j] = symbol
			if terminal, ok := symbol.(Terminal); ok {
				// Define new alias if it doesn't already exist.
				if _, exists := aliasMap[terminal]; !exists {
					aliasMap[terminal] = rejectionSampleNewVariable()
				}

				// Replace with alias.
				newPr.Output[j] = Variable[K]{Value: aliasMap[terminal]}
			}
		}
		newG.productionRules[i] = newPr
	}

	// Append new rules that go from the aliased variables to the original terminals.
	aliasProductionRules := make([]ProductionRule[K], len(aliasMap))
	for terminal, variable := range aliasMap {
		aliasProductionRules = append(aliasProductionRules, ProductionRule[K]{
			Input:  variable,
			Output: []Symbol[K]{terminal},
		})
	}
	newG.productionRules = append(newG.productionRules, aliasProductionRules...)

	return newG
}

// Get the set of all variables in the grammar.
func (g *ContextFreeGrammar[K]) getVariableNameSet() map[K]struct{} {
	variableNameSet := make(map[K]struct{})

	for _, pr := range g.productionRules {
		variableNameSet[pr.Input] = struct{}{}
	}

	return variableNameSet
}

// Assume that produceNewRandomVariable is actually random, and rejection sample
// for a unique new variable name.
func (g *ContextFreeGrammar[K]) rejectionSampleNewVariableFactory(
	produceNewRandomVariable func() K,
) func() K {
	variableNameSet := g.getVariableNameSet()
	return func() K {
		for {
			ret := produceNewRandomVariable()
			if _, ok := variableNameSet[ret]; !ok {
				return ret
			}
		}
	}
}

// Convert a context-free grammar to Chomsky normal form.
// Unfortunately, since K is generic, we need a source of new variable names.
// Get this through produceNewRandomVariable().
func (g *ContextFreeGrammar[K]) ToChomskyNormalForm(
	produceNewRandomVariable func() K,
) *ContextFreeGrammar[K] {
	newG := g

	// Step 1: START. If the start variable appears in a production rule on the RHS,
	// produce an additional start variable.
	newG = newG.applyCNFStartStep(produceNewRandomVariable)

	// Step 2: TERM. For each production rule that with output length > 1,
	// alias the terminals with a variable so all such rules produce only variables.
	newG = newG.applyCNFTermStep(produceNewRandomVariable)

	return newG
}
