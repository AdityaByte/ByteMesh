from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException

class BasePage:
    def __init__(self, driver, timeout=6):
        self.driver = driver
        self.timeout = timeout

    def wait_for_element(self, by_locator):
        try:
            element = WebDriverWait(self.driver, self.timeout).until(
                EC.visibility_of_element_located(by_locator)
            )
            return element
        except TimeoutException:
            print(f"Element {by_locator} not found after {self.timeout} seconds")
            return None

    def wait_for_elements(self, by_locator):
        try:
            elements = WebDriverWait(self.driver, self.timeout).until(
                EC.presence_of_all_elements_located(by_locator)
            )

            return elements
        except TimeoutException:
            print(f"Elements {by_locator} not found after {self.timeout} seconds")
            return None
