package parser

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"strings"
)

type column struct {
	Property  string
	TypeName  string
	ShortType string
	JavaLang  bool

	ColumnName   string
	DataType     string
	Comment      string
	IsPrimaryKey bool
	IsNullable   bool
}

func newColumn(columnName,
	dataType,
	comment,
	colType string,
	primaryKey,
	nullable bool) (*column, error) {
	javaType, shortType, err := transType(dataType)
	if err != nil {
		return nil, err
	}

	if colType == "tinyint(1)" {
		javaType, shortType = "java.lang.Boolean", "Boolean"
	}

	return &column{
		Property:     strcase.ToLowerCamel(columnName),
		TypeName:     javaType,
		ShortType:    shortType,
		JavaLang:     strings.HasPrefix(javaType, "java.lang"),
		ColumnName:   columnName,
		DataType:     dataType,
		Comment:      comment,
		IsPrimaryKey: primaryKey,
		IsNullable:   nullable,
	}, nil
}

/* mybatis typeHandler
类型处理器	Java 类型	JDBC 类型
BooleanTypeHandler	java.lang.Boolean, boolean	数据库兼容的 BOOLEAN
ByteTypeHandler	java.lang.Byte, byte	数据库兼容的 NUMERIC 或 BYTE
ShortTypeHandler	java.lang.Short, short	数据库兼容的 NUMERIC 或 SMALLINT
IntegerTypeHandler	java.lang.Integer, int	数据库兼容的 NUMERIC 或 INTEGER
LongTypeHandler	java.lang.Long, long	数据库兼容的 NUMERIC 或 BIGINT
FloatTypeHandler	java.lang.Float, float	数据库兼容的 NUMERIC 或 FLOAT
DoubleTypeHandler	java.lang.Double, double	数据库兼容的 NUMERIC 或 DOUBLE
BigDecimalTypeHandler	java.math.BigDecimal	数据库兼容的 NUMERIC 或 DECIMAL
StringTypeHandler	java.lang.String	CHAR, VARCHAR
ClobReaderTypeHandler	java.io.Reader	-
ClobTypeHandler	java.lang.String	CLOB, LONGVARCHAR
NStringTypeHandler	java.lang.String	NVARCHAR, NCHAR
NClobTypeHandler	java.lang.String	NCLOB
BlobInputStreamTypeHandler	java.io.InputStream	-
ByteArrayTypeHandler	byte[]	数据库兼容的字节流类型
BlobTypeHandler	byte[]	BLOB, LONGVARBINARY
DateTypeHandler	java.util.Date	TIMESTAMP
DateOnlyTypeHandler	java.util.Date	DATE
TimeOnlyTypeHandler	java.util.Date	TIME
SqlTimestampTypeHandler	java.sql.Timestamp	TIMESTAMP
SqlDateTypeHandler	java.sql.Date	DATE
SqlTimeTypeHandler	java.sql.Time	TIME
ObjectTypeHandler	Any	OTHER 或未指定类型
EnumTypeHandler	Enumeration Type	VARCHAR 或任何兼容的字符串类型，用以存储枚举的名称（而不是索引值）
EnumOrdinalTypeHandler	Enumeration Type	任何兼容的 NUMERIC 或 DOUBLE 类型，存储枚举的序数值（而不是名称）。
SqlxmlTypeHandler	java.lang.String	SQLXML
InstantTypeHandler	java.time.Instant	TIMESTAMP
LocalDateTimeTypeHandler	java.time.LocalDateTime	TIMESTAMP
LocalDateTypeHandler	java.time.LocalDate	DATE
LocalTimeTypeHandler	java.time.LocalTime	TIME
OffsetDateTimeTypeHandler	java.time.OffsetDateTime	TIMESTAMP
OffsetTimeTypeHandler	java.time.OffsetTime	TIME
ZonedDateTimeTypeHandler	java.time.ZonedDateTime	TIMESTAMP
YearTypeHandler	java.time.Year	INTEGER
MonthTypeHandler	java.time.Month	INTEGER
YearMonthTypeHandler	java.time.YearMonth	VARCHAR 或 LONGVARCHAR
JapaneseDateTypeHandler	java.time.chrono.JapaneseDate	DATE
*/

func transType(dataType string) (string, string, error) {
	switch dataType {
	case "bigint":
		return "java.lang.Long", "Long", nil
	case "char", "varchar", "text", "nvarchar", "nchar", "mediumtext", "json":
		return "java.lang.String", "String", nil
	case "blob":
		return "java.lang.Byte", "Byte[]", nil
	case "int", "tinyint", "smallint", "mediumint":
		return "java.lang.Integer", "Integer", nil
	case "float":
		return "java.lang.Float", "Float", nil
	case "decimal":
		return "java.math.BigDecimal", "BigDecimal", nil
	case "double":
		return "java.lang.Double", "Double", nil
	case "time", "datetime", "timestamp":
		return "java.util.Date", "Date", nil
	case "bit", "boolean":
		return "java.lang.Boolean", "Boolean", nil
	default:
		return "", "", fmt.Errorf("unsupport data type %s", dataType)
	}
}

func isPrimaryKey(colKey string) bool {
	return colKey == "PRI"
}

func isNullable(nullable string) bool {
	return nullable == "YES"
}
