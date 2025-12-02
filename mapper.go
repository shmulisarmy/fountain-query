package main

type Mapper struct {
	Observable
	transformer   func(RowType) RowType
	subscribed_to ObservableI
}

func (this *Mapper) set_subscribed_to(observable ObservableI) {
	this.subscribed_to = observable
}

func (this *Mapper) on_add(row RowType) {
	this.publish_add(this.transformer(row))
}

func (this *Mapper) on_remove(row RowType) {
	this.publish_remove(this.transformer(row))
}

func (this *Mapper) on_update(old_row RowType, new_row RowType) {
	this.publish_publish(this.transformer(old_row), this.transformer(new_row))
}
