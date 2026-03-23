# API Reference

Base URL: `http://localhost:8000`

All endpoints are exposed via grpc-gateway. Request/response bodies use JSON.

---

## Racing

### ListRaces

```REST
POST /v1/list-races
```

Returns a list of races. All filter fields are optional.

**Request body:**

```json
{
  "filter": {
    "meetingIds": [1, 2],
    "visibleOnly": true
  },
  "orderBy": "advertised_start_time DESC"
}
```

- `meetingIds` -- filter by meeting IDs.
- `visibleOnly` -- when `true`, only visible races are returned.
- `orderBy` -- format: `<field> [ASC|DESC]`. Allowed fields: `advertised_start_time`, `id`, `meeting_id`, `name`, `number`, `visible`. Defaults to `advertised_start_time ASC`.

**Response:**

```json
{
  "races": [
    {
      "id": "1",
      "meetingId": "1",
      "name": "Race A",
      "number": "3",
      "visible": true,
      "advertisedStartTime": "2026-03-23T10:00:00Z",
      "status": "OPEN"
    }
  ]
}
```

`status` is a derived field: `OPEN` if `advertisedStartTime` is in the future, `CLOSED` otherwise.

---

### GetRace

```REST
GET /v1/races/{id}
```

Returns a single race. Returns `404 Not Found` if the ID does not exist.

**Response:** Same Race object as above.

---

## Sports

### ListEvents

```REST
POST /v1/list-events
```

Returns a list of sports events.

**Request body:**

```json
{
  "filter": {
    "visibleOnly": true
  }
}
```

- `visibleOnly` -- when `true`, only visible events are returned.

**Response:**

```json
{
  "events": [
    {
      "id": "1",
      "name": "Grand Final",
      "sport": "football",
      "visible": true,
      "advertisedStartTime": "2026-03-23T10:00:00Z"
    }
  ]
}
```

Results are ordered by `advertised_start_time ASC`.
