import $file.Meta
import Meta._

import $file.Dsl
import Dsl._

val uuid = golangStringRawType // we can use it to automatically generate id value in `create` method

val int = simpleType("Int", "int", "0", "N", emptyValueIsNormal = true)

val string = golangStringRawType // == simpleType("String", "string", "S")

val optionString = simpleType("Option[String]", "string", "\"\"", "S", isOptional = true)

val timestamp = timestampRawType

val optionTimestamp = optionTimestampRawType

val boolean = simpleType("Boolean", "bool", "false", "BOOL", emptyValueIsNormal = true)

val bool = boolean

val optionBoolean = simpleType("Option[Boolean]", "bool", "false", "BOOL", isOptional = true)

val optionStringArray = simpleType("Option[Array[String]]", "[]string", "nil", "SS")

val goTypes: Map[String, TypeInfo] = Map(
    "bool" -> boolean,
    "string" -> optionString,
    "int" -> int
)
