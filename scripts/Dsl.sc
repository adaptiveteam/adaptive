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

def defaultPackage(table: Table, imports: Imports): Package = 
  Package(table.entity.name, 
    List(
      daoModule(table, imports),
      fieldNamesModule(table, imports)
    )
  )

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

def fieldNamesModule(table: Table, imports: Imports): Module = 
  Module(Filename(table.entity.name ++ "Names".camel, ".go"), 
    List(GoModulePart(
      List(),//imports.importClauses,
      List(
        StringBasedEnum("FieldName".camel, table.entity.fields.map(fieldToStringBasedEnumItem))
      )
    ))
  )


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
