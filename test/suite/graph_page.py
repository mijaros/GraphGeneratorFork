from seleniumpagefactory import PageFactory


class GraphPage(PageFactory):
    def __init__(self, driver):
        self.driver = driver
        self.timeout = 10

