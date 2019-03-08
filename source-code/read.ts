import * as path from 'path';
import * as fs from 'fs';
import { Item } from './models/item';

export default async function read() {
    const items = await readDirAsync(path.join(__dirname, "..", "files-to-be-transfered"))
    const filtered = items.filter((elem: string) => (elem.endsWith('.zip') || elem === "7zipSetup.exe"));
    const files: Array<Item> = await filtered.map((e: string) => {
        return { name: e, dir: path.join(__dirname, "..", "files-to-be-transfered", e) };
    })
    console.table(files);
    return files;
}

function readDirAsync(path: string): Promise<Array<string>> {
    return new Promise<Array<string>>((resolve, reject) => {
        fs.readdir(path, "utf-8", (err, data: Array<any>) => {
            err ? reject(err) : resolve(data);
        });
    });
}