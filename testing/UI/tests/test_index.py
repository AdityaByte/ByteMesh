from pages.index_page import IndexPage

# Writing tests for the index page ok.
def test_index_success(driver):
    # Getting the webpage
    driver.get("http://localhost:5501/webclient/src/index.html")
    # Here we have to create an object of the page model that we have created.
    index_page = IndexPage(driver)
    # Now what we have to do we have to just check it out some tests.
    assert index_page.get_title() == "ByteMesh webclient"
    assert index_page.get_logo_img_src() == "http://localhost:5501/webclient/assets/bytemesh-logo.png"

    assert index_page.get_logo_text() == "ByteMesh."

    assert index_page.get_header_text() == "Distributed File Storage System"

    assert index_page.get_header_subtext() == "Secure, Scalable, and Decentralized Storage for the Modern Web"

    # For the links what we need to do.
    links = index_page.get_footer_links() # Here we gets a list.
    assert len(links) == 2
    # Next we need to check that the links are correct or not.
    github_link = links[0]
    linkedin_link = links[1]

    assert github_link.get_attribute("href") == "https://github.com/AdityaByte/ByteMesh.git"

    assert linkedin_link.get_attribute("href") == "https://www.linkedin.com/in/aditya-pawar-557a56332/"

    assert len(index_page.get_cards()) == 3

    index_page.click_on_button() # Running this at last
