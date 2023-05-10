import requests
from typing import Dict, Optional
from types import SimpleNamespace


class GraphRequestBuilder:
    def nodes(self, nodes: int) -> 'GraphRequestBuilder':
        self.p_nodes = nodes
        return self

    def node_degree(self, degree: int) -> 'GraphRequestBuilder':
        self.p_node_degree = degree
        return self

    def node_degree_max(self, max_degree: int) -> 'GraphRequestBuilder':
        self.p_node_degree_max = max_degree
        return self

    def node_degree_average(self, node_degree_average: float) -> 'GraphRequestBuilder':
        self.p_node_degree_average = node_degree_average
        return self

    def type(self, type: str) -> 'GraphRequestBuilder':
        self.p_type = type
        return self

    def connected(self, con: bool) -> 'GraphRequestBuilder':
        self.p_connected = con
        return self

    def seed(self, seed: int) -> 'GraphRequestBuilder':
        self.p_seed = seed
        return self


    def weighted(self, weigh: bool)-> 'GraphRequestBuilder':
        self.p_weighted = weigh
        return self

    def reset(self):
        to_clear = [k for k in dir(self) if k.startswith("p_")]
        for v in to_clear:
            self.__delattr__(v)

    def build(self) -> Dict:
        result = {}
        for k in dir(self):
            if k.startswith("p_"):
                result[k[2:]] = self.__getattribute__(k)
        return result


class LimitsClient:
    def __init__(self, root_url):
        self._url = f"{root_url}/api/v1/limits"
        self._initialized = False

    def _resolve(self):
        if self._initialized:
            return
        res = requests.get(self._url)
        if res.status_code != 200:
            raise EnvironmentError()
        body = res.json()
        if not body['max_nodes']:
            raise EnvironmentError()
        if not body['max_batch_size']:
            raise EnvironmentError()
        j = res.json()
        self._max_nodes = j['max_nodes']
        self._max_batch_size = j['max_batch_size']
        self._initialized = True

    @property
    def max_nodes(self):
        self._resolve()
        return self._max_nodes

    @property
    def max_batch_size(self):
        self._resolve()
        return self._max_nodes


class GraphResponse(object):

    def __init__(self, raw_response: requests.Response, client: 'GraphClient'):
        self._raw_response = raw_response
        self._unpacked = None
        self._errored = False
        self._client = client
        self._text  = None

        if raw_response.status_code >= 200 or raw_response.status_code == 201:
            self._unpacked = raw_response.json(object_hook=lambda v: SimpleNamespace(**v))
        if raw_response.status_code >= 400:
            self._errored = True
            self._error = raw_response.json(object_hook=lambda v: SimpleNamespace(**v))

    def errored(self) -> bool:
        return self._errored

    def __getattr__(self, name):
            return getattr(self._unpacked, name)

    @property
    def status_code(self):
        return self._raw_response.status_code

    def refresh(self) -> bool:
        newObj = self._client.get_graph(self._unpacked.id)
        if newObj.errored():
            return False
        self._unpacked = newObj._unpacked

    def download(self) -> bool:
        resp = self._client.download_graph(self._unpacked.id)
        if resp.status_code != 200:
            return False
        self._text = resp.text
        return True

    @property
    def text(self) -> Optional[str]:
        return self._text


class GraphClient:
    def __init__(self, root_url):
        self._url = f"{root_url}/api/v1/graph"
        self._builder = GraphRequestBuilder()

    def reset_builder(self):
        self._builder.reset()

    def get_builder(self):
        return self._builder

    def create_graph(self):
        resp = requests.post(self._url, json=self._builder.build())
        return GraphResponse(resp, self)

    def get_graph(self, id):
        resp = requests.get(self._url + f"/{id}")
        return GraphResponse(resp, self)

    def post_any(self, obj) -> requests.Response:
        resp = requests.post(self._url, json=obj)
        return resp

    def download_graph(self, id):
        return requests.get(self._url + f"/{id}/download")

    def refresh(self, req):
        id = req.json()['id']
        return requests.get(self._url + f"/{id}")


