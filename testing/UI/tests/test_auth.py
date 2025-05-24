from pages.auth_page import AuthPage
from os import getenv

# This test is only for UI test it doesn't tests login and signup functionality.
def test_auth(driver):
    # Loads the authentication page.
    # driver.get("http://localhost:5501/webclient/src/auth.html") # Old way
    driver.get(f"{getenv('WEBCLIENT_URL')}/auth")

    # Creating object of AuthPage
    auth_page = AuthPage(driver)

    # Tests
    assert auth_page.get_title() == "Authentication Page"

    driver.execute_script("localStorage.setItem('authView', 'login')") # Setting the default view to login.

    auth_page.type_login_username("   ")
    auth_page.type_login_password("   ")

    auth_page.click_login_btn()

    login_info = auth_page.get_login_info_element()

    assert login_info.text == "Fields can't be empty"
    assert login_info.value_of_css_property("color") == "rgba(255, 0, 0, 1)" # red gets converted to this.

    auth_page.click_switch_to_signup_link()

    auth_page.type_signup_username("    ")
    auth_page.type_signup_password("    ")

    auth_page.click_signup_btn()

    signup_info = auth_page.get_signup_info_element()

    assert signup_info.text == "Fields can't be empty"
    assert signup_info.value_of_css_property("color") == "rgba(255, 0, 0, 1)" # Red gets converted to this.

    auth_page.type_signup_username("hello")
    auth_page.type_signup_password("world")

    auth_page.click_signup_btn()

    signup_info = auth_page.get_signup_info_element()

    assert signup_info.text == "Password should be atleast of 6 characters"
    assert signup_info.value_of_css_property("color") == "rgba(255, 0, 0, 1)"

    auth_page.click_back_btn()