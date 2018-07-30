package kendoui

import (
	"fmt"
	"encoding/json"
	"errors"
)

type GroupFilter struct {
	Logic string
	Filters []KendoFilterable
}

type SingleFilter struct {
	Operator string
	Field string
	Value interface{}
}

type KendoFilterable interface {
	KendoFilterNode() string
}

func (g *GroupFilter) KendoFilterNode() string {
	return "GroupFilter"
}

func (s *SingleFilter) KendoFilterNode() string {
	return "SingleFilter"
}

// ref. http://gregtrowbridge.com/golang-json-serialization-with-interfaces/
// ref. https://golang.org/pkg/encoding/json/#example_RawMessage_unmarshal
func (g *GroupFilter) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}
	
	var logic string
	err = json.Unmarshal(*objMap["logic"], &logic)
	if err != nil {
		return err
	} 
	g.Logic = logic
	
	var rawMessagesForFilters []*json.RawMessage
	err = json.Unmarshal(*objMap["filters"], &rawMessagesForFilters)
	if err != nil {
		return err
	}
	
	g.Filters = make([]KendoFilterable, len(rawMessagesForFilters))
	
	var m map[string]interface{}
	for index, rawMessage := range rawMessagesForFilters {
		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}
		
		if _, ok := m["logic"]; ok {
			var gf GroupFilter
			err := json.Unmarshal(*rawMessage, &gf)
			if err != nil {
				return err
			}
			g.Filters[index] = &gf
			continue
		}
		
		if _, ok := m["field"]; ok {
			var sf SingleFilter
			err := json.Unmarshal(*rawMessage, &sf)
			if err != nil {
				return err
			}
			g.Filters[index] = &sf
			continue
		}
		
		return errors.New("Unsupported type found!")
	}
	return nil
}

func PrintResult(indent string, filters *[]KendoFilterable) {
	for _, f := range *filters {
		switch v := f.(type) {
		case *SingleFilter:
			fmt.Println(indent, "operator:", v.Operator)
			fmt.Println(indent, "value:", v.Value)
			fmt.Println(indent, "field:", v.Field)
		case *GroupFilter:
			fmt.Println(indent, "logic:", v.Logic)
			fmt.Println(indent, "filters:")
			i := indent + indent
			PrintResult(i, &v.Filters)		
		}
				
	}
}