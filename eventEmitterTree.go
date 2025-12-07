package main

import (
	pubsub "sql-compiler/pub_sub"
	"sql-compiler/rowType"
)

const path_separator = "/"

type SyncType string

const (
	SyncTypeUpdate = "update"
	SyncTypeAdd    = "add"
	SyncTypeRemove = "remove"
)

type SyncMessage struct {
	Type SyncType
	Data string
	Path string
}

type eventEmitterTree struct {
	on_message func(SyncMessage)
}

func (receiver *eventEmitterTree) syncFromObservable(obs *pubsub.Mapper, path string) {
	// switch obs := obs.(type) {
	// case *pubsub.Mapper:
	// 	for row := range obs.Pull {
	// 		pubsub.RowTypeToJson(&row, obs.RowSchema.Unwrap())
	// 	}
	// default:
	// 	panic("expected mapper")
	// }
	obs.Add_sub(&pubsub.CustomSubscriber{
		OnAddFunc: func(item rowType.RowType) {
			primary_key := item[0].(string)
			receiver.on_message(SyncMessage{Type: SyncTypeAdd, Data: pubsub.RowTypeToJson(&item, obs.RowSchema.Unwrap()), Path: path + path_separator + primary_key})
			receiver.syncFromObservable_row(item, path+path_separator+primary_key, obs.RowSchema.Unwrap())
		},
		OnRemoveFunc: func(item rowType.RowType) {
			primary_key := item[0].(string)
			receiver.on_message(SyncMessage{Type: SyncTypeRemove, Data: pubsub.RowTypeToJson(&item, obs.RowSchema.Unwrap()), Path: path + path_separator + primary_key})
		},
		OnUpdateFunc: func(oldItem, newItem rowType.RowType) {
			primary_key := oldItem[0].(string)
			receiver.on_message(SyncMessage{Type: SyncTypeUpdate, Data: pubsub.RowTypeToJson(&newItem, obs.RowSchema.Unwrap()), Path: path + path_separator + primary_key})
		},
	})
	for row := range obs.Pull {
		receiver.syncFromObservable_row(row, path, obs.RowSchema.Unwrap())
	}

}

func (receiver *eventEmitterTree) syncFromObservable_row(row rowType.RowType, path string, row_schema rowType.RowSchema) {
	for i, col := range row {
		switch col := col.(type) {
		case string, int, bool:
		case pubsub.ObservableI:
			switch col := col.(type) {
			case *pubsub.Mapper:
				receiver.syncFromObservable(col, path+path_separator+row_schema[i].Name)
			default:
				panic("should be mapper")
			}
		default:
			panic("unhandled")
		}
	}
}
