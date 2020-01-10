package queries

const Vision = `
select

vision.vision,
user.display_name as "owner",
client_config.platform_org as "company_name"

from

vision,
user,
client_config

where

vision.platform_id = ? AND
client_config.platform_id = vision.platform_id AND
user.id = vision.advocate
`
