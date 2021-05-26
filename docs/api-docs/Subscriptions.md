## Subscriptions 📰
This section covers the subscriptions. These are mainly used to deliver SSE kind of event-based trigger communication. In Axiom, the match notifications play such a role.

### `subscribeMatchNotifications: MatchNotification!`
Allows tutors to subscribe to match notifications so that they are informed of any on-demand matches when they are online and available.

Response parameters :repeat: :
```graphql
MatchNotification {
  student: Returns a `Student` for the match
  subject: Returns the subject the student is interested in learning
  token: Returns a token used for the matching in `acceptOnDemandMatch`
}
```