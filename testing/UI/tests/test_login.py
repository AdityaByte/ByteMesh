from pages.auth_page import AuthPage
from time import sleep
from os import getenv

def test_login(driver):
    # Load the page
    # driver.get("http://localhost:5501/webclient/src/auth.html")
    driver.get(f"{getenv('WEBCLIENT_URL')}/auth")

    # Creating the object of the page.
    auth_page = AuthPage(driver)

    # clearing the localstorage.
    driver.execute_script("localStorage.clear()")

    driver.refresh()

    sleep(2)

    # If the credentials are not valid.
    auth_page.type_login_username("xyz") # username not exists
    auth_page.type_login_password("xyz123")

    auth_page.click_login_btn()

    sleep(4) # Waits for 2 second for the message to display.

    login_info = auth_page.get_login_info_element()

    assert login_info.text == "ERROR: Invalid Credentials"
    assert login_info.value_of_css_property("color") == "rgba(255, 0, 0, 1)"

    # When the credentials are valid
    username = "aditya"
    auth_page.type_login_username(username)
    auth_page.type_login_password("aditya123")

    auth_page.click_login_btn()

    # Now what we have to check we have to check that we are getting the token value or not.
    sleep(4)

    storedUsername = driver.execute_script("return localStorage.getItem('user')")
    tokenValue = driver.execute_script("return localStorage.getItem('token')")

    assert storedUsername == username
    assert tokenValue.strip() != ""