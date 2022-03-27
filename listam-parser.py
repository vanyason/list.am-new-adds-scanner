# Logic:
#
# 1. Parse setup config
# 2. Get data
# 3. Parse data
# 4. If file with data not exist:
#       Create file. Notify about file creation. Exit
#    Else
#       Compare new arrived data with the data in file. Save diffs locally. Update file
# 5. Show diffs if they exist

from bs4 import BeautifulSoup
import requests
import os
import logging
import time


def getParams(price, roomsAmount):
    headers = {'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36'}
    listAmUrl = f"https://www.list.am/ru/category/56/1?pfreq=1&po=1&n=0&price2={price}&crc=1&_a4={roomsAmount}"
    listAmUrlOriginal = f"https://www.list.am/category/56?pfreq=1&po=1&n=0&price2={price}&crc=1&_a4={roomsAmount}"
    return {'headers': headers, 'listAmUrl': listAmUrl, 'listAmUrlOriginal': listAmUrlOriginal}


def getAdsFromListAm(params):
    pages = []
    pageNum = 1
    url = params['listAmUrl']

    while True:
        resp = requests.get(url, headers=params['headers'])
        if (params['listAmUrlOriginal'] == resp.url):
            break
        if (resp.status_code != 200):
            raise ConnectionError('list am didn`t return code 200')
        pages.append(resp.text)
        url = url.replace(f'/{pageNum}?', f'/{pageNum + 1}?')
        pageNum = pageNum + 1

    links = []
    for page in pages:
        soup = BeautifulSoup(page, 'lxml')
        for link in soup.find_all('a'):
            if "/item/" in str(link):
                links.append('www.list.am' + link.get('href'))

    return links


if __name__ == "__main__":
    refreshRate = 5*60
    price = 450
    roomsAmout = 3
    dbFileName = f'db_price_{price}_rooms_{roomsAmout}.txt'
    loggingFileName = f'logs_price_{price}_rooms_{roomsAmout}.log'
    logging.basicConfig(format='%(asctime)s - %(levelname)s:%(message)s',
                        datefmt='%m/%d/%Y %I:%M:%S %p',
                        filename=loggingFileName,
                        level=logging.DEBUG)

    logging.info(f'\n\r\n\rStarted: price: {price}, rooms: {roomsAmout}, refresh rate(s): {refreshRate}')

    while True:
        try:
            params = getParams(price, roomsAmout)
            links = getAdsFromListAm(params)

            if not os.path.exists(dbFileName):
                logging.info(f'BD doesn`t exist. Creating: {dbFileName}')
                with open(dbFileName, "a") as db:
                    db.write('\n'.join(links))
                logging.info(f'DB Filled')
                break

            time.sleep(refreshRate)

        except Exception as ex:
            logging.error(ex)
