from pages.auth_page import AuthPage
from time import sleep
from os import getenv

def test_signup(driver):
    # Loads the URL
    # driver.get("http://localhost:5501/webclient/src/auth.html")
    driver.get(f"{getenv('WEBCLIENT_URL')}/auth")

    # Creating the object of the AuthPage
    auth_page = AuthPage(driver)

    # Tests

    # Setting the default page to
    driver.execute_script("localStorage.setItem('authView', 'signup')")

    # Now we need to referesh the page
    driver.refresh()
    sleep(2) # Sleeps for a while so the page gets reload

    # Assuming that the user already exists who is aditya.
    auth_page.type_signup_username("aditya")
    auth_page.type_signup_password("demopassword")

    auth_page.click_signup_btn()

    # Now it gives an error saying username already exists so we have to check the error message
    # The UI is correctly displaying the UI message or not.

    # After a while this will shows up an error message
    sleep(2) # Sleeps for a while
    signup_info = auth_page.get_signup_info_element()
    assert signup_info.text == "ERROR: Username already exists! Try another one"
    assert signup_info.value_of_css_property("color") == "rgba(255, 0, 0, 1)"


    # Now trying with a unique one.
    auth_page.type_signup_username("nobody")
    auth_page.type_signup_password("nobody123")

    # Cliks the signup button
    auth_page.click_signup_btn()

    # So if everything goes correctly the UI will shows up a message.
    sleep(3)

    # Pausing the main thread to execute for a while
    signup_info = auth_page.get_signup_info_element()
    assert signup_info.text == "Signup successful, Now you can login"
    assert signup_info.value_of_css_property("color") == "rgba(0, 128, 0, 1)"

    driver.execute_script("localStorage.removeItem('authView')")