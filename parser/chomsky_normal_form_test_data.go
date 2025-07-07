package parser

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

	palindromeAfterStart = NewContextFreeGrammar[string](
		"V0", /* start */
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
			{"V0", []Symbol[string]{Variable[string]{Value: "S"}}},
		},
	)

	palindromeAfterTerm = NewContextFreeGrammar[string](
		"V0", /* start */
		[]ProductionRule[string]{
			{"S", []Symbol[string]{
				Variable[string]{Value: "V1"},
				Variable[string]{Value: "S"},
				Variable[string]{Value: "V1"},
			}},
			{"S", []Symbol[string]{
				Variable[string]{Value: "V2"},
				Variable[string]{Value: "S"},
				Variable[string]{Value: "V2"},
			}},
			{"S", StringToTerminals[string]("a")},
			{"S", StringToTerminals[string]("b")},
			{"S", []Symbol[string]{}},
			{"V0", []Symbol[string]{Variable[string]{Value: "S"}}},
			{"V1", StringToTerminals[string]("a")},
			{"V2", StringToTerminals[string]("b")},
		},
	)

	palindromeAfterBin = NewContextFreeGrammar[string](
		"V0", /* start */
		[]ProductionRule[string]{
			{"S", []Symbol[string]{
				Variable[string]{Value: "V1"},
				Variable[string]{Value: "V3"},
			}},
			{"V3", []Symbol[string]{
				Variable[string]{Value: "S"},
				Variable[string]{Value: "V1"},
			}},
			{"S", []Symbol[string]{
				Variable[string]{Value: "V2"},
				Variable[string]{Value: "V4"},
			}},
			{"V4", []Symbol[string]{
				Variable[string]{Value: "S"},
				Variable[string]{Value: "V2"},
			}},
			{"S", StringToTerminals[string]("a")},
			{"S", StringToTerminals[string]("b")},
			{"S", []Symbol[string]{}},
			{"V0", []Symbol[string]{Variable[string]{Value: "S"}}},
			{"V1", StringToTerminals[string]("a")},
			{"V2", StringToTerminals[string]("b")},
		},
	)

	palindromeAfterDel = NewContextFreeGrammar[string](
		"V0", /* start */
		[]ProductionRule[string]{
			{"S", []Symbol[string]{
				Variable[string]{Value: "V1"},
				Variable[string]{Value: "V3"},
			}},
			{"V3", []Symbol[string]{
				Variable[string]{Value: "S"},
				Variable[string]{Value: "V1"},
			}},
			{"V3", []Symbol[string]{Variable[string]{Value: "V1"}}},
			{"S", []Symbol[string]{
				Variable[string]{Value: "V2"},
				Variable[string]{Value: "V4"},
			}},
			{"V4", []Symbol[string]{
				Variable[string]{Value: "S"},
				Variable[string]{Value: "V2"},
			}},
			{"V4", []Symbol[string]{Variable[string]{Value: "V2"}}},
			{"S", StringToTerminals[string]("a")},
			{"S", StringToTerminals[string]("b")},
			{"V0", []Symbol[string]{Variable[string]{Value: "S"}}},
			{"V1", StringToTerminals[string]("a")},
			{"V2", StringToTerminals[string]("b")},
		},
	)

	palindromeAfterUnit = NewContextFreeGrammar[string](
		"V0", /* start */
		[]ProductionRule[string]{
			{"V0", []Symbol[string]{
				Variable[string]{Value: "V1"},
				Variable[string]{Value: "V3"},
			}},
			{"V0", []Symbol[string]{
				Variable[string]{Value: "V2"},
				Variable[string]{Value: "V4"},
			}},
			{"V0", StringToTerminals[string]("a")},
			{"V0", StringToTerminals[string]("b")},
			{"S", []Symbol[string]{
				Variable[string]{Value: "V1"},
				Variable[string]{Value: "V3"},
			}},
			{"V3", []Symbol[string]{
				Variable[string]{Value: "S"},
				Variable[string]{Value: "V1"},
			}},
			{"V3", StringToTerminals[string]("a")},
			{"S", []Symbol[string]{
				Variable[string]{Value: "V2"},
				Variable[string]{Value: "V4"},
			}},
			{"V4", []Symbol[string]{
				Variable[string]{Value: "S"},
				Variable[string]{Value: "V2"},
			}},
			{"V4", StringToTerminals[string]("b")},
			{"S", StringToTerminals[string]("a")},
			{"S", StringToTerminals[string]("b")},
			{"V1", StringToTerminals[string]("a")},
			{"V2", StringToTerminals[string]("b")},
		},
	)
)
