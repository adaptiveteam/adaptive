# Request coach workflow

This workflow is triggered when an IDO user selects a coach. This will
create a `PostponedEvent` with the link to this workflow (with the id of 
IDO to coach).

This workflow starts and then
 - checks if the IDO already has coach, then just logs about it and stops.
 - if there is no coach, shows the coaching invitation to the requested user.
   In the invitation there are Accept/Reject buttons.
 - if Accepted, 
   - check again if the IDO already has coach. If the coach is different, then
     report to the current user that IDO has already got another coach and stops.
   - if there is no coach, assign the coach to IDO and notify IDO owner.
