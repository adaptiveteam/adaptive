# Issue Closeout workflow

This workflow is started when the user tries to closeout an issue. 
This requires a confirmation from the coach or accounatility partner.

1. Alice tries to closeout the issue. If there is no coach - the attempt is interrupted.
2. The issue is not changed (?) or is moved to closeout-attempt state.
3. Coach receives notification (via a postponed event). Coach can approve the closeout, or
reject with a comment.
4. If closeout is approved, coachee gets a happy notification.
5. If closeout is rejected, coachee gets a notification with the coach's comment.
