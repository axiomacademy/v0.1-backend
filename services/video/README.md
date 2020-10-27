# video service

Interacts with the [Twilio Rooms API](https://www.twilio.com/docs/video/api/rooms-resource) to create rooms and auth tokens to access those rooms.

## Endpoints

All these endpoints take a single string as an input: the id of the room in question. I advise that this id be the lesson id, but anything goes.

### `createLessonRoom(string) string`

Creates the room and then returns an auth token to the room.

### `endLessonRoom(string) string`

Ends the specified room, 

### `getLessonRoom(string) string`

Returns an auth token to the specified room. Does not return anything.

## The auth token

The auth token returned is a JWT signed with the API key and HMAC256. Being the way it is, the client is not supposed to be able to verify the token since the key also happens to be the API key. 

The token's body should have a single field:

```json
{
	"room": "<room id>"
}
```
