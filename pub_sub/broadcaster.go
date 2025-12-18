package pubsub

import (
	"sql-compiler/compiler/rowType"
	"sql-compiler/unwrap"
)

// sum
type SumBroadCaster struct {
	Observable
	Current_sum   int
	subscribed_to ObservableI
	RowSchema     unwrap.Option[rowType.RowSchema] //should be a single int
}

func (this *SumBroadCaster) set_subscribed_to(observable ObservableI) {
	this.subscribed_to = observable
}

func (this *SumBroadCaster) Pull(yield func(rowType.RowType) bool) {
	this.Current_sum = 0
	for row := range this.subscribed_to.Pull {
		this.Current_sum += row[0].(int)
	}
	if !yield(rowType.RowType{this.Current_sum}) {
		return
	}
}
func (this *SumBroadCaster) on_Add(row rowType.RowType) {
	println("yay on_Add called")
	this.Current_sum += row[0].(int)
	this.Publish_Add(rowType.RowType{this.Current_sum})
}

func (this *SumBroadCaster) on_remove(row rowType.RowType) {
	println("yay on_remove called")
	this.Current_sum -= row[0].(int)
	this.Publish_Add(rowType.RowType{this.Current_sum})
}

func (this *SumBroadCaster) on_update(old_row rowType.RowType, new_row rowType.RowType) {
	println("yay on_update called")
	this.Current_sum -= old_row[0].(int)
	this.Current_sum += new_row[0].(int)
	this.Publish_Update(rowType.RowType{this.Current_sum}, rowType.RowType{this.Current_sum})
}

func (this *SumBroadCaster) String() string {
	res := "["
	for row := range this.Pull {
		res += RowTypeToJson(&row, this.GetRowSchema()) + ","
	}
	return res + "]"
}
func (this *SumBroadCaster) GetRowSchema() rowType.RowSchema {
	return this.RowSchema.Unwrap()
}
