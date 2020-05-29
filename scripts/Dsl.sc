import $file.Meta
import Meta._

import $file.Templates
import Templates._

lazy val golangStringRawType = simpleType("String", "string", "\"\"", "S")

lazy val timestampRawType = simpleType("Timestamp", "string", "\"\"", "S")

lazy val optionTimestampRawType = simpleType("Option[Timestamp]", "string", "\"\"", "S", isOptional = true)

//def tpe[T<:ClassTag](): TypeInfo = TypeInfo
def spacedName(name: String): SimpleName = SimpleName(name.split(" ").toList)

// Cortesy of https://stackoverflow.com/a/7594052/2646454 by @NPE
val camelPartsRegex = "(?<!(^|[A-Z]))(?=[A-Z])|(?<!^)(?=[A-Z][a-z])".r

def camelName(name: String): SimpleName = SimpleName(camelPartsRegex.split(name).map(_.toLowerCase).toList)

def underscoredName(name: String): SimpleName = SimpleName(name.split("_").toList)

def simpleName(name: String): SimpleName = SimpleName(List(name))

def simpleType(name: String, golangType: String, defaultValueLiteral: String, dynamoType: String, isOptional: Boolean = false,
  emptyValueIsNormal: Boolean = false): TypeInfo = 
  SimpleTypeInfo(name, simpleName(golangType), defaultValueLiteral, dynamoType, isOptional, emptyValueIsNormal: Boolean)

implicit class SourceFileOps(sf: SourceFile) {
  def :=(content: String): FileWithContent = FileWithContent(sf, content)
}

implicit class SimpleNameOps(name: SimpleName) {
  def ++(other: SimpleName): SimpleName = 
    SimpleName(name.parts ++ other.parts)

  def init: SimpleName = SimpleName(name.parts.init)

  def ^^(value: String): StringBasedEnumItem = 
    StringBasedEnumItem(name, value)
}

implicit class EntityOps(entity: Entity) {
  def parentField: Field = 
    Field(entity.name, StructTypeInfo(entity), Nil)

  def \\(comment: String): Entity = 
    entity.copy(comments = entity.comments :+ comment)

  def supports(tr: Trait): Boolean = 
    entity.traits.contains(tr)

}

def defaultPackage(table: Table, imports: Imports): Package = {
  Package(table.entity.name, 
    List(
      daoModule(table, imports),
      daoConnectionModule(table, removeUnusedImportForConnection(table, imports)),
      fieldNamesModule(table)
    )
  )
}

def daoModule(table: Table, imports: Imports): Module = 
  Module(Filename(table.entity.name, ".go"), 
    List(GoModulePart(
      imports.importClauses,
      List(
        Struct(table.entity),
        Dao(table)
      )
    ))
  )

def fieldToStringBasedEnumItem(field: Field): StringBasedEnumItem = 
  StringBasedEnumItem(field.name, snakeCaseName(field.dbName))
def indexToStringBasedEnumItem(index: Index): StringBasedEnumItem = 
  StringBasedEnumItem(index.name, goPublicName(index.name))

def fieldNamesModule(table: Table): Module = 
  Module(Filename(table.entity.name ++ "Names".camel, ".go"), 
    List(GoModulePart(
      List(),// no imports are needed for string constants
      List(
        List(StringBasedEnum("FieldName".camel, table.entity.fields.map(fieldToStringBasedEnumItem))),
        table.indices.headOption.toList.map(_ => 
          StringBasedEnum("IndexName".camel, table.indices.map(indexToStringBasedEnumItem))
        )
      ).flatten
    ))
  )

lazy val awsUtilsUrl = "github.com/adaptiveteam/adaptive/aws-utils-go"

@deprecated("Create valid imports from the very beginning", "2020-05-25")
def removeUnusedImportForConnection(table: Table, imports: Imports): Imports = {
  val unusedImports = List(
    "encoding/json", 
    "strings", 
    "github.com/adaptiveteam/adaptive/engagement-builder/model",
  )
  val importClauses = imports.importClauses.filterNot(ic => unusedImports.contains(ic.url))
  val importClauses2 = if(table.indices.isEmpty) 
    importClauses.filterNot(_.url == awsUtilsUrl) 
  else
    importClauses
  Imports(importClauses2)
}
def daoConnectionModule(table: Table, imports: Imports): Module = {
  Module(Filename(table.entity.name ++ "ConnectionBased".camel, ".go"), 
    List(GoModulePart(
      imports.importClauses,
      List(
        ConnectionBasedDao(table)
      )
    ))
  )
}

def goFieldParser(goTypes: Map[String, TypeInfo])(fieldDeclaration: String): Field = {
    try {
        val parts = fieldDeclaration.split("\\w").toList.filterNot(_.isEmpty)
        val List(name, tpe) = parts
        Field(camelName(name), goTypes.getOrElse(tpe, 
            throw new IllegalArgumentException("Couldn't parse go lang declaration" + parts.toString)), Nil)
    } catch {
        case t: Throwable =>
            println(t)
            Field(camelName("unk"), goTypes("int"), Nil)
    }
}

implicit class TypeInfoOps(tpe: TypeInfo){
    def ::(name: SimpleName): Field = Field(name, tpe, Nil)
}

implicit class FieldOps(fld: Field) {
    def \\(comment: String): Field = 
      fld.copy(comments = fld.comments :+ comment)
    def dbName(name: SimpleName): Field =
      fld.copy(dynamoDbFieldRename = Some(name))
}

implicit class StringOps(str: String) {
    def camel: SimpleName = camelName(str)
}

implicit class ImportsOps(imports: Imports) {
    def ::(imp: ImportClause): Imports = Imports(imp :: imports.importClauses)
    def :::(imps: List[ImportClause]): Imports = Imports(imps ::: imports.importClauses)
}

implicit class ImportClauseOps(importClause: ImportClause) {
    def within(n: Name): QualifiedName = QualifiedName(Name(List(importClause.name)), n)

    def struct(n: Name): TypeInfo = StructTypeInfo(within(n))

    def simpleType(name: String, golangTypeName: Name, defaultValueLiteral: String, dynamoType: String, isOptional: Boolean = false,
      emptyValueIsNormal: Boolean = false): TypeInfo = 
        SimpleTypeInfo(name, within(golangTypeName), defaultValueLiteral: String, dynamoType, isOptional, emptyValueIsNormal)
}

// implicit class PackageOps(p: Package) {
//     def ::(definition: GoDefinition): Package = p.copy(definitions = p.declarations)
// }

implicit class ModuleOps(m: Module) {
    def ::(definition: GoDefinition): Module = m.copy(parts = GoModulePart(List(), List(definition)) :: m.parts)
}

implicit class EnumOps(en: StringBasedEnum) {
  def typeAlias: TypeAlias = TypeAlias(en.name, golangStringRawType)
  def typeAliasTypeInfo: TypeAliasTypeInfo = TypeAliasTypeInfo(typeAlias)
}
