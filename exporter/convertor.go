package exporter

type TypeConverter func(string) string

var _ TypeConverter = typescriptTypeConverter

var tsProtocolTypeConverter TypeConverter = func(o string) string {
	switch o {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "string", "decimal.Decimal":
		return "string"
	default:
		return "any"
	}
}

func typescriptTypeConverter(o string) string {
	switch o {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "string", "decimal.Decimal":
		return "string"
	default:
		return o
	}
}
