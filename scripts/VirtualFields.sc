import $file.Meta
import Meta._

import $file.Dsl
import Dsl._

lazy val deactivatedOnField = "DeactivatedOn".camel :: golangStringRawType

lazy val createdAtField  = ("CreatedAt" .camel :: timestampRawType) \\ "Automatically maintained field"

lazy val modifiedAtField = ("ModifiedAt".camel :: timestampRawType) \\ "Automatically maintained field"

implicit class VirtualFieldsOps(entity: Entity) {

    def virtualFields: List[Field] = 
      entity.traits.flatMap{
        case DeactivationTrait =>
          List(deactivatedOnField)
        case CreatedModifiedTimesTrait => 
          List(createdAtField, modifiedAtField)
      }
}
