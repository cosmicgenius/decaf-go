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
			Variable[K]{Value: g.start},
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
	aliasProductionRules := make([]ProductionRule[K], 0, len(aliasMap))
	for terminal, variable := range aliasMap {
		aliasProductionRules = append(aliasProductionRules, ProductionRule[K]{
			Input:  variable,
			Output: []Symbol[K]{terminal},
		})
	}
	newG.productionRules = append(newG.productionRules, aliasProductionRules...)

	return newG
}

// Step 3: BIN. Split up any production rules that produce more than 2 symbols on the RHS.
func (g *ContextFreeGrammar[K]) applyCNFBinStep(
	produceNewRandomVariable func() K,
) *ContextFreeGrammar[K] {
	newG := g.Clone()
	rejectionSampleNewVariable := g.rejectionSampleNewVariableFactory(produceNewRandomVariable)

	getSplitProductionRules := func(pr ProductionRule[K]) []ProductionRule[K] {
		// Given a production rule of the form A -> X0 X1 ... X(n-1),
		// split it into
		//     A      -> X0 A0
		//     A0     -> X1 A1
        //     ...
        //     A(n-4) -> X(n-3) A(n-3)
        //     A(n-3) -> X(n-2) X(n-1)
		n := len(pr.Output)

		if n <= 2 {
            return []ProductionRule[K]{pr}
        }

		splitProductionRules := make([]ProductionRule[K], 0, n-1)
		newVariables := make([]K, n-2)
        for i := 0; i < n-2; i++ {
            newVariables[i] = rejectionSampleNewVariable()
        }

		// Add the A -> X0 A0 rule
		splitProductionRules = append(splitProductionRules, ProductionRule[K]{
            Input:  pr.Input,
            Output: []Symbol[K]{
				pr.Output[0],
				Variable[K]{Value: newVariables[0]},
			},
        })
		// Add the rules of the form A(i) -> X(i+1) A(i+1)
		for i := 0; i < n-3; i++ {
			splitProductionRules = append(splitProductionRules, ProductionRule[K]{
				Input:  newVariables[i],
				Output: []Symbol[K]{
					pr.Output[i+1],
					Variable[K]{Value: newVariables[i+1]},
				},
			})
		}
		// Add the A(n-3) -> X(n-2) X(n-1) rule
		splitProductionRules = append(splitProductionRules, ProductionRule[K]{
			Input:  newVariables[n-3],
			Output: []Symbol[K]{pr.Output[n-2], pr.Output[n-1]},
		})

		return splitProductionRules
	}

	// Just put some lower bound reservation
	newG.productionRules = make([]ProductionRule[K], 0, len(g.productionRules))

	for _, pr := range g.productionRules {
		newG.productionRules = append(newG.productionRules, getSplitProductionRules(pr)...)
	}

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
