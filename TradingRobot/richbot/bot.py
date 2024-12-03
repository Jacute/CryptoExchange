import random
import pprint
import asyncio

from api.api import CryptoExchangeAPI, rub_courses, pairs
from utils import count_rubles


class RichBot:
    checked_order_ids = []
    def __init__(self, api: CryptoExchangeAPI, timeout: int):
        self.api = api
        self.timeout = timeout
    
    async def start(self):
        self.bot_id, token = self.api.register('richbot', 10)
        print("Token:", token)
        
        balance = self.api.get_balance()
        print(balance)
        print(count_rubles(balance))
        
        await asyncio.gather(self.bot_cycle())
            
    async def bot_cycle(self):
        for _ in range(10):
            orders = self.api.get_orders()
            
            for order in orders:
                id = order['id']
                
                if id in self.checked_order_ids:
                    continue
                
                price = float(order['price'])
                pair_id = int(order['pair_id'])
                user_id = int(order['user_id'])
                pair = pairs[pair_id]
                quantity = float(order['quantity'])
                closed = order["closed"]
                buy_lot, sell_lot = pair[0], pair[1]
                op_type = order['type']
                
                if closed == "0" and user_id != self.bot_id:
                    if op_type == 'sell':
                        if sell_lot == 'RUB' or buy_lot == 'RUB':
                            if sell_lot == 'RUB' and price < rub_courses[buy_lot]: # buy any wallet for rubles
                                pass
                            elif buy_lot == 'RUB' and price > rub_courses[sell_lot]: # buy rubles for any wallet
                                result = self.api.create_order(pair_id, quantity, price, "buy")
                                if result["status"] == "OK":
                                    print("Order bought. ID:", id, "Pair ID:", pair_id, "Quantity:", quantity, "Price:", price)
                                self.checked_order_ids.append(id)
                        elif rub_courses[sell_lot] * price * quantity < rub_courses[buy_lot] * quantity: # buy any wallet
                            result = self.api.create_order(pair_id, quantity, price, "buy")
                            if result["status"] == "OK":
                                print("Order bought. ID:", id, "Pair ID:", pair_id, "Quantity:", quantity, "Price:", price)
                                self.checked_order_ids.append(id)
                    elif op_type == 'buy':
                        if sell_lot == 'RUB' or buy_lot == 'RUB':
                            if buy_lot == 'RUB' and price > rub_courses[sell_lot]: # sell rubles for any wallet
                                pass
                            elif sell_lot == 'RUB' and price < rub_courses[buy_lot]: # sell any wallet for rubles
                                result = self.api.create_order(pair_id, quantity, price, "sell")
                                if result["status"] == "OK":
                                    print("Order sold. ID:", id, "Pair ID:", pair_id, "Quantity:", quantity, "Price:", price)
                                    self.checked_order_ids.append(id)
                        elif rub_courses[sell_lot] * price * quantity > rub_courses[buy_lot] * quantity: # sell any wallet
                            result = self.api.create_order(pair_id, quantity, price, "sell")
                            if result["status"] == "OK":
                                print("Order sold. ID:", id, "Pair ID:", pair_id, "Quantity:", quantity, "Price:", price)
                                self.checked_order_ids.append(id)
            await asyncio.sleep(self.timeout)
        balance = self.api.get_balance()
        print(balance)
        print(count_rubles(balance))


def start_bot(host: str, port: int, timeout: int):
    api = CryptoExchangeAPI(host, port)
    bot = RichBot(api, timeout)
    
    print("Starting sellbot on crypto exchange: {host}:{port}".format(host=host, port=port))
    asyncio.run(bot.start())
