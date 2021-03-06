package main

import (
	"errors"
	"fmt"

	elastic "github.com/danlocke/elastic-go"
)

type Decision struct {

	Id string `json:"id"`

	CaseName string `json:"name"`

	DateFiled string `json:"date_filed,omitempty"`

	Html string `json:"html,omitempty"`

	Relevance string `json:"relevance,omitempty"`

	Stored bool `json:"stored,omitempty"`

	// PlainText string `json:"plain_text"`
}

type ApiCaseResponse struct {

	Score float64 `json:"score"`

	Id string `json:"id"`

	CaseName string `json:"case_name"`

	DateFiled string `json:"date_filed,omitempty"`

	Html string `json:"html,omitempty"`

	Relevance string `json:"relevance"`

	Stored bool `json:"stored"`

}

type ApiSearchResponse struct {

	TotalHits int

	Results []ApiCaseResponse
}

type ApiGetResponse Decision

func parseDecisionFromMap(m map[string]interface{}) (Decision, error) {
	dec := Decision{}
	var ok bool
	for k, v := range m {
		switch k {
			case "id":
				f, ok := v.(float64)
				dec.Id = fmt.Sprintf("%f", f)
				if !ok {
					return Decision{}, errors.New("Could not pass Id as int64")
				}
			case "name":
				dec.CaseName, ok = v.(string)
				if !ok {
					return Decision{}, errors.New("Could not pass Case name as string")
				}
			case "date_filed":
				dec.DateFiled, ok = v.(string)
				if !ok {
					return Decision{}, errors.New("Could not pass Date filed as string")
				}
			case "html":
				dec.Html, ok = v.(string)
				if !ok {
					return Decision{}, errors.New("Could not pass html as string")
				}
		}
	}
	return dec, nil
}

func elasticGetToApiResponse(s *elastic.GetResponse) (*ApiGetResponse, error) {
	res, err := parseDecisionFromMap(s.Source.(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	r := ApiGetResponse(res)

	return &r, nil
}

// func (i *Instance) elasticSearchToApiCaseResponse(userId int64, topicId string, s *elastic.SearchResponse) (*SearchResponse, error) {
// 	res := make([]ApiCaseResponse, 0)

// 	assessed, err := dbGetAssessedPerTopic(i.db, userId, topicId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for j := range s.Hits.Hits {
// 		hit, err := parseDecisionFromMap(s.Hits.Hits[j].Source.(map[string]interface{}))
// 		if err != nil {
// 			return nil, err
// 		}
// 		stored := false
// 		relevance := ""
// 		if k, ok := assessed[s.Hits.Hits[j].Id]; ok {
// 			relevance = k
// 			stored = true
// 		}
// 		res = append(res, ApiCaseResponse {
// 			Score : s.Hits.Hits[j].Score,
// 			Id : s.Hits.Hits[j].Id,
// 			CaseName : hit.CaseName,
// 			DateFiled : hit.DateFiled,
// 			Html : hit.Html,
// 			Stored: stored,
// 			Relevance: relevance,
// 		})
// 	}

// 	r := SearchResponse {
// 		TotalHits: int(s.Hits.Total),
// 		Results: res,
// 	}

// 	return &r, nil
// }

func (i *Instance) elasticSearchToApiSearchResponse(userId int64, topicId string, s *elastic.SearchResponse) (*ApiSearchResponse, error) {
	res := make([]ApiCaseResponse, 0)

	assessed, err := dbGetAssessedPerTopic(i.db, userId, topicId)
	if err != nil {
		return nil, err
	}

	for j := range s.Hits.Hits {
		hit, err := parseDecisionFromMap(s.Hits.Hits[j].Source.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		stored := false
		relevance := ""
		if k, ok := assessed[s.Hits.Hits[j].Id]; ok {
			relevance = k
			stored = true
		}
		res = append(res, ApiCaseResponse {
			Score : s.Hits.Hits[j].Score,
			Id : s.Hits.Hits[j].Id,
			CaseName : hit.CaseName,
			DateFiled : hit.DateFiled,
			Html : hit.Html,
			Stored: stored,
			Relevance: relevance,
		})
	}

	r := ApiSearchResponse {
		TotalHits: int(s.Hits.Total),
		Results: res,
	}

	return &r, nil
}

// func elasticSearchToApiQueryResponse(query []byte, s *elastic.SearchResponse) (*ApiQueryResponse, error) {
// 	res := ApiQueryResponse{
// 		Query : string(query),
// 		TotalHits: s.Hits.Total,
// 		Hits : make([]ApiCaseResponse, len(s.Hits.Hits)),
// 	}
// 	for i := range s.Hits.Hits {
// 		hit, err := parseDecisionFromMap(s.Hits.Hits[i].Source.(map[string]interface{}))
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		res.Hits[i] = ApiCaseResponse {
// 			Score : s.Hits.Hits[i].Score,
// 			Id : s.Hits.Hits[i].Id,
// 			CaseName : hit.CaseName,
// 			DateFiled : hit.DateFiled,
// 			Html : hit.Html,
// 			Excerpt : s.Hits.Hits[i].Highlights.Highlight,
// 		}
//
// 	}
// 	return &res, nil
// }

// TODO - parsing structs from generic map[string]interefaces
// const trpTag string = "trp"

// type structFieldMap map[string]struct{
// 	Index int
// 	Type reflect.Type
// }
//
// // Adapted from https://stackoverflow.com/questions/26744873/converting-map-to-struct
// // changing map[string]interface{} to struct ... and encoding/json marshal
// func generateStructTagFieldMap(obj interface{}) (structFieldMap, error) {
// 	tagFieldIndex := map[string]struct{
// 		Index int
// 		Type reflect.Type
// 	}{}
// 	structVal := reflect.ValueOf(obj).Elem()
// 	// fmt.Println(structVal)
// 	for i := 0; i < structVal.NumField(); i++ {
// 		tagFieldIndex[structVal.Type().Field(i).Tag.Get("json")] = struct{Index int; Type reflect.Type}{
// 			Index: i,
// 			Type: structVal.Field(i).Type(),
// 		}
// 	}
// 	return tagFieldIndex, nil
// }
//
// func (s structFieldMap)structFromMap(obj map[string]interface{}, res interface{}) error {
// 	resStruct := reflect.ValueOf(res).Elem()
// 	for k, v := range obj {
// 		fmt.Println("Here -", k, v)
// 		if i, ok := s[k]; !ok {
// 			return fmt.Errorf("Error parsing to struct. No struct attribute %s", k)
// 		} else {
// 			resStructVal := resStruct.Field(i.Index)
// 			fmt.Println("Here", resStructVal.Type())
// 			if !resStructVal.CanSet() {
// 				return fmt.Errorf("Cannot set %s field value", k)
// 			}
// 			structFieldType := resStructVal.Type()
// 			val := reflect.ValueOf(v)
// 			if structFieldType != val.Type() {
// 				fmt.Println(val.Type())
// 				return fmt.Errorf("Provided value type didn't match obj field type")
// 			}
// 			resStructVal.Set(val)
// 		}
// 	}
// 	return nil
// }
