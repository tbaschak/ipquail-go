import { copyFile, mkdir, readFile, readdir, rm, writeFile } from "node:fs/promises";
import path from "node:path";

const root = path.resolve(import.meta.dirname, "..");
const publicDir = path.join(root, "ipquail", "public");
const cssDir = path.join(publicDir, "css");
const jsDir = path.join(publicDir, "js", "vendor");

const readJSON = async (file) => JSON.parse(await readFile(file, "utf8"));
const bootstrap = await readJSON(path.join(root, "node_modules", "bootstrap", "package.json"));
const jquery = await readJSON(path.join(root, "node_modules", "jquery", "package.json"));

const bootstrapCSS = `bootstrap-${bootstrap.version}.min.css`;
const bootstrapMap = `${bootstrapCSS}.map`;
const jqueryJS = `jquery-${jquery.version}.min.js`;
const jqueryMap = `jquery-${jquery.version}.min.map`;

await mkdir(cssDir, { recursive: true });
await mkdir(jsDir, { recursive: true });

const removeStale = async (directory, pattern, keep) => {
  for (const entry of await readdir(directory)) {
    if (pattern.test(entry) && !keep.has(entry)) {
      await rm(path.join(directory, entry));
    }
  }
};

await removeStale(
  cssDir,
  /^bootstrap-.*\.min\.css(?:\.map)?$/,
  new Set([bootstrapCSS, bootstrapMap]),
);
await removeStale(
  jsDir,
  /^jquery-.*\.min\.(?:js|map)$/,
  new Set([jqueryJS, jqueryMap]),
);

const bootstrapSource = path.join(root, "node_modules", "bootstrap", "dist", "css");
const jquerySource = path.join(root, "node_modules", "jquery", "dist");

let css = await readFile(path.join(bootstrapSource, "bootstrap.min.css"), "utf8");
css = css.replace("sourceMappingURL=bootstrap.min.css.map", `sourceMappingURL=${bootstrapMap}`);
await writeFile(path.join(cssDir, bootstrapCSS), css);
await copyFile(path.join(bootstrapSource, "bootstrap.min.css.map"), path.join(cssDir, bootstrapMap));

let js = await readFile(path.join(jquerySource, "jquery.min.js"), "utf8");
js = js.replace("sourceMappingURL=jquery.min.map", `sourceMappingURL=${jqueryMap}`);
await writeFile(path.join(jsDir, jqueryJS), js);
await copyFile(path.join(jquerySource, "jquery.min.map"), path.join(jsDir, jqueryMap));

const indexFile = path.join(publicDir, "index.html");
let html = await readFile(indexFile, "utf8");
html = html.replace(/css\/bootstrap-[^/"']+\.min\.css/g, `css/${bootstrapCSS}`);
html = html.replace(/js\/vendor\/jquery-[^/"']+\.min\.js/g, `js/vendor/${jqueryJS}`);
await writeFile(indexFile, html);

console.log(`Deployed Bootstrap ${bootstrap.version} and jQuery ${jquery.version}`);
