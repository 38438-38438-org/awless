package display

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/rdf"
	"github.com/wallix/awless/revision"
)

func TestDisplayCommit(t *testing.T) {
	now := time.Now()
	t0 := parseTriple(`/region<eu-west-1>	"has_type"@[]	"/region"^^type:text`)
	t1 := parseTriple(`/instance<inst_1>	"has_type"@[]	"/instance"^^type:text`)
	t2 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Id","Value":"inst_1"}"^^type:text`)
	t3 := parseTriple(`/instance<inst_2>	"has_type"@[]	"/instance"^^type:text`)
	t4 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_1>`)
	t5 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_2>`)
	graph := rdf.NewGraphFromTriples([]*triple.Triple{t0, t1, t2, t3, t4, t5})
	diff := rdf.NewEmptyDiffFromGraph(graph)
	diff.AddDeleted(t2, rdf.ParentOfPredicate)
	diff.AddDeleted(t3, rdf.ParentOfPredicate)
	diff.AddDeleted(t5, rdf.ParentOfPredicate)
	t6 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Id","Value":"new_id"}"^^type:text`)
	t7 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_3>`)
	t8 := parseTriple(`/instance<inst_3>	"has_type"@[]	"/instance"^^type:text`)
	t9 := parseTriple(`/instance<inst_3>	"property"@[]	"{"Key":"Id","Value":"inst_3"}"^^type:text`)
	diff.AddInserted(t6, rdf.ParentOfPredicate)
	diff.AddInserted(t7, rdf.ParentOfPredicate)
	diff.AddInserted(t8, rdf.ParentOfPredicate)
	diff.AddInserted(t9, rdf.ParentOfPredicate)
	revisionDiff := revision.CommitDiff{Time: now, Commit: "Commit msg", GraphDiff: diff}

	rootNode, err := node.NewNodeFromStrings("/region", "eu-west-1")
	if err != nil {
		t.Fatal(err)
	}
	table, err := TableFromBuildCommit(&revisionDiff, rootNode)
	if err != nil {
		t.Fatal(err)
	}
	var print bytes.Buffer
	table.Fprint(&print)
	expected := `+-----------+---------+----------+--------+
|  TYPE ▲   | NAME/ID | PROPERTY | VALUE  |
+-----------+---------+----------+--------+
| /instance | inst_1  |          |        |
+           +         +----------+--------+
|           |         | Id       | inst_1 |
+           +         +          +--------+
|           |         |          | new_id |
+           +---------+----------+--------+
|           | inst_2  |          |        |
+           +---------+----------+--------+
|           | inst_3  |          |        |
+-----------+---------+----------+--------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
}

func parseTriple(s string) *triple.Triple {
	t, err := triple.Parse(s, literal.DefaultBuilder())
	if err != nil {
		panic(err)
	}

	return t
}
