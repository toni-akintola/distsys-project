import time
import requests
from concurrent.futures import ThreadPoolExecutor
import itertools


def testGetStockLatency():
    url: str = "http://localhost:9444/single-stock/AAPL"
    start_time = time.time()
    response = requests.get(url)

    if response.status_code == 200:
        end_time = time.time()
        return end_time - start_time  # Return elapsed time
    else:
        print(response.json())
        return None  # Return None if the request fails


def testBuyOrderLatency():
    url: str = "http://localhost:9444/order/"


def testSellOrderLatency():
    pass


def testClients(numClients: int, testFunc):
    with ThreadPoolExecutor(max_workers=numClients) as executor:
        results = executor.map(
            lambda _: testFunc(), itertools.repeat(None, numClients)
        )  # Collect results
    return [res for res in results if res is not None]  # Filter out None


if __name__ == "__main__":
    times = testClients(30, testGetStockLatency)
    print(times)
