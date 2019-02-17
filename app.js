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
app.use(express.static(__dirname + '/public'));

//DATEN WERDEN BEREITGESTELLT, zwei Alternativen wählbar. Eine muss auskommentiert werden
// ALTERNATIVE 1: NUR ZIPs UND 7zip INSTALLATION BEREITSTELLEN
let files = [];

fs.readdir(path.join(__dirname, ".."), "utf-8", (err, items) => {
    console.table(items);
    filtered = items.filter(elem => (elem === "7zipSetup.exe" || elem.endsWith('.zip')));
    filtered.forEach(element => {
        files.push({ name: element, dir: path.join(__dirname, "..", element) })
    });
})

// ACHTUNG!!! Alternative 2 funktioniert (noch) nicht.
// ALTERNATIVE 2: ALLE DATEIEN UND ORDNER BEREITSTELLEN
// fs.readdir(path.join(__dirname, ".."), "utf-8", (err, items) => {
//     items.forEach(element => {
//         files.push({ name: element, dir: path.join(__dirname, "..", element) })
//     });
// })


app.get('/', (req, res) => {
    res.render('index', { files: files });
});

app.post('/dl/:name', (req, res) => {
    let file = req.params.name;
    files.forEach(e => {
        if(e.name === file){
            file = e.dir;
        }
    })
    res.download(file);
});
app.listen(3000, () => {
    console.log('listening on port 3000!');
});
