# 🐾 catApi

> **catApi** is a minimal, self-hostable API endpoint for serving pictures of, you guessed it, cats!
> but it could be used to serve anything, really

## Usage

There are two ways to "use" **catApi** - you can either serve it, blessing the world with pictures of your cat
or be served. Below is the API documentation for visiting an existing instance of catApi.

### API Documentation

**catApi** exposes several endpoints.

#### ID

`/api/id` will return the image with the associated ID.

For example **`http://localhost:3000/api/id?id=3`** will return the image with the ID of "3".

#### List

`/api/list` will return eturn a JSON object containing data about the images within the /images directory

For example, **`http://localhost:3000/api/random`** will return a JSON object that might be as follows

> `[{"id":"0","url":"/api/id?id=0"},{"id":"1","url":"/api/id?id=1"},{"id":"2","url":"/api/id?id=2"}]`

#### Random

`/api/random` will return a random image from the list of available images

### Self-hosting

TODO

## License

> **catApi** is licensed under the [MIT](https://github.com/NotAShelf/catApi/blob/v2/LICENSE) license.
