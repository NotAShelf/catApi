 
const express = require('express');
const app = express();
const fs = require('fs'); 
require('dotenv').config();

// Use tycrek's logging library for fancy messages
// also configure the timestamps so they don't appear
// timestamps look ugly
const TLog = require('@tycrek/log')
const logger2 = new TLog({
    timestamp: {
        enabled: false
    }
});

// Check if a custom port is set in .env
// if it is, then check if a custom url is defined. Paste the correct IP:port combo or the custom url.
// You will need to handle domain -> raw IP redirection yourself.
if(process.env.PORT) { 
    const port = process.env.PORT; // Set port variable to the value in .env
    logger2
        .success('Custom port is set. Using custom port ' + port)
    app.listen(port, () => {
            if(process.env.CUSTOMURL) {
                const customurl = process.env.CUSTOMURL;
                logger2
                    .success('Custom URL defined as ' + customurl)
                    .info('App listening at ' + 'https://' + customurl)
                    .comment('Godspeed, little fella!')
                const getRandomImageApi = () => {
                    const path = getRandomImagePath();
                        const id = path.substring(0, path.length - 4).split('-')[1];
                        return {
                            id: parseInt(id),
                            url: 'https://' + customurl + '/' + id,
                        };
                    }

                    app.get('/api/random', (req, res) => {
                        return res.json(getRandomImageApi());
                    });

                    app.get('/api/:id', (req, res) => {
                        const image = getImageById(req.params.id, false);

                        if (image) {
                            return res.json({
                                id: parseInt(req.params.id),
                                url: 'https://' + customurl + '/' + req.params.id,
                            });

                        } 
                        else {
                            return res.json({
                                error: 'Image not found'
                            });
                        }
                    });
                    
            }
            else {
                logger2
                    .error('Custom URL not defined. Falling back to localhost')
                    .info('App listening at http://localhost:' + port)
                    .comment('Godspeed, little fella!')
                const getRandomImageApi = () => {
                    const path = getRandomImagePath();
                        const id = path.substring(0, path.length - 4).split('-')[1];
                        return {
                            id: parseInt(id),
                            url: 'http://localhost:' + port + '/' + id,
                        };
                    }
                    app.get('/api/random', (req, res) => {
                        return res.json(getRandomImageApi());
                    });
                    app.get('/api/:id', (req, res) => {
                        const image = getImageById(req.params.id, false);
                        if (image) {
                            return res.json({
                                id: parseInt(req.params.id),
                                url: 'http://localhost:' + port + '/' + req.params.id,
                            });
                        } else {
                            return res.json({
                                error: 'Image not found'
                            });
                        }
                    });
                    
            }
    });
}

else { 
    const port = 3005;
    logger2
        .error('Custom port is not set. Falling back to port ' + port); 
    app.listen(port, () => {
        if(process.env.CUSTOMURL) {
            logger2
                .success('Custom URL defined. Using custom url ' + customurl)
                .success('App listening at ' + 'https://' + customurl); 
            console.log('Godspeed, little fella!');
            const getRandomImageApi = () => {
                const path = getRandomImagePath();
                    const id = path.substring(0, path.length - 4).split('-')[1];
                    return {
                        id: parseInt(id),
                        url: 'http://' + customurl + '/' + id,
                    };
                }
                app.get('/api/random', (req, res) => {
                    return res.json(getRandomImageApi());
                });
                
        }
        else {
            logger2
                .error('Custom URL not defined. Falling back to localhost')
                .info('App listening at ' + 'http://localhost:' + port)
                .comment('Godspeed, little fella!')
            const getRandomImageApi = () => {
                const path = getRandomImagePath();
                    const id = path.substring(0, path.length - 4).split('-')[1];
                    return {
                        id: parseInt(id),
                        url: 'http://localhost:' + port + '/' + id,
                    };
                }

                app.get('/api/random', (req, res) => {
                    return res.json(getRandomImageApi());
                });
                
        }
    });
}

app.get('/', (req, res) => {
    res.sendFile(getRandomImagePath(), (err) => {
        if (err) {
            res.status(err.status).end();
        }
    });
});

app.get('/:id', (req, res) => {
    const image = getImageById(req.params.id);
    if (image) {
        res.sendFile(image, (err) => {
            if (err) {
                res.status(err.status).end();
            }
        });
    } else {
        return res.json({
            error: 'Image not found'
        });
    }
});

app.get('/api/list', (req, res) => {
   return res.json(getAllImageIds());
});


const getRandomImagePath = () => {
    const images = fs.readdirSync('./images');
    return  __dirname + '/images/' + images[Math.floor(Math.random() * images.length)];
};

const getImageById = (id, path = true) => {
    const images = fs.readdirSync('./images');

    for (const image of images) {
        if (image.substring(0, image.length - 4).split('-')[1] === id) {
            if (path) {
                return __dirname + '/images/' + image;
            } else {
                return image;
            }
        }
    }

    return null;
};

const getAllImageIds = () => {
    const ids = [];

    fs.readdirSync('./images').forEach(image => {
        ids.push(parseInt(image.substring(0, image.length - 4) .split('-')[1]));
    });

    return ids.sort((a, b) => a - b);
};
