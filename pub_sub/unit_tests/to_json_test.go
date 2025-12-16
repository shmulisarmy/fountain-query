package main

import (
	"encoding/json"
	"sql-compiler/compiler/rowType"
	pubsub "sql-compiler/pub_sub"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func GenTodoRow() gopter.Gen {
	return gopter.CombineGens(
		gen.AnyString(), // Name
		gen.Bool(),      // Age
		gen.Int(),       // Email
	).Map(func(values []any) rowType.RowType {
		return rowType.RowType{
			values[0].(string),
			values[1].(bool),
			values[2].(int),
		}
	})
}

func TestThat_Table_To_Json_Is_Valid_Json_with_property_testing(t *testing.T) {
	//a mapper should have the same amount of rows as the the table that its pointed to

	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)
	properties.Property("TestProperties", prop.ForAll(func(rows []rowType.RowType) bool {

		row_schema := rowType.RowSchema{
			rowType.ColInfo{
				Type: rowType.String,
				Name: "title",
			},
			rowType.ColInfo{
				Type: rowType.Bool,
				Name: "completed",
			},
			rowType.ColInfo{
				Type: rowType.Int,
				Name: "person_id",
			},
		}
		todo_table := pubsub.New_R_Table(row_schema)

		todo_table.Add(rowType.RowType{
			"clean the room", true, 1,
		})
		todo_table.Add(rowType.RowType{
			"clean the other room", true, 1,
		})
		todo_table.Add(rowType.RowType{
			"take out the trash", false, 1,
		})
		for _, row := range rows {
			todo_table.Add(row)
		}

		mapper := todo_table.Map_on(func(rt rowType.RowType) rowType.RowType {
			return rowType.RowType{rt[0].(string) + "!"}
		})

		json_string := pubsub.ObserverToJson(mapper, row_schema)
		var parsed_json map[string]map[string]any
		err := json.Unmarshal([]byte(json_string), &parsed_json)
		if err != nil {
			// t.Fatal(err)
			return false
		}

		if parsed_json == nil {
			// t.Fail()
			t.Log(json_string)
			return false
		}

		return true
	}, gen.SliceOf(GenTodoRow())))

	if !properties.Run(gopter.ConsoleReporter(true)) {
		t.Fail()
	}

}
