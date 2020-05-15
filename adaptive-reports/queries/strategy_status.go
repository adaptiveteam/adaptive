package queries

const SelectStrategyStatusByPlatformID = `
SELECT objective.id                                      AS 'objective_id',
       initiative_updates.id                             AS 'initiative_id',
       objective.name                                    AS 'Objective Name',
       IF(initiative_updates.name is null, 'No Initiatives',
          initiative_updates.objective_community)        AS 'Objective Community',
       objective.type                                    AS 'Objective Type',
       DATE(objective.created_at)                        AS 'Objective Created On',
       DATE(objective_updates.db_updated_at)             AS 'Objective Updated On',
       objective.expected_end_date                       AS 'Objective End Date',
       IF(user_objective.completed = 1, 'TRUE', 'FALSE') AS 'Objective Completed',
       CONCAT(ROUND(DATEDIFF(objective.expected_end_date, CURDATE()) /
                    DATEDIFF(objective.expected_end_date, objective.created_at) * 100, 0),
              '%')                                       AS 'Objective Time Left',
       objective_advocates.display_name                  AS 'Objective Advocate',
       CASE
           WHEN objective_updates.status_color is null THEN 'No Status'
           WHEN objective_updates.status_color = 'Red' THEN 'Off Track'
           WHEN objective_updates.status_color = 'Yellow' THEN 'At Risk'
           WHEN objective_updates.status_color = 'Green' THEN 'On Track'
           END                                           AS 'Objective Status',
       IF(objective_updates.comments is null, 'No Updates',
          objective_updates.comments)                    AS 'Objective Update',
       objective.description                             AS 'Objective Description',
       IF(initiative_updates.name is null, 'No Initiatives',
          initiative_updates.name)                       AS 'Initiative Name',
       IF(initiative_updates.initiative_community is null, 'No Initiatives',
          initiative_updates.initiative_community)       AS 'Initiative Community',
       IF(initiative_advocates.display_name is null, 'No Initiatives',
          initiative_advocates.display_name)             AS 'Initiative Advocate',
       CASE
           WHEN initiative_updates.status_color is null THEN 'No Status'
           WHEN initiative_updates.status_color = 'Red' THEN 'Off Track'
           WHEN initiative_updates.status_color = 'Yellow' THEN 'At Risk'
           WHEN initiative_updates.status_color = 'Green' THEN 'On Track'
           END                                           AS 'Initiative Status',
       IF(initiative_updates.created_at is null, 'No Initiatives',
          DATE(initiative_updates.created_at))           AS 'Initiative Created On',
       IF(initiative_updates.db_updated_at is null, 'No Initiatives',
          DATE(initiative_updates.last_updated))         AS 'Initiative Updated On',
       IF(initiative_updates.expected_end_date is null, 'No Initiatives',
          DATE(initiative_updates.expected_end_date))    AS 'Initiative End Date',
       IF(initiative_updates.expected_end_date is null, 'No Initiatives',
          CONCAT(
                  ROUND(
                                  DATEDIFF(
                                          initiative_updates.expected_end_date,
                                          CURDATE()
                                      ) /
                                  DATEDIFF(
                                          initiative_updates.expected_end_date,
                                          initiative_updates.created_at
                                      ) * 100, 0
                      ),
                  '%'
              )
           )                                             AS 'Initiative Time Left',
       IF(initiative_updates.comments is null, 'No Updates',
          initiative_updates.comments)                   AS 'Initiative Update',
       IF(initiative_updates.description is null, 'No Initiatives',
          initiative_updates.description)                AS 'Initiative Description'

FROM user AS objective_advocates,
     user_objective,
     objective

         LEFT JOIN (
         SELECT user_objective_progress.created_date AS db_updated_at,
                user_objective_progress.objective_id,
                user_objective_progress.status_color,
                user_objective_progress.comments

         FROM user_objective_progress,
              objective

         WHERE user_objective_progress.objective_id = objective.id
           AND user_objective_progress.platform_id = objective.platform_id
     ) AS objective_updates ON
         objective_updates.objective_id = objective.id

         LEFT JOIN (
         SELECT initiative.id,
                initiative.name,
                initiative.description,
                initiative.advocate,
                initiative_community.name          as initiative_community,
                initiative.created_at,
                initiative.db_updated_at,
                initiative.expected_end_date,
                initiative.capability_objective_id AS objective_id,
                objective_community.name           as objective_community,
                latest_update.status_color,
                latest_update.comments,
                latest_update.created_date         as last_updated

         FROM initiative_community,
              objective_community,
              initiative

                  LEFT JOIN (SELECT a.*
                             FROM user_objective_progress a
                                      LEFT OUTER JOIN user_objective_progress b
                                                      ON SUBSTRING_INDEX(a.id, ':', 1) = SUBSTRING_INDEX(b.id, ':', 1) AND
                                                         a.created_date < b.created_date
                             WHERE b.id IS NULL) as latest_update ON
                      initiative.id = latest_update.objective_id AND
                      initiative.platform_id = latest_update.platform_id
         WHERE initiative.initiative_community_id = initiative_community.id
           AND initiative_community.capability_community_id = objective_community.id
     ) AS initiative_updates ON
         objective.id = initiative_updates.objective_id

         LEFT JOIN user AS initiative_advocates ON
         initiative_advocates.id = initiative_updates.advocate

WHERE objective.platform_id = ?
  AND objective_advocates.id = objective.advocate
  AND objective.id = user_objective.id

ORDER BY user_objective.completed ASC,
         objective.type,
         objective.name
`
