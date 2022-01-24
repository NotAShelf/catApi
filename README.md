# 🐾 catApi
> **catApi** is *scuffed* website (fully equipped with an API, obviously) to display pictures of cats from the `/images/` directory 
which are supplied by the kind people who sent me pictures of their cats, and also my own. More cats are always welcome in the paw zone.
###### I was initially going to name this project "Onlypaws" but I am too lazy to implement a donation system.


## Self-Hosting
If you are, for some reason, interested in self-hosting your own version of this website; feel absolutely free to. Fork it, download it, mirror it. Not just that I *don't care* what you do with the code, I will also 
appreciate any kind of contributions that you might decide to make. I have been told multiple times that my code is crude and scuffed, so feel free to fix that. Good luck, though.

### Dependencies
* [git](https://git-scm.com/downloads)
* [NodeJS](https://nodejs.org/en/download/)
* npm or [Yarn](https://classic.yarnpkg.com/lang/en/docs/install/) (I recommend Yarn)

### With Yarn/npm

```sh
# Clone this git repository 
git clone https://github.com/NotAShelf/catApi.git

# Move to the new directory
cd catApi

# Install dependencies
yarn 

# Alternatively, to install dependencies
npm install

# Create the .env file for environmental variables (Optional, there are fallback options for both environmental variables)
cp .sample.env .env

# Start the application with
yarn start

# Alternatively, to start the application
npm run start
```

### With Docker Compose
TODO

## Versioning
catApi is a new idea I am planning to work on (while the idea, in general, is not new, it is for me in terms of working on a project. Leave me alone.) 
and thus, we are currently on the V1 branch. V2 will *probably* introduce proper FrontEnd because let's admit it, a bare website that displays 
**literally a picture** and nothing else looks bad. When, you may ask? Whenever the heck I want, I answer.

## Contributing
TODO

## API Documentation (WIP)
* **`/api/id`**
> return data about the specified image, by ID 
> - **e.g. `http://localhost:3000/api/3`**
* **`/api/list`**
> return a list of all images on the server (explicitly from the `/images` directory)
> - e.g. `http://localhost:3000/api/list`
* **`/api/random`** 
> returns data about a random image on the server
> - e.g `http://localhost:3000/api/random`

## License

**catApi** is licensed under the [MIT](https://github.com/NotAShelf/catApi/blob/v1/LICENSE) license.
