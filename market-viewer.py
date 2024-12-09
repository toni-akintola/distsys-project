import time
import os

# user can run this in a terminal to see a continuously-updated market
# pulling from the log so only displaying ticker and price but this can be changed

log = "market_log.txt"

seen = {}


def read_log():
    with open(log, "r") as fp:
        lines = fp.readlines()
    return lines


def process(lines):
    # stocks [ticker] = price
    stocks = {}
    for line in lines:
        if line.startswith("LOG:"):
            parts = line.strip().split(" ")
            ticker = parts[4].strip(",")
            price = parts[7]
            stocks[ticker] = float(price)
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
    lines = read_log()
    stocks = process(lines)

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
