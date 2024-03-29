import $file.Meta
import Meta._

import $file.Dsl
import Dsl._

lazy val deactivatedAtField = "DeactivatedAt".camel :: optionTimestampRawType

lazy val createdAtField  = ("CreatedAt" .camel :: timestampRawType) \\ "Automatically maintained field"

lazy val modifiedAtField = ("ModifiedAt".camel :: optionTimestampRawType) \\ "Automatically maintained field"

implicit class VirtualFieldsOps(entity: Entity) {

    def virtualFields: List[Field] = 
      entity.traits.flatMap{
        case DeactivationTrait =>
          List(deactivatedAtField)
        case CreatedModifiedTimesTrait => 
          List(createdAtField, modifiedAtField)
      }
}
