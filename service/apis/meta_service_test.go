package apis

import (
	"encoding/json"
	"testing"

	"github.com/intuit/katlas/service/db"
	"github.com/stretchr/testify/assert"
)

func TestMetaService(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")

	q := NewQueryService(dc)
	m := NewMetaService(dc)
	e := NewEntityService(dc)
	// create pod metadata
	podMeta := `{
		"name": "pod_metadata",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldName": "name",
				"fieldType": "json",
				"mandatory": true,
				"index": true,
				"cardinality": "One",
                                "upsert": true,
                                "tokenizer": [
                                        "term",
                                        "trigram"
                                ]
			},
			{
				"fieldName": "status",
				"fieldType": "string",
				"mandatory": true,
				"index": false,
				"cardinality": "One"
			},
			{
				"fieldName": "containers",
				"fieldType": "relationship",
				"refDataType": "K8scontainer",
				"mandatory": false,
				"index": false,
				"cardinality": "Many"
			}
		]
	}`
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "name", PType: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", PType: "string", Index: true, Tokenizer: []string{"term"}})
	// create pod metadata
	dataMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(podMeta), &dataMap)
	if err != nil {
		panic(err)
	}
	m.CreateMetadata(dataMap)
	// query to get created pod metadata
	qm := map[string][]string{"name": {"pod_metadata"}, "objtype": {"metadata"}}
	n, _ := q.GetQueryResult(qm)
	o := n["objects"].([]interface{})[0].(map[string]interface{})
	// cleanup after test
	defer e.DeleteEntity(o["uid"].(string))
	assert.Equal(t, o["name"], "pod_metadata", "query return doesn't match pod_metadata")
	for _, fields := range o["fields"].([]interface{}) {
		rid := fields.(map[string]interface{})["uid"]
		defer e.DeleteEntity(rid.(string))
	}
	// get all fields
	fs, err := m.GetMetadataFields("pod_metadata")
	if err != nil {
		assert.Fail(t, "Failed to get meta fields")
	}
	assert.Equal(t, 3, len(fs), "return fields don't match metadata define")
}
