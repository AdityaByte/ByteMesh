from utils.utils import BasePage
from selenium.webdriver.common.by import By

class AuthPage(BasePage):
    def __init__(self, driver):
        super().__init__(driver)
        self.back_btn = (By.ID, "back-btn")
        self.login_input_username = (By.ID, "login-username")
        self.login_input_password = (By.ID, "login-password")
        self.login_btn = (By.CSS_SELECTOR, "#login-form button[type=submit]")
        self.signup_input_username = (By.ID, "signup-username")
        self.signup_input_password = (By.ID, "signup-password")
        self.signup_btn = (By.CSS_SELECTOR, "#signup-form button[type=submit]")
        self.switch_to_signup_link = (By.ID, "switch-to-signup")
        self.switch_to_login_link = (By.ID, "switch-to-login")
        self.login_info_element = (By.ID, "login-info")
        self.signup_info_element = (By.ID, "signup-info")

    def get_title(self):
        return self.driver.title

    def click_back_btn(self):
        self.wait_for_element(self.back_btn).click()

    def type_login_username(self, username):
        inp = self.wait_for_element(self.login_input_username)
        inp.clear()
        inp.send_keys(username)

    def type_login_password(self, password):
        inp = self.wait_for_element(self.login_input_password)
        inp.clear()
        inp.send_keys(password)

    def type_signup_username(self, username):
        inp = self.wait_for_element(self.signup_input_username)
        inp.clear()
        inp.send_keys(username)

    def type_signup_password(self, password):
        inp = self.wait_for_element(self.signup_input_password)
        inp.clear()
        inp.send_keys(password)

    def click_switch_to_signup_link(self):
        self.wait_for_element(self.switch_to_signup_link).click()

    def click_switch_to_login_link(self):
        self.wait_for_element(self.switch_to_login_link).click()

    def get_login_info_element(self):
        return self.wait_for_element(self.login_info_element)

    def get_signup_info_element(self):
        return self.wait_for_element(self.signup_info_element)

    def click_login_btn(self):
        self.wait_for_element(self.login_btn).click()

    def click_signup_btn(self):
        self.wait_for_element(self.signup_btn).click()