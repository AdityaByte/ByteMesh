# Navigating to the dashboard we first need login
from os import getenv
from time import sleep
from pages.auth_page import AuthPage
from pages.dashboard_page import DashboardPage
from selenium.webdriver.support.wait import WebDriverWait

def test_dashboard(driver):
    driver.get(f"{getenv('WEBCLIENT_URL')}/auth")

    driver.execute_script("localStorage.setItem('authView', 'login')")
    driver.refresh()
    sleep(2)

    auth_page = AuthPage(driver)
    # Valid credentials
    auth_page.type_login_username("aditya")
    auth_page.type_login_password("aditya123")
    auth_page.click_login_btn()

    sleep(3)

    # Testing just UI components.
    dashboard_page = DashboardPage(driver)

    # Switching to the upload content
    dashboard_page.click_upload_switcher_btn()

    # When it switches to the upload page the upload-content must be visible and some keys are saved in the localstorage we need to check them.
    upload_content = dashboard_page.get_upload_content()
    # main_content = dashboard_page.get_main_content()

    assert upload_content.value_of_css_property("display") == "flex"
    # assert main_content.value_of_css_property("display") == "none"

    content = driver.execute_script("return localStorage.getItem('content')")

    assert content == "upload-content"


    # Switching to the home content
    dashboard_page.click_home_switcher_btn()

    # Some checks

    # upload_content = dashboard_page.get_upload_content()
    main_content = dashboard_page.get_main_content()

    # assert upload_content.value_of_css_property("display") == "none"
    assert main_content.value_of_css_property("display") == "flex"

    content = driver.execute_script("return localStorage.getItem('content')")

    assert content == "main-content"

    # checking the logout button functionality

    dashboard_page.click_logout_btn()

    # If the user log out now we need to check of the token and the user localstorage are cleared out or not.
    user = driver.execute_script("return localStorage.getItem('user')")
    token = driver.execute_script("return localStorage.getItem('token')")

    assert user == None
    assert token == None