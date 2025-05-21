// When the document is loaded initailly i need to set the content ok.
document.addEventListener("DOMContentLoaded", function () {
    // Here we need to fetch out the content from the backend server.
    // So here what i have to do i have to do some thing like search for the token thing n all.

    const homeSwitchingBtn = document.getElementById("home-switcher")
    const uploadSwitchingBtn = document.getElementById("upload-switcher")
    const logout = document.querySelector("#logout-btn")

    logout.addEventListener("click", function (event) {
        event.preventDefault()
        // Here we need to clear the localstorage
        localStorage.removeItem("user")
        localStorage.removeItem("token")
        window.location.href = "/webclient/src/index.html"
    })

    content = localStorage.getItem("content")
    if (content === "main-content") {
        document.querySelector(".main-content").style.display = "flex"
        document.querySelector(".upload-content").style.display = "none"
    } else {
        document.querySelector(".main-content").style.display = "none"
        document.querySelector(".upload-content").style.display = "flex"
    }


    homeSwitchingBtn.addEventListener("click", function (event) {
        event.preventDefault()
        // Now we need to make the other content hidden and show up the main one and have to do other task too like making the backend api request and show up the content.
        document.querySelector(".main-content").style.display = "flex"
        document.querySelector(".upload-content").style.display = "none"
        homeSwitchingBtn.classList.add("active")
        uploadSwitchingBtn.classList.remove("active")
        localStorage.setItem("content", "main-content")
    })

    uploadSwitchingBtn.addEventListener("click", function (event) {
        event.preventDefault()
        document.querySelector(".main-content").style.display = "none"
        document.querySelector(".upload-content").style.display = "flex"
        uploadSwitchingBtn.classList.add("active")
        homeSwitchingBtn.classList.remove("active")
        localStorage.setItem("content", "upload-content")
    })


    token = localStorage.getItem("token")
    if (!token || token.trim() === "") {
        // In that case alert that u cannot access the pages ok.
        alert("Unauthorized")
        window.location.href = "auth.html"
        return
    }

    // console.log("token:", token)
    console.log("Authorized")
    // Else you can load the content.
    user = localStorage.getItem("user")

    console.log(user)
    if (!user || user.trim() === "") {
        console.error("Username is empty")
    }

    document.getElementById("greet").textContent = `Hey ${localStorage.getItem("user")}`

    const gatewayURL = "http://localhost:4444"

    fetch(`${gatewayURL}/fetchall?user=${user}`, {
        method: "GET",
    })
        .then(response => {
            if (response.status != 200) {
                throw new Error("Failed to fetch files")
            }
            return response.json()
        })
        .then(jsonResponse => {
            // Here we need to handle the 201 response.
            // <div class="file-card">
            // <img width="12px" alt="">
            // <span id="filename"></span>
            // <span id="file-upload-date">12/02/2025</span>
            // <button>Download</button>
            // </div>

            // Create a card something like that ok.
            const fileList = document.querySelector(".file-list")
            const imgUrl = "https://imgs.search.brave.com/Ses02BbgloQAO7UoRZnwh1kPMpy0f_WYYJz8polU10Y/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly9zdGF0/aWMtMDAuaWNvbmR1/Y2suY29tL2Fzc2V0/cy4wMC9pbWFnZS1m/aWxlLWljb24tNzgw/eDEwMjQtZnVjOXY1/cmwucG5n"


            // Now i have to set the content

            console.log(jsonResponse)
            console.log(typeof jsonResponse)

            jsonResponse.forEach(data => {
                // here we have to create a structure ok.
                const fileCard = document.createElement("div")
                fileCard.className = "file-card"
                const img = document.createElement("img")
                img.src = `${imgUrl}`
                img.width = "12px"
                const filename = document.createElement("span")
                filename.className = "filename"
                filename.textContent = data.Filename
                const uploadDate = document.createElement("span")
                uploadDate.className = "file-upload-date"
                uploadDate.textContent = data.UploadDate
                const size = document.createElement("span")
                size.className = "filesize"
                size.textContent = data.Size
                const downloadBtn = document.createElement("button")
                downloadBtn.className = "download-btn"
                downloadBtn.innerHTML = "Download"
                fileCard.append(img, filename, size, uploadDate, downloadBtn)
                fileList.append(fileCard)
                console.log("data appended")
            });
        })
        .catch(error => {
            console.error(error)
        })

    document.querySelector(".file-list").addEventListener("click", function (e) {
        if (e.target && e.target.classList.contains("download-btn")) {
            const fileCard = e.target.closest(".file-card");
            const filenameSpan = fileCard.querySelector(".filename");
            const filename = filenameSpan.textContent;
            downloadFile(filename)
        }
    })

    const downloadFile = async (filename) => {
        const response = await fetch(`${gatewayURL}/download?user=${user}&filename=${filename}`, {
            method: "GET"
        })
        .then( async response => {
            if (!response.ok) {
                const errorText = await response.text()
                throw new Error(errorText)
            }
            return response
        })
        .then(async response => {
             // Since if the response is ok we need to extract the binary which was sent by the backend server.
             const blob = await response.blob()
             const url = window.URL.createObjectURL(blob)
             const a = document.createElement("a")
             a.href = url
             a.download = filename
             document.body.appendChild(a)
             a.click()
             a.remove()
             window.URL.revokeObjectURL(url)
        })
        .catch(err => {
            console.error("Download error:", err.message)
            alert("Download failed:"+err.message)
        })
    }


    // Upload functionality
    const uploadForm = document.querySelector("#upload-form")
    const uploadInfo = document.querySelector("#upload-info")

    uploadForm.addEventListener("submit", async function(event) {
        event.preventDefault()
        // Here what we need to do we need to fetch the file.
        const file = document.querySelector("input[name=file]").files[0]

        if (!file) {
            uploadInfo.textContent = "Please upload a file first"
            uploadInfo.style.color = "red"
            uploadInfo.style.visibility = "visible"
            return
        }

        // Else if the user selects the file we need to create a form data
        const formData = new FormData()
        formData.append("file", file)

        // Here we have to send a request to the server.
        await fetch(`${gatewayURL}/upload?user=${user}`, {
            method: "POST",
            body: formData,
        })
        .then(async response => {
            return response.text().then(text => {
                if (response.status != 201) {
                    throw new Error(text)
                }
                return text
            })
        })
        .then(async response => {
            uploadInfo.textContent = response.text()
            uploadInfo.style.color = "white"
            uploadInfo.style.visibility = "visible"
        })
        .catch(err => {
            console.error(err)
            uploadInfo.textContent = "Upload failed: " + err.message;
            uploadInfo.style.color = "red";
            uploadInfo.style.visibility = "visible";
        })
    })
})


