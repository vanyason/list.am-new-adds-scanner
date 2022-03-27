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


def findUniqueElements(oldList, newList):
    return list(set(newList) - set(oldList))


def getParams(price, roomsAmount):
    headers = {'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36'}
    listAmUrl = f"https://www.list.am/ru/category/56/1?pfreq=1&po=1&n=0&price2={price}&crc=1&_a4={roomsAmount}"
    listAmUrlOriginal = f"https://www.list.am/category/56?pfreq=1&po=1&n=0&price2={price}&crc=1&_a4={roomsAmount}"
    maxErrorsToStop = 10
    return {'headers': headers, 'listAmUrl': listAmUrl, 'listAmUrlOriginal': listAmUrlOriginal, 'maxErrorsToStop': maxErrorsToStop}


def getAdsFromListAm(params):
    pages = []
    pageNum = 1
    url = params['listAmUrl']

    # get all the html pages possible
    # when pages end - list am will respomd with the original url
    while True:
        resp = requests.get(url, headers=params['headers'])
        if (params['listAmUrlOriginal'] == resp.url):
            break
        if (resp.status_code != 200):
            raise ConnectionError('list am didn`t return code 200')
        pages.append(resp.text)
        url = url.replace(f'/{pageNum}?', f'/{pageNum + 1}?')
        pageNum = pageNum + 1

    # exctract links from the pages
    links = []
    for page in pages:
        soup = BeautifulSoup(page, 'lxml')
        for link in soup.find_all('a'):
            if "/item/" in str(link):
                links.append('https://www.list.am' + link.get('href'))

    return links


if __name__ == "__main__":
    # Setup params
    errorCounter = 0
    refreshRate = 5*60
    price = 450
    roomsAmout = 3
    dbFileName = f'db_price_{price}_rooms_{roomsAmout}.txt'
    params = getParams(price, roomsAmout)
    loggingFileName = f'logs_price_{price}_rooms_{roomsAmout}.log'
    logging.basicConfig(format='%(asctime)s - %(levelname)s:%(message)s',
                        datefmt='%m/%d/%Y %I:%M:%S %p',
                        filename=loggingFileName,
                        level=logging.INFO)

    logging.info(f'\n\r\n\rStarted: price: {price}, rooms: {roomsAmout}, refresh rate(s): {refreshRate}')
    print(f'Started: price: {price}, rooms: {roomsAmout}, refresh rate(s): {refreshRate}')

    while True:
        try:
            # Get data from List am
            links = getAdsFromListAm(params)

            # If db doesn`t exist - create and repeat
            if not os.path.exists(dbFileName):
                logging.info(f'BD doesn`t exist. Creating: {dbFileName}')
                with open(dbFileName, "w") as db:
                    db.write('\n'.join(links))
                logging.info(f'DB Filled')

            # Get old data from db and compare with the new one
            with open(dbFileName, "r") as db:
                oldlinks = [line.rstrip('\n') for line in db]

            newLinks = findUniqueElements(oldlinks, links)
            if newLinks:
                logging.info(f'NEW LINKS!!! : {newLinks}')

            with open(dbFileName, "w") as db:
                db.write('\n'.join(links))

        except Exception as ex:
            logging.exception(ex)
            print(ex)

            errorCounter = errorCounter + 1
            if (errorCounter >= params['maxErrorsToStop']):
                logging.critical('Some serious shit happens. IT`S TIME TO STOP')
                break

        time.sleep(refreshRate)