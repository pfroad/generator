package parser

import (
	"github.com/smartystreets/goconvey/convey"
	"os"
	"strings"
	"testing"
)

func TestCurrentFile(t *testing.T) {
	convey.Convey("test current file is failed", t, func() {
		fileDir := CurrentFile()
		t.Log(fileDir[:strings.LastIndex(fileDir, string(os.PathSeparator))])
		convey.So(fileDir[strings.LastIndex(fileDir, string(os.PathSeparator)):], convey.ShouldEqual, "/util_test.go")
	})
}

func TestGenJavaModel(t *testing.T) {
	convey.Convey("failed to generate java file", t, func() {
		var err error
		var t *table
		t, err = parseTable("airparking", "ap_lease", "com.airparking.cloud.ecenter", "model")
		convey.So(err, convey.ShouldBeNil)

		err = GenJavaFile(t, "com.airparking.cloud.pcenter", "price", "", "")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestGenForJava(t *testing.T) {
	convey.Convey("failed to generate java files", t, func() {
		err := GenForJava("airparking", "u_invoice_info", "com.airparking.cloud.old")
		convey.So(err, convey.ShouldBeNil)
	})
}
