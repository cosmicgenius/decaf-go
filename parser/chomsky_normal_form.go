package parser

import "errors"

var errBINNotPerformed = errors.New("production rule with output length > 2 detected. Did you perform BIN?")

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
	// Just clone here even though it is technically a waste of time.
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
			Input: pr.Input,
			Output: []Symbol[K]{
				pr.Output[0],
				Variable[K]{Value: newVariables[0]},
			},
		})
		// Add the rules of the form A(i) -> X(i+1) A(i+1)
		for i := 0; i < n-3; i++ {
			splitProductionRules = append(splitProductionRules, ProductionRule[K]{
				Input: newVariables[i],
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

// Step 4: DEL. Find and inline (pushdown) all \varepsilon rules.
// This assumes that BIN has already been performed. Otherwise, this pushdown would cause an.
// exponential increase in the number of production rules. We will error in that case.
func (g *ContextFreeGrammar[K]) applyCNFDelStep() (*ContextFreeGrammar[K], error) {
	// Just clone here even though it is technically a waste of time.
	newG := g.Clone()

	// First, find the set of variables that can derive \varepsilon.
	// This is possible iff there exists some production rule A -> X0 X1 ... X(n-1)
	// where all X(i) are nullable.
	//
	// We'll just do some O(#variables * #production rules) loop because efficiency doesn't really matter here,
	// and also this is actually quite tricky to get right (e.g. we need to avoid dying due to cycles).
	nullable := make(map[K]bool)

	// Iterate until nothing changes. Note that this iterates at most #variables times
	// since #nullable is otherwise strictly increasing.
	changed := true
	for changed {
		changed = false

		for _, pr := range g.productionRules {
			if nullable[pr.Input] {
				continue
			}

			// If there is a production rule A -> \varepsilon, clearly A is nullable.
			if len(pr.Output) == 0 {
				nullable[pr.Input] = true
				changed = true
				continue
			}

			allNullable := true
			// If pr is a production rule A -> X0 X1 ... X(n-1) such that
			// all X(i) are nullable, then A is nullable.
			for _, symbol := range pr.Output {
				switch s := symbol.(type) {
				case Terminal:
					allNullable = false
					break
				case Variable[K]:
					if !nullable[s.Value] {
						allNullable = false
						break
					}
				default:
					panic("unreachable")
				}
			}
			if allNullable {
				nullable[pr.Input] = true
				changed = true
			}
		}
	}

	getPushedDownProductionRules := func(pr ProductionRule[K]) ([]ProductionRule[K], error) {
		// Delete any production rule that produces \varepsilon if it is not from the starting variable.
		if len(pr.Output) == 0 {
			if pr.Input == g.start {
				return []ProductionRule[K]{pr}, nil
			}
			return []ProductionRule[K]{}, nil
		}

		// Passthrough any production rule with output length 1
		if len(pr.Output) == 1 {
			return []ProductionRule[K]{pr}, nil
		}

		if len(pr.Output) > 2 {
			return nil, errBINNotPerformed
		}

		// Otherwise, pr is of the form A -> B C
		newProductionRules := []ProductionRule[K]{pr}

		// If B is a nullable variable, add the production rule A -> C
		if variable, ok := pr.Output[0].(Variable[K]); ok && nullable[variable.Value] {
			newProductionRules = append(newProductionRules, ProductionRule[K]{
				Input:  pr.Input,
				Output: []Symbol[K]{pr.Output[1]},
			})
		}

		// Similarly, if C is a nullable variable, add the production rule A -> B
		if variable, ok := pr.Output[1].(Variable[K]); ok && nullable[variable.Value] {
			newProductionRules = append(newProductionRules, ProductionRule[K]{
				Input:  pr.Input,
				Output: []Symbol[K]{pr.Output[0]},
			})
		}

		return newProductionRules, nil
	}

	// Just put some lower bound reservation
	newG.productionRules = make([]ProductionRule[K], 0, len(g.productionRules))

	for _, pr := range g.productionRules {
		pushedDownProductionRules, err := getPushedDownProductionRules(pr)
		if err != nil {
			return nil, err
		}
		newG.productionRules = append(newG.productionRules, pushedDownProductionRules...)
	}

	return newG, nil
}

// Step 5: UNIT. Pushdown all rules of the form A -> B.
func (g *ContextFreeGrammar[K]) applyCNFUnitStep() *ContextFreeGrammar[K] {
	// Clone even though it's a bit useless.
	newG := g.Clone()

	// Construct a graph of all unital production rules, i.e. ones of the form A -> B.
	unitalProductionRuleGraph := make(map[K]map[K]bool)
	variables := g.getVariableNameSet()

	// Initialize adjacency matrix with just the identity
	for v := range variables {
		unitalProductionRuleGraph[v] = make(map[K]bool)
		for u := range variables {
			unitalProductionRuleGraph[v][u] = false
		}
		unitalProductionRuleGraph[v][v] = true
	}

	// Save a map of all non unital production rules.
	nonUnitalProductionRules := make(map[K][]ProductionRule[K])

	for _, pr := range g.productionRules {
		// Unital
		if len(pr.Output) == 1 {
			if output, ok := pr.Output[0].(Variable[K]); ok {
				unitalProductionRuleGraph[pr.Input][output.Value] = true
				continue
			}
		}

		nonUnitalProductionRules[pr.Input] = append(nonUnitalProductionRules[pr.Input], pr)
	}

	// Compute the connectivity graph using Floyd-Warshall. This is fast enough for us.
	for k := range unitalProductionRuleGraph {
		for i := range unitalProductionRuleGraph {
			for j := range unitalProductionRuleGraph {
				if unitalProductionRuleGraph[i][j] {
					continue
				}

				unitalProductionRuleGraph[i][j] = unitalProductionRuleGraph[i][k] && unitalProductionRuleGraph[k][j]
			}
		}
	}

	// For each v, compute all the production rules that can be obtained
	// by repeatedly applying unital rules, and then applying 1 non-unital rule.
	allPulledBackProductionRules := make(map[K][]ProductionRule[K])
	for v := range unitalProductionRuleGraph {
		for u := range unitalProductionRuleGraph[v] {
			if unitalProductionRuleGraph[v][u] {
				// Change the input from u to v for each non-unital production rule with input u
				pulledBackProductionRules := make([]ProductionRule[K], len(nonUnitalProductionRules[u]))
				for i, pr := range nonUnitalProductionRules[u] {
					pulledBackProductionRules[i] = ProductionRule[K]{
						Input:  v,
						Output: pr.Output,
					}
				}
				allPulledBackProductionRules[v] = append(allPulledBackProductionRules[v], pulledBackProductionRules...)
			}
		}
	}

	// Prune by keeping only the pulled back production rules
	// that are reachable from the starting variable.
	reachable := make(map[K]struct{})
	var traverse func(v K)
	traverse = func(v K) {
		if _, ok := reachable[v]; ok {
			return
		}
		reachable[v] = struct{}{}
		for _, pr := range allPulledBackProductionRules[v] {
			for _, output := range pr.Output {
				if outputVariable, ok := output.(Variable[K]); ok {
					traverse(outputVariable.Value)
				}
			}
		}
	}
	traverse(g.start)

	newG.productionRules = make([]ProductionRule[K], 0)
	for v := range reachable {
		newG.productionRules = append(newG.productionRules, allPulledBackProductionRules[v]...)
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

	// Apply each of the steps in succession
	var err error
	newG = newG.applyCNFStartStep(produceNewRandomVariable)
	newG = newG.applyCNFTermStep(produceNewRandomVariable)
	newG = newG.applyCNFBinStep(produceNewRandomVariable)
	newG, err = newG.applyCNFDelStep()
	// Should only happen if applyCNFDelStep is called before applyCNFBinStep
	// which is obviously not possible here.
	if err != nil {
		panic(err)
	}
	newG = newG.applyCNFUnitStep()

	return newG
}
