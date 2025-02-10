# 🐾 catApi

catApi is a minimal, self-hostable API (and frontend) for serving pictures of...
you guessed, it cats! Who doesn't like cats?

## API Documentation

There are several API endpoints that you may query.

### `/api/id`

`/api/id` will return the image with the associated ID.

For example **`http://localhost:3000/api/id?id=3`** will return the image with
the ID of "3".

### `/api/list`

`/api/list` will return return a JSON object containing data about the images
within the /images directory

For example, querying **`http://localhost:3000/api/random`** will return a JSON
object that might be as follows

```json
[
  { "filename": "0.jpg", "id": "0", "url": "/api/id?id=0" },
  { "filename": "1.jpg", "id": "1", "url": "/api/id?id=1" },
  { "filename": "10.jpg", "id": "2", "url": "/api/id?id=2" },
  { "filename": "11.jpg", "id": "3", "url": "/api/id?id=3" }
]
```

### `/api/random`

`/api/random` will return a random image from the list of available images

## License

**catApi** is licensed under the [MIT License](./LICENSE)
