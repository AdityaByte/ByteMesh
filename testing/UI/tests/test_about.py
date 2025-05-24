from pages.about_page import AboutPage
from re import search
from os import getenv

# Tests of the about page
def test_about(driver):
    driver.get(f"{getenv('WEBCLIENT_URL')}/about")

    # Creating object of the About Page
    about_page = AboutPage(driver)

    # Tests
    elements = about_page.get_headings()
    assert len(elements) == 2
    assert elements[0].text == "💻 Bytemesh Web Client"
    assert elements[1].text == "💡 System architecture"

    # Now we have to check the system architecture image url.
    bg_img = about_page.get_architecture_img_url()
    url = search(r'url\("(.*?)"\)', bg_img)
    assert url.group(1) == "http://localhost:5501/assets/webclient-bytemesh-architecture.png"

    # At last we need to check the back functionality
    about_page.click_back_btn()
