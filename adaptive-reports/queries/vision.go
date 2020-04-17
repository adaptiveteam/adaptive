package queries

const SelectVisionByPlatformID = `
SELECT
  vision.vision,
  user.display_name AS "owner",
  client_config.platform_org AS "company_name"
FROM
  user,
  vision LEFT OUTER JOIN client_config 
    ON client_config.platform_id = vision.platform_id
WHERE
  vision.platform_id = ? AND
  user.id = vision.advocate
`
