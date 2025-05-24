import pytest
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from dotenv import load_dotenv # Library for loading the .env file.

# Config method for loading the driver.
@pytest.fixture(scope="session") # Session cause we are only kept one browser for each test request.
def driver():
    # paths
    brave_path = r"C:\Program Files\BraveSoftware\Brave-Browser\Application\brave.exe"
    driver_path = r"F:\program files\selenium webdriver\chromedriver-brave\chromedriver.exe"

    options = Options()
    options.binary_location = brave_path # Setted the binary location cause the brave is too chromium based.
    options.add_argument("--start-maximized")

    service = Service(executable_path=driver_path)
    driver = webdriver.Chrome(options=options, service=service)

    yield driver # Returning the driver to all tests.

    driver.quit() # Shuts down after all tests.

# Pytest will automatically call it.
def pytest_configure():
    load_dotenv(dotenv_path=".env.tests")