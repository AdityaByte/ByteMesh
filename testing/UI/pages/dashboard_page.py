from utils.utils import BasePage
from selenium.webdriver.common.by import By

class DashboardPage(BasePage):
    def __init__(self, driver):
        super().__init__(driver)
        self.logout_btn = (By.ID, "logout-btn")
        self.greet = (By.ID, "greet")
        self.home_switcher_btn = (By.CSS_SELECTOR, ".container .base-bar #home-switcher")
        self.upload_switcher_btn = (By.CSS_SELECTOR, ".container .base-bar #upload-switcher")
        self.main_content = (By.CSS_SELECTOR, ".container .main-content")
        self.upload_content = (By.CSS_SELECTOR, ".container .upload-content")

    def get_greet_text(self):
        greet = self.wait_for_element(self.greet)
        return greet.text if greet else None

    def click_logout_btn(self):
        self.wait_for_element(self.logout_btn).click()

    def click_home_switcher_btn(self):
        element = self.wait_for_element(self.home_switcher_btn)
        element.click()

    def click_upload_switcher_btn(self):
        element = self.wait_for_element(self.upload_switcher_btn)
        element.click()

    def get_main_content(self):
        return self.wait_for_element(self.main_content)

    def get_upload_content(self):
        return self.wait_for_element(self.upload_content)