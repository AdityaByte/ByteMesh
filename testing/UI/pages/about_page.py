from selenium.webdriver.common.by import By
from utils.utils import BasePage

class AboutPage(BasePage):
    def __init__(self, driver):
        super().__init__(driver)
        self.back_btn = (By.ID, "back-btn")
        self.headings = (By.TAG_NAME, "h1")
        self.architecture_img = (By.CLASS_NAME, "architecture")

    def get_title(self):
        return self.driver.title

    def get_headings(self):
        elements = self.wait_for_elements(self.headings)
        return elements # Here we get the list of two elements.

    def get_architecture_img_url(self):
        element = self.wait_for_element(self.architecture_img)
        # Here we need to fetch the URL of the background.
        return element.value_of_css_property("background-image")

    def click_back_btn(self):
        btn = self.wait_for_element(self.back_btn)
        btn.click()