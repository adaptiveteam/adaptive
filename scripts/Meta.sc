
sealed trait Name

case class SimpleName(parts: List[String]) extends Name

case class QualifiedName(packageName: Name, name: Name) extends Name

object Name {
  def apply(parts: List[String]): SimpleName = SimpleName(parts)

}

sealed trait TypeInfo

case class SimpleTypeInfo(name: String,
  golangType: Name,
  defaultValueLiteral: String, // go value that is used by default
  dynamoType: String,
  /*sqlType: String, */
  isOptional: Boolean = false,
  emptyValueIsNormal: Boolean = false) extends TypeInfo

case class StructTypeInfo(name: Name) extends TypeInfo

object StructTypeInfo {
  def apply(entity: Entity): StructTypeInfo = StructTypeInfo(entity.name)
}

case class TypeAliasTypeInfo(ta: TypeAlias) extends TypeInfo

case class Field(name: SimpleName, tpe: TypeInfo, comments: List[String] = Nil, 
  dynamoDbFieldRename: Option[SimpleName] = None) {
    def dbName: SimpleName = dynamoDbFieldRename.getOrElse(name)

}

sealed trait Trait

// Deactivation adds a field "DeactivatedOn: timestamp" and
// changes "delete" logic to update with timestamp.
// Also all reads will include filter of only active rows.
case object DeactivationTrait extends Trait

// CreatedModifiedTimesTrait adds a couple of fields - 
// created_at: timestamp
// modified_at: timestamp
// and respective logic to maintain these field values.
case object CreatedModifiedTimesTrait extends Trait

case class Entity(
  name: SimpleName,
  primaryKeyFields: List[Field],
  otherFields: List[Field],
  comments: List[String] = Nil,
  traits: List[Trait] = Nil
) {
  def fields: List[Field] = primaryKeyFields ::: otherFields
}

// Part of a go module. A few parts can be merged together
case class GoModulePart(importClauses: List[ImportClause], definitions: List[GoDefinition])

trait ConvertibleToGoDefinitions[A] {
  def convertToGoDefinitions(a: A): GoModulePart
}

implicit class ConvertibleOps[A: ConvertibleToGoDefinitions](a: A) {
  def toGo: GoModulePart = 
    implicitly[ConvertibleToGoDefinitions[A]].convertToGoDefinitions(a)
}

sealed trait Definition

sealed trait GoDefinition extends Definition

case class TypeAlias(name: SimpleName, tpe: TypeInfo) extends GoDefinition {
  def getType: TypeInfo = TypeAliasTypeInfo(this)
}
case class Const(name: SimpleName, tpe: TypeInfo, value: String) extends GoDefinition

case class StringBasedEnumItem(name: SimpleName, value: String)
// defines a type alias and a few valid values in `const`
case class StringBasedEnum(name: SimpleName, values: List[StringBasedEnumItem]) extends GoDefinition

case class Struct(entity: Entity) extends GoDefinition
// A list of constants with the single `const` header
// case class ConstBlock(consts: List[Const]) extends GoDefinition

case class ImportClause(nameOpt: Option[String], url: String) {
  def name: String = nameOpt.getOrElse{
    val idx = url.lastIndexOf('/')
    if(idx < 0) url 
    else url.substring(idx + 1, url.length)
  }
}

case class Imports(importClauses: List[ImportClause]) extends GoDefinition

/**
 * @param suffix is dot + extension ('.go', '.tf')
 */
case class Filename(name: SimpleName, suffix: String)

case class Module(filename: Filename, parts: List[GoModulePart]) {
  def imports: Imports =
    Imports(parts.flatMap(_.importClauses).distinct)
  def definitions: List[GoDefinition] = 
    parts.flatMap(_.definitions)
}

case class Package(name: SimpleName, modules: List[Module])

sealed trait ProjectionType

object ProjectionType {
  object ALL extends ProjectionType
  object INCLUDE extends ProjectionType
}

case class Index(hashKey: Field, rangeKey: Option[Field], 
  projectionType: ProjectionType = ProjectionType.ALL, 
  nonKeyAttributes: List[Field] = Nil
) {
  def fields: List[Field] =
    hashKey :: rangeKey.toList

  def name: SimpleName = 
    SimpleName(fields.flatMap(_.dbName.parts) ::: List("index"))
}

case class Table(
  entity: Entity,
  defaultIndex: Index,
  indices: List[Index],
  encrypted: Boolean = false
) {
    def hashKey: Field = defaultIndex.hashKey
    
    def rangeKey: Option[Field] = defaultIndex.rangeKey

    def allIndices: List[Index] = defaultIndex :: indices

    def fieldsInIndices: List[Field] = {
      val set = allIndices.flatMap(_.fields.map(_.name)).toSet
      (entity.fields.filter(f => set.contains(f.name))).distinct
    }
}

sealed trait DaoOperation

case object DaoCreateRow extends DaoOperation
case object DaoReadRow extends DaoOperation
case object DaoReadChildren extends DaoOperation
case object DaoReadOrEmptyRow extends DaoOperation
case object DaoUpdateRow extends DaoOperation
case object DaoDeleteRow extends DaoOperation
case class DaoQueryRow(index: Index) extends DaoOperation
case class DaoQueryRowByHashKey(index: Index) extends DaoOperation

case class Dao(table: Table) extends GoDefinition {
    def operations: List[DaoOperation] = List(
        DaoCreateRow,
        DaoReadRow,
        DaoReadOrEmptyRow,
        DaoUpdateRow,
        DaoDeleteRow,
        // DaoQueryRow(table.defaultIndex)
    ) ::: table.indices.map(DaoQueryRow)
}

case class ConnectionBasedDao(table: Table) extends GoDefinition {
    def operations: List[DaoOperation] = List(
        DaoCreateRow,
        DaoReadRow,
        DaoReadOrEmptyRow,
        DaoUpdateRow,
        DaoDeleteRow,
        // DaoQueryRow(table.defaultIndex)
    ) ::: (
      table.defaultIndex.rangeKey. // if there is a range key, then we should also add a method to read by only hash key
        map(_ => DaoReadChildren).
        toList
    ) ::: table.indices.map(DaoQueryRow) ::: (
      table.indices.
        filter(_.rangeKey.isDefined).
        filterNot(_.hashKey == table.defaultIndex.hashKey). // because we already generate the query in ReadChildren
        groupBy(_.hashKey).
        map{ case (_, indices) => 
          DaoQueryRowByHashKey(indices.head) 
        }.
        toList
    )
}

case class SourceFile(path: String)

case class FileWithContent(sf: SourceFile, content: String)

type ProjectContent = List[FileWithContent] // (fileContentMap:

sealed trait ProjectFolder

case class GoProjectFolder(path: String, packages: List[Package]) extends ProjectFolder

case class TerraformProjectFolder(path: String, tables: List[Table]) extends ProjectFolder

type Workspace = List[ProjectFolder]
