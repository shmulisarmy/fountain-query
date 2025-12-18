package main

import (
	"fmt"
	"sql-compiler/compiler/rowType"
	compiler_runtime "sql-compiler/compiler/runtime"
	"sql-compiler/db_tables"
	"sql-compiler/display"
	event_emitter_tree "sql-compiler/eventEmitterTree"
	"sql-compiler/local_live_db"
	pubsub "sql-compiler/pub_sub"
	"sql-compiler/unwrap"
	"testing"
	"time"
)

func TestPlayground(t *testing.T) {
	query := `SELECT person.name, (select Sum(person.age) from todo where todo.person_id == person.id) from person `

	obs := compiler_runtime.Query_to_observer(query)
	obs.To_display(unwrap.None[rowType.RowSchema]())
	display.DisplayStruct(obs)
	live_db := local_live_db.LocalLiveDB{}
	tree := event_emitter_tree.EventEmitterTree{
		On_message: func(message event_emitter_tree.SyncMessage) {
			display.DisplayStruct(message)
			live_db.HandleUpdate(message)
		},
	}
	tree.SyncFromObservable(obs, "")

	tree.On_message(event_emitter_tree.SyncMessage{
		Type:      event_emitter_tree.LoadInitialData,
		Data:      pubsub.ObserverToJson(obs, obs.GetRowSchema()),
		Path:      "",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
	})
	{
		tables := db_tables.Tables
		tables.Get("person").Insert(rowType.RowType{"teddy", "teddyemail@gmail.com", 22, "state", tables.Get("person").Next_row_id(), "https://api.dicebear.com/7.x/avataaars/svg?seed=teddy"})
		tables.Get("person").Insert(rowType.RowType{"ariana", "arianaemail@gmail.com", 23, "state", tables.Get("person").Next_row_id(), "https://api.dicebear.com/7.x/avataaars/svg?seed=ajay"})
		tables.Get("person").Insert(rowType.RowType{"the-doo-er", "the-doo-eremail@gmail.com", 25, "state", tables.Get("person").Next_row_id(), "https://api.dicebear.com/7.x/avataaars/svg?seed=the-doo-er"})
		tables.Get("todo").Insert(rowType.RowType{"write tests", "description", true, 1, true, tables.Get("todo").Next_row_id()})
		tables.Get("todo").Insert(rowType.RowType{"write code", "description", true, 1, true, tables.Get("todo").Next_row_id()})
		tables.Get("todo").Insert(rowType.RowType{"take a nap", "description", true, 1, false, tables.Get("todo").Next_row_id()})
		tables.Get("todo").Insert(rowType.RowType{"write code", "description", true, 2, true, tables.Get("todo").Next_row_id()})
		tables.Get("todo").Insert(rowType.RowType{"write tests", "description", true, 2, true, tables.Get("todo").Next_row_id()})
		tables.Get("todo").Insert(rowType.RowType{"take a nap", "description", true, 2, false, tables.Get("todo").Next_row_id()})

	}
	fmt.Printf("%v\n", live_db.Data)

}
