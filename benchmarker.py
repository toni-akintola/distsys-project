import time
import requests
from concurrent.futures import ThreadPoolExecutor
import itertools


def testClient():
    url: str = "http://localhost:9444/single-stock/AAPL"
    start_time = time.time()
    response = requests.get(url)

    if response.status_code == 200:
        end_time = time.time()
        return end_time - start_time  # Return elapsed time
    else:
        print(response.json())
        return None  # Return None if the request fails


def testClientScale(numClients: int):
    with ThreadPoolExecutor(max_workers=numClients) as executor:
        results = executor.map(lambda _: testClient(), itertools.repeat(
            None, numClients))  # Collect results
    return [res for res in results if res is not None]  # Filter out None


if __name__ == "__main__":
    times = testClientScale(30)
    print(times)
