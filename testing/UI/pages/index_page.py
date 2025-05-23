from selenium.webdriver.common.by import By
from utils.utils import BasePage

class IndexPage(BasePage):
    def __init__(self, driver):
        super().__init__(driver)
        self.logo_img = (By.CSS_SELECTOR, ".content nav .logo img")
        self.logo_text = (By.CSS_SELECTOR, ".content nav .logo span")
        self.header_text = (By.CSS_SELECTOR, "header h1")
        self.header_subtext = (By.CSS_SELECTOR, "header p")
        self.footer_link = (By.CSS_SELECTOR, ".links a") # By this we get a list.
        self.section_cards = (By.CSS_SELECTOR, ".section1 .card")
        self.button = (By.ID, "get-started")

    def get_title(self):
        return self.driver.title

    def get_logo_img_src(self):
        # return self.driver.find_element(*self.logo_img).get_attribute("src")
        element = self.wait_for_element(self.logo_img)
        return element.get_attribute("src")

    def get_logo_text(self):
        element = self.wait_for_element(self.logo_text)
        if element :
            return element.text
        else :
            return None

    def get_header_text(self):
        element = self.wait_for_element(self.header_text)
        return element.text if element else None

    def get_header_subtext(self):
        element = self.wait_for_element(self.header_subtext)
        return element.text if element else None

    def get_footer_links(self):
        links = self.driver.find_elements(*self.footer_link)
        print(type(links)) # Here it returns the list.
        return links

    def get_cards(self):
        cards = self.wait_for_elements(self.section_cards)
        return cards


    def click_on_button(self):
        self.driver.find_element(*self.button).click() # Clicking on the button.