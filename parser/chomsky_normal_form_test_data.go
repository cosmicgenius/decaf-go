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

var (
	arithmeticOrig = NewContextFreeGrammar[string](
		"expr", /* start */
		[]ProductionRule[string]{
			{"expr", []Symbol[string]{Variable[string]{Value: "term"}}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "expr"},
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},
			{"term", []Symbol[string]{Variable[string]{Value: "factor"}}},
			{"term", []Symbol[string]{
				Variable[string]{Value: "term"},
				Variable[string]{Value: "mul_op"},
				Variable[string]{Value: "factor"},
			}},
			{"factor", []Symbol[string]{Variable[string]{Value: "primary"}}},
			{"factor", append(append(
				[]Symbol[string]{Variable[string]{Value: "factor"}},
				StringToTerminals[string]("^")...),
				[]Symbol[string]{Variable[string]{Value: "primary"}}...)},
			{"primary", append(append(
				StringToTerminals[string]("("),
				[]Symbol[string]{Variable[string]{Value: "expr"}}...),
				StringToTerminals[string](")")...)},
			{"primary", StringToTerminals[string]("0")},
			{"primary", StringToTerminals[string]("1")},
			{"primary", StringToTerminals[string]("a")},
			{"primary", StringToTerminals[string]("b")},
			{"add_op", StringToTerminals[string]("+")},
			{"add_op", StringToTerminals[string]("-")},
			{"mul_op", StringToTerminals[string]("*")},
			{"mul_op", StringToTerminals[string]("/")},
		},
	)

	arithmeticAfterStart = NewContextFreeGrammar[string](
		"Vstart", /* start */
		[]ProductionRule[string]{
			{"Vstart", []Symbol[string]{Variable[string]{Value: "expr"}}},
			{"expr", []Symbol[string]{Variable[string]{Value: "term"}}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "expr"},
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},
			{"term", []Symbol[string]{Variable[string]{Value: "factor"}}},
			{"term", []Symbol[string]{
				Variable[string]{Value: "term"},
				Variable[string]{Value: "mul_op"},
				Variable[string]{Value: "factor"},
			}},
			{"factor", []Symbol[string]{Variable[string]{Value: "primary"}}},
			{"factor", append(append(
				[]Symbol[string]{Variable[string]{Value: "factor"}},
				StringToTerminals[string]("^")...),
				[]Symbol[string]{Variable[string]{Value: "primary"}}...)},
			{"primary", append(append(
				StringToTerminals[string]("("),
				[]Symbol[string]{Variable[string]{Value: "expr"}}...),
				StringToTerminals[string](")")...)},
			{"primary", StringToTerminals[string]("0")},
			{"primary", StringToTerminals[string]("1")},
			{"primary", StringToTerminals[string]("a")},
			{"primary", StringToTerminals[string]("b")},
			{"add_op", StringToTerminals[string]("+")},
			{"add_op", StringToTerminals[string]("-")},
			{"mul_op", StringToTerminals[string]("*")},
			{"mul_op", StringToTerminals[string]("/")},
		},
	)

	arithmeticAfterTerm = NewContextFreeGrammar[string](
		"Vstart", /* start */
		[]ProductionRule[string]{
			{"Vstart", []Symbol[string]{Variable[string]{Value: "expr"}}},
			{"expr", []Symbol[string]{Variable[string]{Value: "term"}}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "expr"},
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},
			{"term", []Symbol[string]{Variable[string]{Value: "factor"}}},
			{"term", []Symbol[string]{
				Variable[string]{Value: "term"},
				Variable[string]{Value: "mul_op"},
				Variable[string]{Value: "factor"},
			}},
			{"factor", []Symbol[string]{Variable[string]{Value: "primary"}}},
			{"factor", []Symbol[string]{
				Variable[string]{Value: "factor"},
				Variable[string]{Value: "Vpow_op"},
				Variable[string]{Value: "primary"},
			}},
			{"primary", []Symbol[string]{
				Variable[string]{Value: "Vopen"},
				Variable[string]{Value: "expr"},
				Variable[string]{Value: "Vclose"},
			}},
			{"primary", StringToTerminals[string]("0")},
			{"primary", StringToTerminals[string]("1")},
			{"primary", StringToTerminals[string]("a")},
			{"primary", StringToTerminals[string]("b")},
			{"add_op", StringToTerminals[string]("+")},
			{"add_op", StringToTerminals[string]("-")},
			{"mul_op", StringToTerminals[string]("*")},
			{"mul_op", StringToTerminals[string]("/")},
			{"Vpow_op", StringToTerminals[string]("^")},
			{"Vopen", StringToTerminals[string]("(")},
			{"Vclose", StringToTerminals[string](")")},
		},
	)

	arithmeticAfterBin = NewContextFreeGrammar[string](
		"Vstart", /* start */
		[]ProductionRule[string]{
			{"Vstart", []Symbol[string]{Variable[string]{Value: "expr"}}},
			{"expr", []Symbol[string]{Variable[string]{Value: "term"}}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "expr"},
				Variable[string]{Value: "Vadd_op_term"},
			}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},
			{"term", []Symbol[string]{Variable[string]{Value: "factor"}}},
			{"term", []Symbol[string]{
				Variable[string]{Value: "term"},
				Variable[string]{Value: "Vmul_op_factor"},
			}},
			{"factor", []Symbol[string]{Variable[string]{Value: "primary"}}},
			{"factor", []Symbol[string]{
				Variable[string]{Value: "factor"},
				Variable[string]{Value: "Vpow_op_primary"},
			}},
			{"primary", []Symbol[string]{
				Variable[string]{Value: "Vopen"},
				Variable[string]{Value: "Vexpr_close"},
			}},
			{"primary", StringToTerminals[string]("0")},
			{"primary", StringToTerminals[string]("1")},
			{"primary", StringToTerminals[string]("a")},
			{"primary", StringToTerminals[string]("b")},
			{"add_op", StringToTerminals[string]("+")},
			{"add_op", StringToTerminals[string]("-")},
			{"mul_op", StringToTerminals[string]("*")},
			{"mul_op", StringToTerminals[string]("/")},
			{"Vpow_op", StringToTerminals[string]("^")},
			{"Vopen", StringToTerminals[string]("(")},
			{"Vclose", StringToTerminals[string](")")},
			{"Vadd_op_term", []Symbol[string]{
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},
			{"Vmul_op_factor", []Symbol[string]{
				Variable[string]{Value: "mul_op"},
				Variable[string]{Value: "factor"},
			}},
			{"Vpow_op_primary", []Symbol[string]{
				Variable[string]{Value: "Vpow_op"},
				Variable[string]{Value: "primary"},
			}},
			{"Vexpr_close", []Symbol[string]{
				Variable[string]{Value: "expr"},
				Variable[string]{Value: "Vclose"},
			}},
		},
	)

	arithmeticAfterDel = arithmeticAfterBin

	arithmeticAfterUnit = NewContextFreeGrammar[string](
		"Vstart", /* start */
		[]ProductionRule[string]{
			{"Vstart", StringToTerminals[string]("0")},
			{"Vstart", StringToTerminals[string]("1")},
			{"Vstart", StringToTerminals[string]("a")},
			{"Vstart", StringToTerminals[string]("b")},
			{"Vstart", []Symbol[string]{
				Variable[string]{Value: "Vopen"},
				Variable[string]{Value: "Vexpr_close"},
			}},
			{"Vstart", []Symbol[string]{
				Variable[string]{Value: "factor"},
				Variable[string]{Value: "Vpow_op_primary"},
			}},
			{"Vstart", []Symbol[string]{
				Variable[string]{Value: "term"},
				Variable[string]{Value: "Vmul_op_factor"},
			}},
			{"Vstart", []Symbol[string]{
				Variable[string]{Value: "expr"},
				Variable[string]{Value: "Vadd_op_term"},
			}},
			{"Vstart", []Symbol[string]{
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},

			{"expr", StringToTerminals[string]("0")},
			{"expr", StringToTerminals[string]("1")},
			{"expr", StringToTerminals[string]("a")},
			{"expr", StringToTerminals[string]("b")},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "Vopen"},
				Variable[string]{Value: "Vexpr_close"},
			}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "factor"},
				Variable[string]{Value: "Vpow_op_primary"},
			}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "term"},
				Variable[string]{Value: "Vmul_op_factor"},
			}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "expr"},
				Variable[string]{Value: "Vadd_op_term"},
			}},
			{"expr", []Symbol[string]{
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},

			{"term", StringToTerminals[string]("0")},
			{"term", StringToTerminals[string]("1")},
			{"term", StringToTerminals[string]("a")},
			{"term", StringToTerminals[string]("b")},
			{"term", []Symbol[string]{
				Variable[string]{Value: "Vopen"},
				Variable[string]{Value: "Vexpr_close"},
			}},
			{"term", []Symbol[string]{
				Variable[string]{Value: "factor"},
				Variable[string]{Value: "Vpow_op_primary"},
			}},
			{"term", []Symbol[string]{
				Variable[string]{Value: "term"},
				Variable[string]{Value: "Vmul_op_factor"},
			}},

			{"factor", StringToTerminals[string]("0")},
			{"factor", StringToTerminals[string]("1")},
			{"factor", StringToTerminals[string]("a")},
			{"factor", StringToTerminals[string]("b")},
			{"factor", []Symbol[string]{
				Variable[string]{Value: "Vopen"},
				Variable[string]{Value: "Vexpr_close"},
			}},
			{"factor", []Symbol[string]{
				Variable[string]{Value: "factor"},
				Variable[string]{Value: "Vpow_op_primary"},
			}},

			{"primary", StringToTerminals[string]("0")},
			{"primary", StringToTerminals[string]("1")},
			{"primary", StringToTerminals[string]("a")},
			{"primary", StringToTerminals[string]("b")},
			{"primary", []Symbol[string]{
				Variable[string]{Value: "Vopen"},
				Variable[string]{Value: "Vexpr_close"},
			}},

			{"add_op", StringToTerminals[string]("+")},
			{"add_op", StringToTerminals[string]("-")},
			{"mul_op", StringToTerminals[string]("*")},
			{"mul_op", StringToTerminals[string]("/")},
			{"Vpow_op", StringToTerminals[string]("^")},
			{"Vopen", StringToTerminals[string]("(")},
			{"Vclose", StringToTerminals[string](")")},

			{"Vadd_op_term", []Symbol[string]{
				Variable[string]{Value: "add_op"},
				Variable[string]{Value: "term"},
			}},
			{"Vmul_op_factor", []Symbol[string]{
				Variable[string]{Value: "mul_op"},
				Variable[string]{Value: "factor"},
			}},
			{"Vpow_op_primary", []Symbol[string]{
				Variable[string]{Value: "Vpow_op"},
				Variable[string]{Value: "primary"},
			}},
			{"Vexpr_close", []Symbol[string]{
				Variable[string]{Value: "expr"},
				Variable[string]{Value: "Vclose"},
			}},
		},
	)
)
