document.addEventListener('DOMContentLoaded', function (event) {
    event.preventDefault()
    // Here we have to handle authentication sending the request and fetching response all these things.

    const backendApiUrl = "http://localhost:8080"

    const loginForm = document.querySelector("#login-form")
    const signupForm = document.querySelector("#signup-form")

    loginForm.addEventListener("submit", function (event) {
        event.preventDefault()
        // let formData = new FormData(this) // Gives the form data of the current context.
        // Instead of sending the formdata directly we are sending the data via json.
        const username = document.querySelector("[name=username").value.trim()
        const password = document.querySelector("[name=password]").value.trim()

        // Here we have to the api request to the backend.
        if (!username || !password) {
            const error = document.querySelector("#info")
            error.textContent = "Fields can't be empty"
            error.style.color = "red"
            error.style.visibility = "visible"
            return
        } else {
            const info = document.querySelector("#info")
            info.textContent = "Checking details..."
            info.style.color = "green"
            info.style.visibility = "visible"
        }

        const payload = {
            username: username,
            password: password
        }


        const makeRequest = (async () => {

            await fetch(`${backendApiUrl}/login`, {
                method: "POST",
                body: JSON.stringify(payload),
                headers: {
                    "content-type": "application/json",
                }
            })
                .then((response) => {
                    if (response.status === 200) {
                        localStorage.removeItem("authView")
                        // So when the user logs in we need to save the token in the localstorage
                        let token = response.headers.get("Authorization")
                        if (token) {
                            let trimmedToken = token.replace("Bearer ", "");
                            localStorage.setItem("token", trimmedToken)
                            // Also need to set the client username too.
                            localStorage.setItem("user", username)
                            window.location.href = "/webclient/src/dashboard.html"
                        } else {
                            throw new Error("Authorization header missing!")
                        }
                    }
                    throw new Error(`ERROR: Response status ${response.status} and Response text ${response.text}`)
                })
                .catch((error) => {
                    console.error(error)
                })
        })() // IIFE
    })


    signupForm.addEventListener("submit", function (event) {
        event.preventDefault()

        console.log("working..")

        const username = document.querySelector("[name=username]").value.trim()
        const password = document.querySelector("[name=password]").value.trim()
        const info = document.querySelector("#info")

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
            .then(response => {
                if (response == 201) {
                    alert("Credentials created now login")
                    document.querySelector(".signup-content").classList.remove("active")
                    document.querySelector(".login-content").classList.add("active")
                    // window.location.href = "/webclient/src/dashboard.html"
                    return
                }
                // Else we need to generate some error and print back the error
                throw new Error(response.text)
            })
            .catch(error => {
                console.error(error)
            })
    })
})

