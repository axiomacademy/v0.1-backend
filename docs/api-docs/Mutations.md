## Mutations 🧬
This section refers to the GraphQL mutations that are currently available in the backend. Mutations typically involve some change in the database state, like creating a user or scheduling a lesson. Requests that trigger some sort of business logic are also considered mutations in our system.

### `createStudent: input: NewStudent!): String!`
This mutation is fairly self-explanatory, it creates a new student. This is typically invoked during the account creation process when a new user is signing up onto the platform.

Request parameters :speaking_head: :
```
NewStudent {
  username
  firstName
  lastName
  email
  password
  profilePic
}
```

Response parameters :repeat: :
Returns a string which contains a JWT token for login purposes, and request authentication. **This has to be inserted into a cookie called `token`**

### `createTutor: input: NewTutor!): String!`
The equivalent analogue to `createStudent` but enables tutor creation instead.


Request parameters :speaking_head: :
```
NewTutor {
  username
  firstName
  lastName
  email
  password
  profilePic
  hourlyRate: An integer number with the hourly payment charge
  bio
  education: An array of strings containing their alumni institutions
  subjects: The list of subjects they teach
}
```

Response parameters :repeat: :
Returns a string which contains a JWT token for login purposes, and request authentication. **This has to be inserted into a cookie called `token`**

### `loginStudent(input: LoginInfo!): String!`
Logs in any student, based on the credentials provided in `LoginInfo`

Request parameters :speaking_head: :
```
LoginInfo {
  username
  password
}
```

Response parameters :repeat: :
Returns a string which contains a JWT token for login purposes, and request authentication. **This has to be inserted into a cookie called `token`**

### `loginTutor(input: LoginInfo!): String!`
Logs in any tutor, based on the credentials provided in `LoginInfo`

Request parameters :speaking_head: :
```
LoginInfo {
  username
  password
}
```

Response parameters :repeat: :
Returns a string which contains a JWT token for login purposes, and request authentication. **This has to be inserted into a cookie called `token`**

### `refreshToken: String!`
Takes a request authorised with a still valid JWT token, and then returns a new JWT token that expires later

Response parameters :repeat: :
Returns a string which contains a JWT token for login purposes, and request authentication. **This has to be inserted into a cookie called `token`**

### `updateHeartbeat(input: HeartbeatStatus!): String!`
The heartbeat service keeps track of which tutors are online and which of them are accepting on-demand requests. This requires the tutors to send heartbeat requests at regular intervals to keep their status online. Students are unable to access this mutation.

Request parameters :speaking_head: :
```
HeartbeatStatus{
  AVAILABLE
  UNAVAILABLE
}
```

Response parameters :repeat: :
Returns a string which contains a refreshed JWT token for login purposes, and request authentication. **This has to be inserted into a cookie called `token`**

### `createLessonRoom(input: String!): String!`
This takes in a string which is the lesson Id for the lesson you want to create the room for. After that, it returns a auth token for the room. This is only callable by Tutors, only they have the authorisation to start a lesson.

Request parameters :speaking_head: : 
Lesson Id as a string

Response parameters :repeat: :
Auth token as a string

### `endLessonRoom(input: String!): String!`
This takes in a string which is the lesson id for the lesson you want to end the room for. After that, it returns a status string. This can be called by both tutors and student.

Request parameters :speaking_head: : 
Lesson Id as a string

Response parameters :repeat: :
A status string, says `SUCCESS`

### `requestOnDemandMatch(input: OnDemandMatchRequest!): String!`
Creates an on-demand match request. Only students can make this request.

Request parameters :speaking_head: :
```
OnDemandMatchRequest {
  subject: Takes the relevant subject and subject standard the student is looking for
}
```

Response parameters :repeat: :
Returns a string which contains a match id, which can later be used by `checkForMatch` to long poll for a match

### `requestScheduledMatch(input: ScheduledMatchRequest!): String!`
Creates a scheduled match request. Only students can make this request.

Request parameters :speaking_head: :
```
ScheduledMatchRequest {
  tutor: Takes the tutor id as a string
  subject: Takes the relevant subject and subject standard the student is looking for
  time: Takes a `TimeRangeRequest`
}
```

Response parameters :repeat: :
Returns a string which contains a match id, which can later be used by `checkForMatch` to check for a match

### `acceptOnDemandMatch(input: String!): Lesson!`
Accepts an existing match, only accessible by the tutor.

Request parameters :speaking_head: :
Takes in a string containing the match id

Response parameters :repeat: :
Returns a `Lesson` with the appropriate data type

### `acceptScheduledMatch(input: String!): Lesson!`
Accepts an existing match, only accessible by the tutor.

Request parameters :speaking_head: :
Takes in a string containing the match id

Response parameters :repeat: :
Returns a `Lesson` with the appropriate data type

### `updateNotification(input: UpdateNotification!): Notification!`
Updates the notification, primarily meant to update the read status of the notification, but could be extended in the future

Request parameters :speaking_head: :
```graphql
UpdateNotification {
  id: UUID of notification
  read: boolean value indicated whether it has been read
}
```

Response parameters :repeat: :
Returns an updated `Notification`

### `registerPushNotification(input: String!): String!`
Specifically related to Firebase, it accepts the tokens generated by Firebase Messaging on the client-side to be stored and used for push notifications by the backend.

Request parameters :speaking_head: :
A string which is the Firebase APN token

Response parameters :repeat: :
Returns the same string if successful, rather pointless I know :cry: 