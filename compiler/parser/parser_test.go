package parser

import (
	"sql-compiler/compare"
	"sql-compiler/compiler/ast"
	"sql-compiler/compiler/parser/tokenizer"
	"testing"
)

func TestParser(t *testing.T) {
	src := `SELECT person.name, person.email, person.id, (
		SELECT todo.title as epic_title, person.name as author, person.id FROM todo WHERE todo.person_id == person.id
		) as todos FROM person WHERE person.age >= 3 `
	l := tokenizer.NewLexer(src)
	p := Parser{Tokens: l.Tokenize()}
	expected := &ast.Select{
		Table: "person",
		Wheres: []ast.Where{
			{
				Value1:   ast.Table_access{Table_name: "person", Col_name: "age"},
				Operator: tokenizer.GE,
				Value2:   3,
			},
		},
		Selected_values: []ast.Selected_value{
			{Value_to_select: ast.Table_access{Table_name: "person", Col_name: "name"}},
			{Value_to_select: ast.Table_access{Table_name: "person", Col_name: "email"}},
			{Value_to_select: ast.Table_access{Table_name: "person", Col_name: "id"}},
			{
				Value_to_select: ast.Select{
					Table: "todo",
					Wheres: []ast.Where{
						{
							Value1:   ast.Table_access{Table_name: "todo", Col_name: "person_id"},
							Operator: tokenizer.EQ,
							Value2:   ast.Table_access{Table_name: "person", Col_name: "id"},
						},
					},
					Selected_values: []ast.Selected_value{
						{Value_to_select: ast.Table_access{Table_name: "todo", Col_name: "title"}, Alias: "epic_title"},
						{Value_to_select: ast.Table_access{Table_name: "person", Col_name: "name"}, Alias: "author"},
						{Value_to_select: ast.Table_access{Table_name: "person", Col_name: "id"}},
					},
				},
				Alias: "todos",
			},
		},
	}
	got := p.Parse_Select()
	output, err := compare.Compare(expected, got, "")
	println(output)
	if err != nil {
		t.Fatal(err)
	}

}
