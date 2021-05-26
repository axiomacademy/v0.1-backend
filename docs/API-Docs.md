# API Documentation 📓

Thankfully, GraphQL APIs can be introspected using the GraphQL playground. This can be done by accessing `localhost:8080` once running your local development version of the backend as instructed in [:tractor: Getting Started](getting-started). That will display the schema, plus any queries, mutations and subscriptions associated with the backend.

The purpose of this documentation is then, to outline the high-level usage details of the API, and explain how it fits within the context of the program flow. The parameters of the request will also be explained. This will broadly be subdivided into the queries, mutations and subscriptions sections as follows.

## Queries 🤔
* [`self: User!`](api-docs/Queries#self-user)
* [`lessons(input: TimeRangeRequest!): [Lesson!]`](api-docs/Queries#lessonsinput-timerangerequest-lesson)
* [`pendingMatches: [Match!]`](api-docs/Queries#pendingmatches-match)
* [`notifications(input: TimeRangeRequest!): [Notification!]!`](api-docs/Queries#notificationsinput-timerangerequest-notification)
* [`getScheduledMatches(input: ScheduledMatchParameters!): [Tutor!]!`](https://gitlab.solderneer.me/axiom/backend/-/wikis/api-docs/Queries#getscheduledmatchesinput-scheduledmatchparameters-tutor)
* [`checkForMatch(input: String!): Lesson`](api-docs/Queries#checkformatchinput-string-lesson)
* [`getLessonRoom(input: String!): String!`](api-docs/Queries#getlessonroominput-string-string)

## Mutations 🧬
* [`createStudent: input: NewStudent!): String!`](api-docs/Mutations#createstudent-input-newstudent-string)
* [`createTutor: input: NewTutor!): String!`](api-docs/Mutations#createtutor-input-newtutor-string)
* [`loginStudent(input: LoginInfo!): String!`](api-docs/Mutations#loginstudentinput-logininfo-string)
* [`loginTutor(input: LoginInfo!): String!`](api-docs/Mutations#logintutorinput-logininfo-string)
* [`refreshToken: String!`](api-docs/Mutations#refreshtoken-string)
* [`updateHeartbeat(input: HeartbeatStatus!): String!`](api-docs/Mutations#updateheartbeatinput-heartbeatstatus-string)
* [`createLessonRoom(input: String!): String!`](api-docs/Mutations#createlessonroominput-string-string)
* [`endLessonRoom(input: String!): String!`](api-docs/Mutations#endlessonroominput-string-string)
* [`requestOnDemandMatch(input: OnDemandMatchRequest!): String!`](api-docs/Mutations#requestondemandmatchinput-ondemandmatchrequest-string)
* [`requestScheduledMatch(input: ScheduledMatchRequest!): String!`](api-docs/Mutations#requestscheduledmatchinput-scheduledmatchrequest-string)
* [`acceptOnDemandMatch(input: String!): Lesson!`](api-docs/Mutations#acceptondemandmatchinput-string-lesson)
* [`acceptScheduledMatch(input: String!): Lesson!`](api-docs/Mutations#acceptscheduledmatchinput-string-lesson)
* [`updateNotification(input: UpdateNotification!): Notification!`](api-docs/Mutations#updatenotificationinput-updatenotification-notification)
* [`registerPushNotification(input: String!): String!`](api-docs/Mutations#registerpushnotificationinput-string-string)

## Subscriptions 📰
* [`subscribeMatchNotifications: MatchNotification!`](api-docs/Subscriptions#subscribematchnotifications-matchnotification)