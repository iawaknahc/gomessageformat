package messageformat

import (
	"encoding/json"
	// "io/ioutil"
	"reflect"
	"testing"
)

func parse(t *testing.T, s string, expected []Node) {
	actual, err := Parse(s)
	if err != nil {
		t.Errorf("err: %v\n", err)
	} else if !reflect.DeepEqual(actual, expected) {
		actualBytes, _ := json.MarshalIndent(actual, "", "  ")
		expectedBytes, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("expected\n")
		t.Errorf("%s\n", string(expectedBytes))
		t.Errorf("actual\n")
		t.Errorf("%s\n", string(actualBytes))
		// ioutil.WriteFile("expected.json", expectedBytes, 0644)
		// ioutil.WriteFile("actual.json", actualBytes, 0644)
	}
}

func TestParseSelect(t *testing.T) {
	parse(t, "hello {1} {gender, select, male {gentleman} female {lady} other {{kind, select, other{}}}} ", []Node{
		TextNode{"hello "},
		NoneArgNode{Arg: Argument{Index: 1}},
		TextNode{" "},
		SelectArgNode{
			Arg: Argument{Name: "gender"},
			Clauses: []SelectClause{
				SelectClause{
					Keyword: "male",
					Nodes: []Node{
						TextNode{"gentleman"},
					},
				},
				SelectClause{
					Keyword: "female",
					Nodes: []Node{
						TextNode{"lady"},
					},
				},
				SelectClause{
					Keyword: "other",
					Nodes: []Node{
						TextNode{},
						SelectArgNode{
							Arg: Argument{Name: "kind"},
							Clauses: []SelectClause{
								SelectClause{
									Keyword: "other",
									Nodes:   []Node{TextNode{}},
								},
							},
						},
						TextNode{},
					},
				},
			},
		},
		TextNode{" "},
	})
}

func TestParsePlural(t *testing.T) {
	parse(t, "Hello {count, plural, other{}}", []Node{
		TextNode{"Hello "},
		PluralArgNode{
			Arg:    Argument{Name: "count"},
			Kind:   "plural",
			Offset: 0,
			Clauses: []PluralClause{
				PluralClause{
					Keyword: "other",
					Nodes:   []Node{TextNode{}},
				},
			},
		},
		TextNode{""},
	})

	parse(t, "Hello {count, plural, offset:1 other{}}", []Node{
		TextNode{"Hello "},
		PluralArgNode{
			Arg:    Argument{Name: "count"},
			Kind:   "plural",
			Offset: 1,
			Clauses: []PluralClause{
				PluralClause{
					Keyword: "other",
					Nodes:   []Node{TextNode{}},
				},
			},
		},
		TextNode{""},
	})

	parse(t, "Hello {count, plural, offset:1 =0{} other{}}", []Node{
		TextNode{"Hello "},
		PluralArgNode{
			Arg:    Argument{Name: "count"},
			Kind:   "plural",
			Offset: 1,
			Clauses: []PluralClause{
				PluralClause{
					ExplicitValue: 0,
					Nodes:         []Node{TextNode{}},
				},
				PluralClause{
					Keyword: "other",
					Nodes:   []Node{TextNode{}},
				},
			},
		},
		TextNode{""},
	})

	parse(t, "Hello {count, plural, offset:1 =0{} one{{gender, select, other{}}} other{}}", []Node{
		TextNode{"Hello "},
		PluralArgNode{
			Arg:    Argument{Name: "count"},
			Kind:   "plural",
			Offset: 1,
			Clauses: []PluralClause{
				PluralClause{
					ExplicitValue: 0,
					Nodes:         []Node{TextNode{}},
				},
				PluralClause{
					Keyword: "one",
					Nodes: []Node{
						TextNode{},
						SelectArgNode{
							Arg: Argument{Name: "gender"},
							Clauses: []SelectClause{
								SelectClause{
									Keyword: "other",
									Nodes:   []Node{TextNode{}},
								},
							},
						},
						TextNode{},
					},
				},
				PluralClause{
					Keyword: "other",
					Nodes:   []Node{TextNode{}},
				},
			},
		},
		TextNode{""},
	})
}

func TestParseExample1(t *testing.T) {
	parse(t, `{gender_of_host, select,
  female {{
      num_guests, plural, offset:1
      =0 {{host} does not give a party.}
      =1 {{host} invites {guest} to her party.}
      =2 {{host} invites {guest} and one other person to her party.}
      other {{host} invites {guest} and # other people to her party.}}}
  male {{
      num_guests, plural, offset:1
      =0 {{host} does not give a party.}
      =1 {{host} invites {guest} to his party.}
      =2 {{host} invites {guest} and one other person to his party.}
      other {{host} invites {guest} and # other people to his party.}}}
  other {{
      num_guests, plural, offset:1
      =0 {{host} does not give a party.}
      =1 {{host} invites {guest} to their party.}
      =2 {{host} invites {guest} and one other person to their party.}
      other {{host} invites {guest} and # other people to their party.}}}}`, []Node{
		TextNode{},
		SelectArgNode{
			Arg: Argument{Name: "gender_of_host"},
			Clauses: []SelectClause{
				SelectClause{
					Keyword: "female",
					Nodes: []Node{
						TextNode{},
						PluralArgNode{
							Arg:    Argument{Name: "num_guests"},
							Kind:   "plural",
							Offset: 1,
							Clauses: []PluralClause{
								PluralClause{
									ExplicitValue: 0,
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" does not give a party."},
									},
								},
								PluralClause{
									ExplicitValue: 1,
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" invites "},
										NoneArgNode{Arg: Argument{Name: "guest"}},
										TextNode{" to her party."},
									},
								},
								PluralClause{
									ExplicitValue: 2,
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" invites "},
										NoneArgNode{Arg: Argument{Name: "guest"}},
										TextNode{" and one other person to her party."},
									},
								},
								PluralClause{
									Keyword: "other",
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" invites "},
										NoneArgNode{Arg: Argument{Name: "guest"}},
										TextNode{" and "},
										PoundNode{},
										TextNode{" other people to her party."},
									},
								},
							},
						},
						TextNode{},
					},
				},
				SelectClause{
					Keyword: "male",
					Nodes: []Node{
						TextNode{},
						PluralArgNode{
							Arg:    Argument{Name: "num_guests"},
							Kind:   "plural",
							Offset: 1,
							Clauses: []PluralClause{
								PluralClause{
									ExplicitValue: 0,
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" does not give a party."},
									},
								},
								PluralClause{
									ExplicitValue: 1,
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" invites "},
										NoneArgNode{Arg: Argument{Name: "guest"}},
										TextNode{" to his party."},
									},
								},
								PluralClause{
									ExplicitValue: 2,
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" invites "},
										NoneArgNode{Arg: Argument{Name: "guest"}},
										TextNode{" and one other person to his party."},
									},
								},
								PluralClause{
									Keyword: "other",
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" invites "},
										NoneArgNode{Arg: Argument{Name: "guest"}},
										TextNode{" and "},
										PoundNode{},
										TextNode{" other people to his party."},
									},
								},
							},
						},
						TextNode{},
					},
				},
				SelectClause{
					Keyword: "other",
					Nodes: []Node{
						TextNode{},
						PluralArgNode{
							Arg:    Argument{Name: "num_guests"},
							Kind:   "plural",
							Offset: 1,
							Clauses: []PluralClause{
								PluralClause{
									ExplicitValue: 0,
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" does not give a party."},
									},
								},
								PluralClause{
									ExplicitValue: 1,
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" invites "},
										NoneArgNode{Arg: Argument{Name: "guest"}},
										TextNode{" to their party."},
									},
								},
								PluralClause{
									ExplicitValue: 2,
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" invites "},
										NoneArgNode{Arg: Argument{Name: "guest"}},
										TextNode{" and one other person to their party."},
									},
								},
								PluralClause{
									Keyword: "other",
									Nodes: []Node{
										TextNode{},
										NoneArgNode{Arg: Argument{Name: "host"}},
										TextNode{" invites "},
										NoneArgNode{Arg: Argument{Name: "guest"}},
										TextNode{" and "},
										PoundNode{},
										TextNode{" other people to their party."},
									},
								},
							},
						},
						TextNode{},
					},
				},
			},
		},
		TextNode{},
	})
}

func TestPound(t *testing.T) {
	parse(t, "{0, select, a{#} b{{1, plural, one{{0, select, other{#}}} other{#}}} other{#}}", []Node{
		TextNode{},
		SelectArgNode{
			Arg: Argument{Index: 0},
			Clauses: []SelectClause{
				SelectClause{
					Keyword: "a",
					Nodes: []Node{
						TextNode{"#"},
					},
				},
				SelectClause{
					Keyword: "b",
					Nodes: []Node{
						TextNode{},
						PluralArgNode{
							Arg:    Argument{Index: 1},
							Kind:   "plural",
							Offset: 0,
							Clauses: []PluralClause{
								PluralClause{
									Keyword: "one",
									Nodes: []Node{
										TextNode{},
										SelectArgNode{
											Arg: Argument{Index: 0},
											Clauses: []SelectClause{
												SelectClause{
													Keyword: "other",
													Nodes: []Node{
														TextNode{"#"},
													},
												},
											},
										},
										TextNode{},
									},
								},
								PluralClause{
									Keyword: "other",
									Nodes: []Node{
										TextNode{},
										PoundNode{},
										TextNode{},
									},
								},
							},
						},
						TextNode{},
					},
				},
				SelectClause{
					Keyword: "other",
					Nodes: []Node{
						TextNode{"#"},
					},
				},
			},
		},
		TextNode{},
	})
}

func TestParseDatetime(t *testing.T) {
	parse(t, "hello {t, date, short} {t, time, medium} {t, datetime, long} {t, datetime, full}", []Node{
		TextNode{"hello "},
		DateArgNode{Arg: Argument{Name: "t"}, Style: "short"},
		TextNode{" "},
		TimeArgNode{Arg: Argument{Name: "t"}, Style: "medium"},
		TextNode{" "},
		DatetimeArgNode{Arg: Argument{Name: "t"}, Style: "long"},
		TextNode{" "},
		DatetimeArgNode{Arg: Argument{Name: "t"}, Style: "full"},
		TextNode{""},
	})
}
