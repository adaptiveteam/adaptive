package queries

const SelectCompetenciesByPlatformID = `
SELECT
  competency_type as type,
  name,
  description
FROM
  competency
WHERE
  platform_id    = ? AND
  deactivated_on = ""
ORDER BY
  type,
  name
`
