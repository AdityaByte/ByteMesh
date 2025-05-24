const http = require("http")
const fs = require("fs")
const path = require("path")

const port = 5501

const mimeTypes = {
    ".html": "text/html",
    ".js": "text/javascript",
    ".css": "text/css",
    ".png": "image/png",
    ".jpg": "image/jpeg",
    ".svg": "image/svg+xml",
    ".ico": "image/x-icon"
}

const server = http.createServer((req, res) => {
    let filePath = ""

    if (req.url === "/") {
        filePath = path.join(__dirname, "src", "index.html")
    } else if (req.url === "/auth") {
        filePath = path.join(__dirname, "src", "auth.html")
    } else if (req.url === "/about") {
        filePath = path.join(__dirname, "src", "about.html")
    } else if (req.url === "/dashboard") {
        filePath = path.join(__dirname, "src", "dashboard.html")
    } else {
        // Serving static files.
        const urlPath = req.url.startsWith("/assets/")
            ? path.join(__dirname, req.url)
            : path.join(__dirname, "src", req.url)
        filePath = urlPath
    }

    const ext = path.extname(filePath)
    const contentType = mimeTypes[ext] || "text/plain"

    fs.readFile(filePath, (err, content) => {
        if (err) {
            res.writeHead(404)
            res.end("404 not found")
        } else {
            res.writeHead(200, { "Content-Type": contentType });
            res.end(content);
        }
    })
})

server.listen(port, () => {
    console.log(`Server is running at  http://localhost:${port}`)
})