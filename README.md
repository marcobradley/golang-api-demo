# golang-api-demo

## Run tests

From the `demo-api` module directory:

```powershell
go test ./...
```

From the repository root:

```powershell
npm test
```

Verbose output from the repository root:

```powershell
npm run test:verbose
```

## POST route examples

The API currently supports two POST routes:

- `POST /songs`
- `POST /quicksort`

Base URL (local run): `http://localhost:8080`

Base URL (Docker host routing in this repo): `http://go-api.localhost:8080`

If you are using **Windows PowerShell 5.1**, use `curl.exe --%` so JSON is passed to `curl` exactly as written.

### Create a song (success)

```powershell
curl.exe --% -X POST http://localhost:8080/songs -H "Content-Type: application/json" -d "{\"title\":\"Levitating\",\"artist\":\"Dua Lipa\",\"price\":1.49}"
```

Expected status: `201 Created`

Example response:

```json
{
	"id": "4",
	"title": "Levitating",
	"artist": "Dua Lipa",
	"price": 1.49
}
```

### Invalid JSON body

```powershell
curl.exe --% -X POST http://localhost:8080/songs -H "Content-Type: application/json" -d "{\"title\":"
```

Expected status: `400 Bad Request`

Example response:

```json
{
	"message": "invalid request body"
}
```

### Omit `id` (auto-generated)

```powershell
curl.exe --% -X POST http://localhost:8080/songs -H "Content-Type: application/json" -d "{\"title\":\"No ID Song\",\"artist\":\"Unknown\",\"price\":0.99}"
```

Expected status: `201 Created`

Example response:

```json
{
	"id": "4",
	"title": "No ID Song",
	"artist": "Unknown",
	"price": 0.99
}
```

### Another create example (auto-generated `id`)

```powershell
curl.exe --% -X POST http://localhost:8080/songs -H "Content-Type: application/json" -d "{\"title\":\"Duplicate\",\"artist\":\"Tester\",\"price\":0.99}"
```

Expected status: `201 Created`

Example response:

```json
{
	"id": "4",
	"title": "Duplicate",
	"artist": "Tester",
	"price": 0.99
}
```

## Quicksort POST example

### Sort an array

```powershell
curl.exe --% -X POST http://localhost:8080/quicksort -H "Content-Type: application/json" -d "{\"array\":[5,3,8,1,2]}"
```

Expected status: `200 OK`

Example response:

```json
{
	"sorted": [1, 2, 3, 5, 8]
}
```

### Invalid JSON body

```powershell
curl.exe --% -X POST http://localhost:8080/quicksort -H "Content-Type: application/json" -d "{\"array\":"
```

Expected status: `400 Bad Request`

Example response:

```json
{
	"error": "unexpected EOF"
}
```