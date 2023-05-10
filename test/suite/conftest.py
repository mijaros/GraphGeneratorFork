import pytest

from suite.graph_api import GraphClient, LimitsClient
from selenium.webdriver import Firefox, Chrome


def pytest_addoption(parser):
    parser.addoption(
        '--chrome', action='store_true', help='If passed chrome driver is used'
    )
    parser.addoption(
        '--url', action='store', default='http://localhost:8080', help='Set url of SUT'
    )


@pytest.fixture
def server_root(request):
    return request.config.getoption('--url')


@pytest.fixture
def graph_client(server_root):
    return GraphClient(server_root)


@pytest.fixture
def web_driver(request):
    if request.config.getoption('--chrome'):
        _driver = Chrome()
    else:
        _driver = Firefox()
    yield _driver
    _driver.close()


@pytest.fixture
def limits(server_root):
    return LimitsClient(server_root)
