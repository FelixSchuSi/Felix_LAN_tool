const path = require('path');
const express = require('express');
const exphbs = require('express-handlebars');
const app = express();
const fs = require('fs');
//HANDLEBARS
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

//reading relevant Zip-files + 7ZipSetup

let files = [];

fs.readdir(path.join(__dirname, "..", "files-to-be-transfered"), "utf-8", (err, items) => {
    console.table(items);
    filtered = items.filter(elem => (elem.endsWith('.torrent') || elem.endsWith('.zip')));
    filtered.forEach(element => {
        files.push({ name: element, dir: path.join(__dirname, "..", "files-to-be-transfered", element) })
    });
})
fs.readdir(path.join(__dirname), "utf-8", (err, items) => {
    console.table(items);
    filtered = items.filter(elem => (elem === "7zipSetup.exe"));
    filtered.forEach(element => {
        files.push({ name: element, dir: path.join(__dirname, "..", element) })
    });
})


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

app.listen(3000, () => {
    console.log('listening on port 3000!');
});
