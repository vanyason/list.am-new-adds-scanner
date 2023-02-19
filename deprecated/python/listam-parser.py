# Logic:
#
# 1. Parse setup config
# 2. Get data and Parse data
# 3. If file with data do not exist:
#       Create file. Notify about file creation. Exit
#    Else
#       Compare new arrived data with the data in file. Save diffs locally. Update file
# 5. Show diffs if they exist

from enum import unique
from bs4 import BeautifulSoup
import requests
import os
import logging
import time
import telegram_send
import argparse


def usdToAmd(price):
    return price * 500


def findUniqueElements(oldList, newList):
    return list(set(newList) - set(oldList))

# currency = 1 means USD, 0 - AMD
def getAdsFromListAm(price, roomsAmount, currency):
    if currency != 1:
        price = usdToAmd(price)

    headers = {'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36'}
    listAmUrl = f"https://www.list.am/ru/category/56/1?cmtype=1&pfreq=1&po=1&n=1&price2={price}&crc={currency}&_a4={roomsAmount}"
    listAmUrlOriginal = f"https://www.list.am/category/56?cmtype=1&pfreq=1&po=1&n=1&price2={price}&crc={currency}&_a4={roomsAmount}"
    pages = []
    pageNum = 1

    # get all the html pages possible
    # when pages end - list am will respomd with the original url
    while True:
        try:
            resp = requests.get(listAmUrl, headers=headers)
        except Exception as ex:
            raise ConnectionError(f'Some requests problem: info {ex}')

        if (listAmUrlOriginal == resp.url):
            break
        if (resp.status_code != 200):
            raise ConnectionError('list am didn`t return code 200')

        pages.append(resp.text)
        listAmUrl = listAmUrl.replace(f'/{pageNum}?', f'/{pageNum + 1}?')
        pageNum = pageNum + 1

    # exctract links from the pages
    links = set()
    for page in pages:
        soup = BeautifulSoup(page, 'lxml')
        for link in soup.find_all('a'):
            if "/item/" in str(link):
                links.add('https://www.list.am' + link.get('href'))

    return links


if __name__ == "__main__":
    # Parse args
    try:
        parser = argparse.ArgumentParser()
        parser.add_argument('--price', type=int, required=True)
        parser.add_argument('--rooms', type=int, required=True)
        args = parser.parse_args()
    except Exception as ex:
        print(ex)
        quit()

    # Setup params
    refreshRate = 60
    price = args.price
    roomsAmout = args.rooms
    dbFileName = f'db_price_{price}USD_rooms_{roomsAmout}.db'
    loggingFileName = f'logs_price_{price}USD_rooms_{roomsAmout}.log'
    logging.basicConfig(format='%(asctime)s - %(levelname)s:%(message)s',
                        datefmt='%m/%d/%Y %I:%M:%S %p',
                        filename=loggingFileName,
                        level=logging.INFO)

    logging.info(f'\n\r\n\rStarted: price: {price}, rooms: {roomsAmout}, refresh rate(s): {refreshRate}')
    print(f'Started: price: {price}, rooms: {roomsAmout}, refresh rate(s): {refreshRate}')

    while True:
        try:
            # Get data from List am
            links = set()
            links.update(getAdsFromListAm(price, roomsAmout, 1))
            links.update(getAdsFromListAm(price, roomsAmout, 0))

            # If db doesn`t exist - create
            if not os.path.exists(dbFileName):
                logging.info(f'BD doesn`t exist. Creating: {dbFileName}')
                with open(dbFileName, "w") as db:
                    db.write('\n'.join(links))
                logging.info(f'DB Filled')

            # Get old data from db and compare with the new one
            with open(dbFileName, "r") as db:
                oldlinks = [line.rstrip('\n') for line in db]

            # if new elements appear - save them and send to the telegramm
            newLinks = findUniqueElements(oldlinks, links)
            if newLinks:
                logging.info(f'NEW LINKS!!! : {newLinks}')
                telegram_send.send(messages=newLinks)

            # refresh db with the new links
            with open(dbFileName, "w") as db:
                db.write('\n'.join(links))

        except Exception as ex:
            logging.exception(ex)
            print(ex)

        time.sleep(refreshRate)