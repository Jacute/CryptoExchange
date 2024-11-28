import random
import pprint
import asyncio

from api.api import *


class RandomBot:
    def __init__(self, api: CryptoExchangeAPI, timeout: int):
        self.api = api
        self.timeout = timeout
    
    async def start(self):
        self.bot_id, token = self.api.register('randombot', 10)
        print("Token:", token)
        await asyncio.gather(self.bot_cycle(), self.balance_monitor())
        
    async def balance_monitor(self):
        while True:
            balance = self.api.get_balance()
            print("Random bot balance:")
            pprint.pprint(balance)
            await asyncio.sleep(60)
            
    async def bot_cycle(self):
        while True:
            pair_id = random.randint(1, 42)
            quantity = random.randint(5, 25)
            price = random.randint(1, 200) / 10
            if random.randint(0, 1) == 0:
                op_type = 'sell'
            else:
                op_type = 'buy'
            
            data = self.api.create_order(pair_id=pair_id, quantity=quantity, price=price, type=op_type)
            print("Try to create order: {} pair_id, {} quantity, {} price, {} type".format(pair_id, quantity, price, op_type))
            if data["status"] == "OK":
                print("Sell order created:", data)
            else:
                print("Error:", data)
            
            await asyncio.sleep(self.timeout)


def start_bot(host: str, port: int, timeout: int):
    api = CryptoExchangeAPI(host, port)
    bot = RandomBot(api, timeout)
    
    print("Starting randombot on crypto exchange: {host}:{port}".format(host=host, port=port))
    asyncio.run(bot.start())
