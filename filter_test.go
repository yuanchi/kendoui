package kendoui

import (
	"fmt"
	"encoding/json"
	"testing"
)

func TestGroupFilterUnmarshal(t *testing.T) {
	var jsonBlob = []byte(`{"logic":"and","filters":[{"operator":"startswith","value":"AAA","field":"modelId"},{"logic":"and","filters":[{"field":"orderDate","operator":"gt","value":"2018-07-29T00:00:00.000Z"},{"field":"orderDate","operator":"gt","value":"2018-07-31T00:00:00.000Z"}]}]}`)
	var groupFilter GroupFilter
	err := json.Unmarshal(jsonBlob, &groupFilter)
	if err != nil {
		t.Errorf("%s", err.Error())
	} else {
		fmt.Println("logic:", groupFilter.Logic)
		fmt.Println("filters:")
		indent := " "
		PrintResult(indent, &groupFilter.Filters)
/*
logic: and
filters:
  operator: startswith
  value: AAA
  field: modelId
  logic: and
  filters:
   operator: gt
   value: 2018-07-29T00:00:00.000Z
   field: orderDate
   operator: gt
   value: 2018-07-31T00:00:00.000Z
   field: orderDate
*/
	}
}

func genFilterBlock(sf *SingleFilter) string {
	field := sf.Field
	operator := sf.Operator
	value := sf.Value
	
	r := ""
	switch field {
	case "modelId":
		switch operator {
		case "startswith":
			r = field + " LIKE '" + value.(string) + "%'"
		}
	case "orderDate":
		oper := ""
		switch operator {
		case "gt":
			oper = ">"
		default:
			oper = "="
		}
		r = field + " " + oper + " '" + value.(string) + "'"
	}
	return r
}

func TestComposeStmt(t *testing.T) {
	var jsonBlob = []byte(`{"logic":"and","filters":[{"operator":"startswith","value":"AAA","field":"modelId"},{"logic":"and","filters":[{"field":"orderDate","operator":"gt","value":"2018-07-29T00:00:00.000Z"},{"field":"orderDate","operator":"gt","value":"2018-07-31T00:00:00.000Z"}]}]}`)
	var groupFilter GroupFilter
	err := json.Unmarshal(jsonBlob, &groupFilter)
	if err != nil {
		t.Errorf("%s", err.Error())
	} 
	s := ComposeStmt(&groupFilter, genFilterBlock)
	expected := "modelId LIKE 'AAA%' and (orderDate > '2018-07-29T00:00:00.000Z' and orderDate > '2018-07-31T00:00:00.000Z')"
	if s != expected {
		t.Error("Not matched to expected")
	}
}

func TestUnmarshalToGroupFilter(t *testing.T) {
	var jsonBlob = []byte(`{"logic":"and","filters":[{"operator":"startswith","value":"AAA","field":"modelId"},{"logic":"and","filters":[{"field":"orderDate","operator":"gt","value":"2018-07-29T00:00:00.000Z"},{"field":"orderDate","operator":"gt","value":"2018-07-31T00:00:00.000Z"}]}]}`)
	groupFilter, err := UnmarshalToGroupFilter(jsonBlob)
	if err != nil {
		t.Errorf("%s", err.Error())
	} else {
		fmt.Println("logic:", groupFilter.Logic)
		fmt.Println("filters:")
		indent := " "
		PrintResult(indent, &groupFilter.Filters)
	}
}
// custom struct implement UnmarshalJSON function: 32137 ns/op
func Benchmark_GroupFilterUnmarshal(b *testing.B) {
	var jsonBlob = []byte(`{"logic":"and","filters":[{"operator":"startswith","value":"AAA","field":"modelId"},{"logic":"and","filters":[{"field":"orderDate","operator":"gt","value":"2018-07-29T00:00:00.000Z"},{"field":"orderDate","operator":"gt","value":"2018-07-31T00:00:00.000Z"}]}]}`)
	for i := 0; i < b.N; i++ {		
		var groupFilter GroupFilter
		_ = json.Unmarshal(jsonBlob, &groupFilter)
	}
}
// first unmarshal to map, then convert to struct: 7643 ns/op
func Benchmark_UnmarshalToGroupFilter(b *testing.B) {
	var jsonBlob = []byte(`{"logic":"and","filters":[{"operator":"startswith","value":"AAA","field":"modelId"},{"logic":"and","filters":[{"field":"orderDate","operator":"gt","value":"2018-07-29T00:00:00.000Z"},{"field":"orderDate","operator":"gt","value":"2018-07-31T00:00:00.000Z"}]}]}`)
	for i := 0; i < b.N; i++ {		
		_, _ = UnmarshalToGroupFilter(jsonBlob)
	}
}

type (
	FSingleFilter struct {
		Field, Operator string
		Value interface{}
	}
	FGroupFilter struct {
		Logic string
		Filters []KendoFilterNodeComposition
	}
	KendoFilterNodeComposition struct {
		*FSingleFilter
		*FGroupFilter
	}
)	
// struct composing anonymoust struct: 8760 ns/op
func Benchmark_UnmarshalCompositionWithAnonymousStruct(b *testing.B) {
	var jsonBlob = []byte(`{"logic":"and","filters":[{"operator":"startswith","value":"AAA","field":"modelId"},{"logic":"and","filters":[{"field":"orderDate","operator":"gt","value":"2018-07-29T00:00:00.000Z"},{"field":"orderDate","operator":"gt","value":"2018-07-31T00:00:00.000Z"}]}]}`)
	for i := 0; i < b.N; i++ {
		var k KendoFilterNodeComposition
		_ = json.Unmarshal(jsonBlob, &k)
	}
}

func Benchmark_SprintfTypeCheck(b *testing.B) {
	var k KendoFilterable = &SingleFilter{}
	for i := 0; i < b.N; i++ {
		s := fmt.Sprintf("%T", k)
		if s == "SingleFilter" {}
	}
}

func Benchmark_TypeAssertion(b *testing.B) {
	var k KendoFilterable = &SingleFilter{}
	for i := 0; i < b.N; i++ {
		if _, ok := k.(*SingleFilter); ok {}
	}
}