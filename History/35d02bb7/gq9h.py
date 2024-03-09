import datetime

import requests
import psycopg2
from rich.traceback import install
from bs4 import BeautifulSoup
from setup.logger import log
pip
from config.config import DB_HOST, DB_NAME, DB_USER, DB_PASSWORD

install()
connection = psycopg2.connect(host=DB_HOST, dbname=DB_NAME, user=DB_USER, password=DB_PASSWORD)
cursor = connection.cursor()


class FearGreedBot:
    def __init__(self):
        self.url = "https://money.cnn.com/data/fear-and-greed/"

    def run_fear_and_greed_bot(self):
        try:
            fear_and_greed_page = requests.get(self.url)
            log.info(f"Status Code: {fear_and_greed_page.status_code}")
            if fear_and_greed_page.status_code == 200:
                soup = BeautifulSoup(fear_and_greed_page.content, 'html.parser')

                fear_and_greed_selector_value = soup.select("#needleChart > ul > li:nth-child(1)")
                log.info(fear_and_greed_selector_value[0].text)

                fear_and_greed_value = fear_and_greed_selector_value[0].text.split(": ")[1].split(" ")[0]
                log.info(f'Fear and Greed Value: {fear_and_greed_value}')

                query = f"INSERT INTO fear_and_greed_value(timestamp, fear_greed_value) VALUES ('{datetime.datetime.now()}', {fear_and_greed_value});"
                log.info(query)
                cursor.execute(query)
                connection.commit()
                connection.close()
            else:
                log.error("Error in fetching the fear and greed page.")
        except Exception as e:
            log.error(f"Error while trying to run the FII DII Scrapper Bot: {e}", exc_info=True)
            connection.rollback()


if __name__ == '__main__':
    fear_greed_bot = FearGreedBot()
    log.info("Starting running Fear and Greed Bot.")
    fear_greed_bot.run_fear_and_greed_bot()
    log.info("Finished running Fear and Greed Bot!")