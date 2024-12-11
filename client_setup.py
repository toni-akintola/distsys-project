import json
import http.client


def findServer():
    host = "catalog.cse.nd.edu"
    path = "/query.json"
    port = 9097

    conn = http.client.HTTPConnection(host, port)

    conn.request("GET", path)

    response = conn.getresponse()

    data = response.read().decode("utf-8")
    json_data = json.loads(data)
    my_servers = list(
        filter(
            lambda item: type(item) == dict and item.get("type") == "south-bend-bets",
            json_data,
        )
    )
    my_servers.sort(key=lambda item: item["lastheardfrom"], reverse=True)

    return "http://" + my_servers[0]["name"] + ":" + str(my_servers[0]["port"])


if __name__ == "__main__":
    print(f"Here is the executor server URL: {findServer()}")
