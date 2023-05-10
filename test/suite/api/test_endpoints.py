import http.client

import requests


def test_invalid_paths(server_root):
    test_random = requests.get(f"{server_root}/api/asdfve")
    assert test_random.status_code == 404
    test_random = requests.get(f"{server_root}/api/")
    assert test_random.status_code == 404
    test_random = requests.get(f"{server_root}/api/v1")
    assert test_random.status_code == 404

    test_pass = requests.get(f"{server_root}/api/v1/graph")
    assert test_pass.status_code == 200
    test_pass = requests.get(f"{server_root}/api/v1/batch")
    assert test_pass.status_code == 200
    test_pass = requests.get(f"{server_root}/api/v1/limits")
    assert test_pass.status_code == 200


def test_valid_paths(server_root):
    test_get = requests.get(f"{server_root}/api/v1/graph")
    assert test_get.status_code == 200

    test_patch = requests.patch(f"{server_root}/api/v1/graph")
    assert test_patch.status_code == 404

    test_put = requests.put(f"{server_root}/api/v1/graph")
    assert test_put.status_code == 404

    test_batch_get = requests.get(f"{server_root}/api/v1/batch")
    assert test_get.status_code == 200
    test_patch = requests.patch(f"{server_root}/api/v1/batch")
    assert test_patch.status_code == 404

    test_put = requests.put(f"{server_root}/api/v1/batch")
    assert test_put.status_code == 404


def test_invalid_graph_methods(graph_client, server_root):
    graph_client.get_builder().nodes(12).type("complete")
    graph = graph_client.create_graph()

    assert graph.status_code == 201

    test_get = requests.get(f"{server_root}/api/v1/graph/{graph.id}")
    assert test_get.status_code == 200

    test_post = requests.post(f"{server_root}/api/v1/graph/{graph.id}", json={"bleh": "boo"})
    assert test_post.status_code == 404
