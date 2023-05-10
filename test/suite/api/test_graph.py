import json
import time
from typing import List


def verify_graph(mat: List[List[int]]):
    for i in range(len(mat)):
        for j in range(len(mat)):
            if i == j:
                assert mat[i][i] == 0
            else:
                assert mat[i][j] == mat[j][i]


def parse_graph(data: str) -> List[List[int]]:
    return [[int(x) for x in v.strip().split(' ')] for v in data.strip().split("\n")]


def verify_regular(mat: List[List[int]], exp_deg: int):
    counts = [len([v for v in x if v != 0]) for x in mat]

    for k in counts:
        assert k == exp_deg


def verify_between(mat: List[List[int]], min_deg: int, max_deg: int):
    counts = [len([v for v in x if v != 0]) for x in mat]
    for k in counts:
        assert k >= min_deg
        assert k <= max_deg


def verify_complete(mat: List[List[int]]):
    counts = [len([v for v in x if v != 0]) for x in mat]
    nodes = len(mat)

    for k in counts:
        assert k == (nodes - 1)


def verify_connectivity(mat: List[List[int]]):
    found = set()
    queue = [0]

    while len(queue) != 0:
        top = queue[0]
        queue = queue[1:]
        found.add(top)
        for i in range(len(mat)):
            if mat[top][i] != 0 and i not in found:
                queue.append(i)
    assert len(found) == len(mat)


def test_simple_graph_creation(graph_client):
    graph_client.get_builder().type('exact-degree').nodes(8).node_degree(3).connected(True).weighted(False) # {'type': 'exact-degree', 'weighted': False, 'node_degree': 3, 'nodes': 8, 'connected': True}
    response = graph_client.create_graph()

    assert response.status_code == 201
    assert response.type == 'exact-degree'
    assert not response.weighted
    assert response.connected

    graph = graph_client.get_graph(response.id)

    assert response.id == graph.id
    while graph.status != 'finished':
        time.sleep(1)
        graph.refresh()
    if graph.status == 'finished':
        assert graph.download()
        graph_txt = parse_graph(graph.text)
        verify_graph(graph_txt)
        verify_regular(graph_txt, 3)
        verify_connectivity(graph_txt)


def test_same_seed(graph_client):
    graph_client.get_builder().\
        type('average-degree').\
        nodes(25).\
        node_degree_average(6.8).\
        seed(123442233).\
        connected(False)
    resp = graph_client.create_graph()
    assert resp.status_code == 201
    graph_client.get_builder(). \
        type('average-degree'). \
        nodes(25). \
        node_degree_average(6.8). \
        seed(123442233). \
        connected(False)
    resp2 = graph_client.create_graph()
    assert resp2.status_code == 201

    assert resp.id != resp2.id

    while resp.status != 'finished' or resp2.status != 'finished':
        resp.refresh()
        resp2.refresh()
    assert resp.download()
    assert resp2.download()

    graph1 = parse_graph(resp.text)
    verify_graph(graph1)
    graph2 = parse_graph(resp2.text)

    assert graph2 == graph1
