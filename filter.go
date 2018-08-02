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

type groupFilter struct {
	Logic string
	Filters []*json.RawMessage
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
	var lgf groupFilter
	err := json.Unmarshal(b, &lgf)
	if err != nil {
		return err
	}
	 
	g.Logic = lgf.Logic
	g.Filters = make([]KendoFilterable, len(lgf.Filters))
	
	var m map[string]interface{}
	for index, rawMessage := range lgf.Filters {
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

func ComposeStmt(gf *GroupFilter, genFunc func(*SingleFilter) string) string {
	logic := gf.Logic
	filters := gf.Filters
	s := ""
	for _, f := range filters {
		if sf, ok := f.(*SingleFilter); ok {
			s += (" " + logic + " " + genFunc(sf))
			continue
		}
		
		if gfs, ok := f.(*GroupFilter); ok {
			s += (" " + logic + " (" + ComposeStmt(gfs, genFunc) + ")")
			continue
		}
	}
	s = s [(len(logic))+2:]
	return s
}

func ToKendoFilterable(m map[string]interface{}) (KendoFilterable, error) {
	var kf KendoFilterable
	if v, ok := m["logic"]; ok {
		filters := m["filters"].([]interface{})
		gf := GroupFilter{Logic: v.(string), Filters: make([]KendoFilterable, len(filters))}
		for i, f := range filters {
			if filter, ok := f.(map[string]interface{}); ok {
				out, err := ToKendoFilterable(filter)
				if err != nil {
					return kf, err
				}
				gf.Filters[i] = out
			}
		}
		kf = &gf
		return kf, nil
	}
	if v, ok := m["field"]; ok {
		sf := SingleFilter{Field: v.(string), Operator: m["operator"].(string), Value: m["value"]}
		kf = &sf
		return kf, nil
	}
	return kf, errors.New("not found the exact fields")
}

func UnmarshalToGroupFilter(b []byte) (*GroupFilter, error) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	out, err := ToKendoFilterable(m)
	if err != nil {
		return nil, err
	}
	return out.(*GroupFilter), nil
}