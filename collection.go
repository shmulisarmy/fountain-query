package main

type R_Table struct {
	Observable
	rows []RowType
}

func (this *R_Table) add(row RowType) {
	this.rows = append(this.rows, row)
	this.publish_add(row)
}
