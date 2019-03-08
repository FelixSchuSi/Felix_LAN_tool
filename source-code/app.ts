import * as path from 'path';
import express, { Express } from 'express';
import exphbs from 'express-handlebars';
import { Item } from './models/item';
import read from './read';

const port: number = 3000;

function configureApp(app: Express) {
    const engineConfig = {
        extname: '.hbs',
        partialsDir: path.join(__dirname, 'views', 'partials'),
        layoutsDir: path.join(__dirname, 'views', 'layouts'),
        defaultLayout: 'main'
    }
    app.engine('hbs', exphbs(engineConfig));
    app.set('view engine', 'hbs');
    app.set('views', path.join(__dirname, '/views'));
    app.use(express.static(path.join(__dirname + '/public')));
}

async function start() {
    const app: Express = express();
    let files: Array<Item> = await read();
    configureApp(app);
    startHttpServer(app, files);
}

function startHttpServer(app: Express, files: Array<Item>) {
    // rendering standard view
    app.get('/', (req, res) => {
        res.render('index', { files: files });
    });

    // Download handler
    app.post('/dl/:name', (req, res) => {
        let file = req.params.name;
        files.forEach(e => {
            if (e.name === file) {
                file = e.dir;
            }
        })
        res.download(file);
    });

    app.listen(port, () => {
        console.log(`Server running at http://localhost:${port}`);
    });
}

start();