import json
import time
import os

# user can run this in a terminal to see a continuously-updated market
# pulling from the log so only displaying ticker and price but this can be changed

log = "stocks.json"

seen = {}


def read_log():
    with open(log, "r") as fp:
        data = json.load(fp)
    return data


def process(data):
    # stocks [ticker] = price
    stocks = {}
    for stock in data:
        ticker = stock["ticker"]
        price = stock["currentPrice"]
        stocks[ticker] = price
    return stocks


def display(stocks):
    # no scrolling output
    os.system("cls" if os.name == "nt" else "clear")
    print(f"{'TICKER':<10} {'PRICE':>10}")
    for ticker, price in stocks.items():
        print(f"{ticker:<10} {price:>10.2f}")


def detect_changes():
    changed = False

    global seen
    data = read_log()
    stocks = process(data)

    for ticker, price in stocks.items():
        if ticker not in seen or seen[ticker] != price:
            changed = True

    if changed:
        display(stocks)
        seen = stocks


def main():
    while True:
        detect_changes()
        time.sleep(1)


if __name__ == "__main__":
    main()
