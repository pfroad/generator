package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseTable(t *testing.T) {
	Convey("Failed to parse table", t, func() {
		table, err := parseTable("airparking", "ap_bank", "com.airparking.cloud.ecenter", "model")
		So(err, ShouldBeNil)

		for _, col := range table.Columns {
			t.Logf("%s: %s %s", col.ColumnName, col.ShortType, col.Property)
		}
	})
}

func TestDb(t *testing.T) {
	rows, err := db.Query("SELECT count(*) FROM ap_bank")
	if err != nil {
		t.Error(err)
		return
	}

	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			t.Error(err)
			return
		}

		t.Log(count)
	}
}
