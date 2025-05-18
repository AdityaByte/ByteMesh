document.addEventListener("DOMContentLoaded", function (event) {
    event.preventDefault()
    document.getElementById("get-started").addEventListener("click", function (event) {
        event.preventDefault()
        window.location.href = "/src/auth.html"
    })
})