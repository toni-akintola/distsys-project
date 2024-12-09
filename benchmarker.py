import time
import requests
from concurrent.futures import ThreadPoolExecutor
import itertools
import json


def testGetStockLatency():
    url: str = "http://localhost:9445/single-stock/AAPL"
    start_time = time.time()
    response = requests.get(url)

    if response.status_code == 200:
        end_time = time.time()
        return end_time - start_time
    else:
        print(response.json())
        return None


def testBuyOrderLatency():
    url: str = "http://localhost:9445/order/"

    request_body = {"ticker": "AAPL", "quantity": 1, "username": "oakintol"}

    start_time = time.time()
    response = requests.post(url, json=request_body)

    if response.status_code == 200:
        end_time = time.time()
        return 1
    else:
        print(response.json())
        return None


def testSellOrderLatency():
    url: str = "http://localhost:9445/order/"

    request_body = {"ticker": "AAPL", "quantity": -1, "username": "oakintol"}

    start_time = time.time()
    response = requests.post(url, json=request_body)

    if response.status_code == 200:
        end_time = time.time()
        return 1
    else:
        print(response.json())
        return None


def testGetStockThroughput():
    url: str = "http://localhost:9445/single-stock/AAPL"
    response = requests.get(url)

    if response.status_code == 200:
        return 1
    else:
        print(response.json())
        return None


def testBuyOrderThroughput():
    url: str = "http://localhost:9445/order/"

    request_body = {"ticker": "AAPL", "quantity": 1, "username": "oakintol"}

    response = requests.post(url, json=request_body)

    if response.status_code == 200:
        return 1
    else:
        print(response.json())
        return None


def testSellOrderThroughput():
    url: str = "http://localhost:9445/order/"

    request_body = {"ticker": "AAPL", "quantity": -1, "username": "oakintol"}

    response = requests.post(url, json=request_body)

    if response.status_code == 200:
        return 1
    else:
        print(response.json())
        return None


def testLatency(numClients: int, testFunc):
    with ThreadPoolExecutor(max_workers=numClients) as executor:
        results = executor.map(lambda _: testFunc(), itertools.repeat(None, numClients))
    return [res for res in results if res is not None]


def testThroughput(numClients: int, testFunc):
    start_time = time.time()
    res = 0
    while time.time() - start_time < 5:
        with ThreadPoolExecutor(max_workers=numClients) as executor:
            results = executor.map(
                lambda _: testFunc(), itertools.repeat(None, numClients)
            )
        res += sum([res for res in results if res is not None])
    return res


if __name__ == "__main__":
    times = testThroughput(30, testGetStockThroughput)
    print(times)
