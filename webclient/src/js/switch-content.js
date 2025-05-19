const signup = document.querySelector("#switch-to-signup")
const login = document.querySelector("#switch-to-login")
const signupContent = document.querySelector(".signup-content")
const loginContent = document.querySelector(".login-content")

signup.addEventListener("click", (event)=>{
    event.preventDefault()
    // Now what we have to switch the content.
    loginContent.classList.remove("active")
    signupContent.classList.add("active")
    localStorage.setItem("authView", "signup")
})

login.addEventListener("click", (event)=>{
    event.preventDefault()
    // Now we have to switch the content.
    signupContent.classList.remove("active")
    loginContent.classList.add("active")
    localStorage.setItem("authView", "login")
})

window.addEventListener("DOMContentLoaded", ()=>{
    const view = localStorage.getItem("authView")
    if (view == "signup") {
        loginContent.classList.remove("active")
        signupContent.classList.add("active")
    } else {
        signupContent.classList.remove("active")
        loginContent.classList.add("active")
    }
})