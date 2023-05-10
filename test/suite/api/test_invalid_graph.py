
def test_invalid_regular_graph_creation(graph_client):
    graph_client.get_builder().nodes(12).node_degree(12).connected(True)

    response = graph_client.create_graph()

    assert response.status_code == 400
    assert response.errored()
    assert response.error

    graph_client.reset_builder()
    response = graph_client.create_graph()

    assert response.status_code == 400
    assert response.errored
    assert response.error

    graph_client.get_builder().nodes(11).node_degree(3).connected(True)
    response = graph_client.create_graph()
    assert response.status_code == 400
    assert response.errored
    assert response.error


def test_graph_invalid_limits(graph_client, limits):
    graph_client.get_builder().nodes(limits.max_nodes+5).type("exact-degree").node_degree(limits.max_nodes).connected(True)
    res = graph_client.create_graph()
    assert res.errored()
    graph_client.get_builder().nodes(limits.max_nodes+1).type("average-degree").node_degree_average(limits.max_nodes*0.5).connected(True)
    res = graph_client.create_graph()
    assert res.errored()
    graph_client.get_builder().nodes(limits.max_nodes+3).type("between-degree").node_degree(limits.max_nodes-3).node_degree_max(limits.max_nodes-2).connected(True)
    res = graph_client.create_graph()
    assert res.errored()
    graph_client.get_builder().nodes(limits.max_nodes+5).type("exact-degree").node_degree(limits.max_nodes).connected(True)

    res = graph_client.create_graph()
    assert res.errored()


def test_invalid_json(graph_client):
    invalid_jsons = [{"nodes": 4, "bleh": 12},
                     {"nodes": 6, "type": "exact-degree-test", "node_degree": 3},
                     {"nodes": 6, "type": "exact-degree-test", "node_degree": 3},
                     {"nodes": 6, "type": "exact-degree", "node_degree": 3.5},
                     {"nodes": 6, "type": "exact-degree", "node_degree": 3.5, "unknown-field": "field"}]
    for invalid_json in invalid_jsons:
        res = graph_client.post_any(invalid_json)
        assert res.status_code == 400
