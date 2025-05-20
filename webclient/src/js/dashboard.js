const homeSwitchingBtn = document.getElementById("home-switcher")
const uploadSwitchingBtn = document.getElementById("upload-switcher")

homeSwitchingBtn.addEventListener("click", function(event) {
    event.preventDefault()
    // Now we need to make the other content hidden and show up the main one and have to do other task too like making the backend api request and show up the content.
    document.querySelector(".main-content").style.display = "flex"
    document.querySelector(".upload-content").style.display = "none"
    homeSwitchingBtn.classList.add("active")
    uploadSwitchingBtn.classList.remove("active")
    localStorage.setItem("content", "main-content")
})

uploadSwitchingBtn.addEventListener("click", function(event) {
    event.preventDefault()
    document.querySelector(".main-content").style.display = "none"
    document.querySelector(".upload-content").style.display = "flex"
    uploadSwitchingBtn.classList.add("active")
    homeSwitchingBtn.classList.remove("active")
    localStorage.setItem("content", "upload-content")
})

document.addEventListener("DOMContentLoaded", function() {
    // Here we need to read out the data from the local storage.
    content = localStorage.getItem("content")
    if (content === "main-content") {
        document.querySelector(".main-content").style.display = "flex"
        document.querySelector(".upload-content").style.display = "none"
    } else {
        document.querySelector(".main-content").style.display = "none"
        document.querySelector(".upload-content").style.display = "flex"
    }
})