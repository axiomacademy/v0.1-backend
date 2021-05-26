## Queries 🤔

This section looks purely at the available queries. In GraphQL terminology, queries are a **read-only** operations and do not modify server state (although this is entirely up to the discretion of the resolver implementation). However, we do stick to the official recommendations in our implementations.

### `self: User!`
This endpoint simply returns the data of whatever user you happen to be. If you are a `User` of type `Student` you can access all relevant student data, and similarly for `Tutor`. Since there is a form of polymorphism at play here, a note can be made about type fragments in GraphQL. Type fragements are used to parse multiple possible types in a GrapQL response as follows.

```graphql
query {
  self {
    ... on Student {
      username
    }
    ... on Tutor {
      hourlyRate
    }
  }
}
```

However, it is highly likely that the caller will always know what type of `User` they are, therefore will only need one fragment in the request, throwing an error if an invalid type is returned.

### `lessons(input: TimeRangeRequest!): [Lesson!]`
This endpoint returns the list of lessons that the user has, regardless of whether the user is a student of a tutor. The `TimeRangeRequest` is used as a pagination tool to load and cache appropriate blocks on the frontend.

Request parameters :speaking_head: :
```
TimeRangeRequest {
  startTime: Absolute start time
  endTime: Absolute end time
}
```

Response parameters :repeat: :
```graphql
Lesson {
  id: UUID of the lesson, a random string of ASCII characters
  subject: Type of string values of valid subject, standard pairs
  summary: Summary of the lesson, typically set by the tutor after a lesson
  tutor: Tutor structure that is associated with this lesson
  student: Student structure that is associated with this lesson
  scheduled: Returns true if the lesson is scheduled and false if it is an on-demand session
  startTime: Absolute start time if it is on-demand, relative start time if it is scheduled
  endTime: Absolute end time if it is on-demand, relative end time if it is scheduled
}
```
### `pendingMatches: [Match!]`
This endpoint returns the list of pending matches the student or tutor has. Pending matches are typically related to scheduled tutoring sessions. These sessions need to be negotiated between the tutor and the student. Therefore, the tutor needs access to the pending matches to accept/deny them, and the student needs to be able to view the pending matches they have requested.

Response parameters :repeat: :
```graphql
Match {
  id: UUID of the match, a random string of ASCII characters
  status: Match status, can be `MATCHED`, `MATCHING` or `FAILED`
  scheduled: Returns true if the lesson is scheduled and false if it is an on-demand session
  tutor: Tutor structure that is associated with this lesson
  student: Student structure that is associated with this lesson
  subject: Type of string values of valid subject, standard pairs
  startTime: Absolute start time if it is on-demand, relative start time if it is scheduled
  endTime: Absolute end time if it is on-demand, relative end time if it is scheduled
}
```

### `notifications(input: TimeRangeRequest!): [Notification!]!`
This endpoint gets the notifications for a user, within a time window. The time window is used as a pagination tool to load and cache appropriate blocks on the frontend.

Request parameters :speaking_head: :
```
TimeRangeRequest {
  startTime: Absolute start time
  endTime: Absolute end time
}
```

Response parameters :repeat: :
```graphql
Notification {
  id: UUID of the match, a random string of ASCII characters
  title: String title of the notificaiton
  subtitle: String subtitle of the notification
  image: String pointing to an S3 bucket for the image
  created: Absolute time of the notification's creation
}
```

### `getScheduledMatches(input: ScheduledMatchParameters!): [Tutor!]!`
This takes a scheduled tutor request and returns a list of tutors who are available to take the lesson. This is typically used in the request flow of scheduling a tutor, and is made by the student after which the student picks a specific tutor to request a match with. That match can be requested using the mutation `requestScheduledMatch`.

Request parameters :speaking_head: :
```
ScheduledMatchParameters {
  subject: The subject they want a lesson for
  time: `TimeRangeRequest` as shown above
}
```

Response parameters :repeat: :
Returns a list of `Tutor` types

### `checkForMatch(input: String!): Lesson`
This takes in an input which has the match id, and returns a lesson if it is available. Students will have to long poll this to check for matches.

Request parameters :speaking_head: : `Match` Id

Response parameters :repeat: :  A single `Lesson` type if it is successful, else an error message `No Match Found`

### `getLessonRoom(input: String!): String!`
This takes in a string which is the lesson Id for the lesson you want to open the room for. After that, it returns a auth token for the room. This is only callable by Tutors, only they have the authorisation to start a lesson.

Request parameters :speaking_head: : 
Lesson Id as a string

Response parameters :repeat: :
Auth token as a string