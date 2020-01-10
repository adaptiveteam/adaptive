package queries

const Competencies = `
select

competency_type as type,
name,
description

from

competency

where

platform_id = ? AND
deactivated_on = ""

ORDER BY

type,
name
`
