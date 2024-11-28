import random
import pprint
import asyncio

from api.api import CryptoExchangeAPI, rub_courses, pairs


class RichBot:
    checked_order_ids = []
    def __init__(self, api: CryptoExchangeAPI, timeout: int):
        self.api = api
        self.timeout = timeout
    
    async def start(self):
        self.bot_id, token = self.api.register('richbot', 10)
        print("Token:", token)
        await asyncio.gather(self.bot_cycle(), self.balance_monitor())
        
    async def balance_monitor(self):
        while True:
            balance = self.api.get_balance()
            print("RichBot balance:")
            pprint.pprint(balance)
            await asyncio.sleep(60)
            
    async def bot_cycle(self):
        while True:
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
                        if sell_lot == 'RUB' and price < rub_courses[buy_lot]: # buy any wallet for rubles
                            pass
                            # result = self.api.create_order(pair_id, quantity, price, "buy")
                            # if result["status"] == "OK":
                            #     print("Order bought. ID:", id)
                            #     self.checked_order_ids.append(id)
                        elif buy_lot == 'RUB' and price > rub_courses[sell_lot]: # buy rubles for any wallet
                            result = self.api.create_order(pair_id, quantity, price, "buy")
                            if result["status"] == "OK":
                                print("Order bought. ID:", id, "Pair ID:", pair_id, "Quantity:", quantity, "Price:", price)
                                self.checked_order_ids.append(id)
                        elif price < rub_courses[buy_lot]: # buy any wallet
                            result = self.api.create_order(pair_id, quantity, price, "buy")
                            if result["status"] == "OK":
                                print("Order bought. ID:", id, "Pair ID:", pair_id, "Quantity:", quantity, "Price:", price)
                                self.checked_order_ids.append(id)
                    elif op_type == 'buy':
                        if buy_lot == 'RUB' and price > rub_courses[sell_lot]: # sell rubles for any wallet
                            pass
                            # result = self.api.create_order(pair_id, quantity, price, "sell")
                            # if result["status"] == "OK":
                            #     print("Order sold. ID:", id)
                            #     self.checked_order_ids.append(id)
                        elif sell_lot == 'RUB' and price < rub_courses[buy_lot]: # sell any wallet for rubles
                            result = self.api.create_order(pair_id, quantity, price, "sell")
                            if result["status"] == "OK":
                                print("Order sold. ID:", id, "Pair ID:", pair_id, "Quantity:", quantity, "Price:", price)
                                self.checked_order_ids.append(id)
                        elif price > rub_courses[buy_lot]: # sell any wallet
                            result = self.api.create_order(pair_id, quantity, price, "sell")
                            if result["status"] == "OK":
                                print("Order sold. ID:", id, "Pair ID:", pair_id, "Quantity:", quantity, "Price:", price)
                                self.checked_order_ids.append(id)
            await asyncio.sleep(self.timeout)


def start_bot(host: str, port: int, timeout: int):
    api = CryptoExchangeAPI(host, port)
    bot = RichBot(api, timeout)
    
    print("Starting sellbot on crypto exchange: {host}:{port}".format(host=host, port=port))
    asyncio.run(bot.start())
