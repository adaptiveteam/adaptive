package queries

const StrategyStatus = `
SELECT objective.id                                                                                                   AS 'objective_id',
       initiative_updates.id                                                                                          AS 'initiative_id',
       objective.name                                                                                                 AS 'Objective Name',
       objective.type                                                                                                 AS 'Objective Type',
       DATE(objective.created_at)                                                                                     AS 'Objective Created On',
       DATE(objective_updates.db_updated_at)                                                                          AS 'Objective Updated On',
       objective.expected_end_date                                                                                    AS 'Objective End Date',
       CONCAT(ROUND(DATEDIFF(objective.expected_end_date, CURDATE()) /
                    DATEDIFF(objective.expected_end_date, objective.created_at) * 100, 0),
              '%')                                                                                                    AS 'Objective Time Left',
       objective_advocates.display_name                                                                               AS 'Objective Advocate',
       CASE
           WHEN objective_updates.status_color is null THEN 'No Status'
           WHEN objective_updates.status_color = 'Red' THEN 'Off Track'
           WHEN objective_updates.status_color = 'Yellow' THEN 'At Risk'
           WHEN objective_updates.status_color = 'Green' THEN 'On Track'
           END                                                                                                        AS 'Objective Status',
       IF(objective_updates.comments is null, 'No Updates',
          objective_updates.comments)                                                                                 AS 'Objective Update',
       objective.description                                                                                          AS 'Objective Description',
       IF(initiative_updates.name is null, 'No Initiatives',
          initiative_updates.name)                                                                                    AS 'Initiative Name',
       IF(initiative_advocates.display_name is null, 'No Initiatives',
          initiative_advocates.display_name)                                                                          AS 'Initiative Advocate',
       CASE
           WHEN initiative_updates.status_color is null THEN 'No Status'
           WHEN initiative_updates.status_color = 'Red' THEN 'Off Track'
           WHEN initiative_updates.status_color = 'Yellow' THEN 'At Risk'
           WHEN initiative_updates.status_color = 'Green' THEN 'On Track'
           END                                                                                                        AS 'Initiative Status',
       IF(initiative_updates.created_at is null, 'No Initiatives',
          DATE(initiative_updates.created_at))                                                                        AS 'Initiative Created On',
       IF(initiative_updates.db_updated_at is null, 'No Initiatives',
          DATE(initiative_updates.last_updated))                                                                      AS 'Initiative Updated On',
       IF(initiative_updates.expected_end_date is null, 'No Initiatives',
          DATE(initiative_updates.expected_end_date))                                                                 AS 'Initiative End Date',
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
           ) AS 'Initiative Time Left',
       IF(initiative_updates.comments is null, 'No Updates',
          initiative_updates.comments)                                                                                AS 'Initiative Update',
       IF(initiative_updates.description is null, 'No Initiatives',
          initiative_updates.description)                                                                             AS 'Initiative Description'

FROM user AS objective_advocates,
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
                initiative.created_at,
                initiative.db_updated_at,
                initiative.expected_end_date,
                initiative.capability_objective_id AS objective_id,
                latest_update.status_color,
                latest_update.comments,
                MAX(latest_update.created_date)    AS last_updated

         FROM initiative

                  LEFT JOIN (
             SELECT user_objective_progress.objective_id,
                    user_objective_progress.platform_id,
                    user_objective_progress.status_color,
                    user_objective_progress.comments,
                    user_objective_progress.partner_comments,
                    user_objective_progress.created_date
             FROM user_objective_progress
         ) AS latest_update ON
                 initiative.id = latest_update.objective_id AND
                 initiative.platform_id = latest_update.platform_id
         GROUP BY initiative.name
         ORDER BY initiative.name
     ) AS initiative_updates ON
         objective.id = initiative_updates.objective_id

         LEFT JOIN user AS initiative_advocates ON
         initiative_advocates.id = initiative_updates.advocate

WHERE objective.platform_id = ?
  AND objective_advocates.id = objective.advocate

ORDER BY objective.type,
         objective.name
`
