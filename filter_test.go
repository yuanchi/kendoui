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