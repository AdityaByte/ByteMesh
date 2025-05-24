document.addEventListener('DOMContentLoaded', function () {
    // Here we have to handle authentication sending the request and fetching response all these things.

    const backendApiUrl = "http://localhost:8080"
    const loginForm = document.querySelector("#login-form")
    const signupForm = document.querySelector("#signup-form")
    const backBtn = document.querySelector("#back-btn")

    backBtn.addEventListener("click", function(event) {
        event.preventDefault()
        window.history.back() // For going back to the previous page
    })

    loginForm.addEventListener("submit", async function (event) {
        event.preventDefault()
        // let formData = new FormData(this) // Gives the form data of the current context.
        // Instead of sending the formdata directly we are sending the data via json.
        const username = document.querySelector("#login-username").value.trim()
        const password = document.querySelector("#login-password").value.trim()
        const info = document.querySelector("#login-info")

        // Here we have to the api request to the backend.
        if (!username || !password) {
            info.textContent = "Fields can't be empty"
            info.style.color = "red"
            info.style.visibility = "visible"
            return
        } else {
            info.textContent = "Checking details..."
            info.style.color = "green"
            info.style.visibility = "visible"
        }

        const payload = {
            username: username,
            password: password
        }



        await fetch(`${backendApiUrl}/login`, {
            method: "POST",
            body: JSON.stringify(payload),
            headers: {
                "content-type": "application/json",
            }
        })
            .then(async response => {
                const responseText = await response.text() // It's a promise so we have to wait till it gets resolved.

                if (response.status != 200) {
                    throw new Error(responseText)
                    return
                }

                let token = response.headers.get("Authorization")
                if (!token || token.trim() === "") {
                    throw new Error("ERROR: No token found in the headers")
                    return
                }
                localStorage.removeItem("authView")
                token = token.replace("Bearer ", "")
                localStorage.setItem("token", token)
                localStorage.setItem("user", username)
                window.location.href = "/dashboard"
            })
            .catch((error) => {
                info.textContent = error.message
                info.style.color = "red"
                info.style.visibility = "visible"
                console.error(error.message)
            })
    })


    signupForm.addEventListener("submit", function (event) {
        event.preventDefault()

        console.log("Signing up...")

        const username = document.querySelector("#signup-username").value.trim()
        const password = document.querySelector("#signup-password").value.trim()
        const info = document.querySelector("#signup-info")

        if (!username || !password) {
            info.textContent = "Fields can't be empty"
            info.style.color = "red"
            info.style.visibility = "visible"
            return
        }
        else if (password.length < 6) {
            info.textContent = "Password should be atleast of 6 characters"
            info.style.color = "red"
            info.style.visibility = "visible"
            return
        }
        else {
            info.textContent = "Checking details..."
            info.style.color = "green"
            info.style.visibility = "visible"
        }

        const payload = {
            username: username,
            password: password
        }

        fetch(`${backendApiUrl}/signup`, {
            method: "POST",
            headers: {
                "content-type": "application/json"
            },
            body: JSON.stringify(payload)
        })
        .then(async response => {
            const text = await response.text()
            if (response.status != 201) {
                throw new Error(text)
                return
            }
            info.textContent = text + ", Now you can login"
            info.style.color = "green"
            info.style.visibility = "visible"
        })
        .catch(error => {
            info.textContent = error.message
            info.style.color = "red"
            info.style.visibility = "visible"
            console.error(error.message)
        })
    })

})