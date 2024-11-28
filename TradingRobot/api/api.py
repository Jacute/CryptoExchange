import requests
import random
import string

pairs = {
    1: ['RUB', 'BTC'],
    2: ['RUB', 'ETH'],
    3: ['RUB', 'USDT'],
    4: ['RUB', 'USDC'],
    5: ['RUB', 'SOL'],
    6: ['RUB', 'DOGE'],
    7: ['BTC', 'RUB'],
    8: ['BTC', 'ETH'],
    9: ['BTC', 'USDT'],
    10: ['BTC', 'USDC'],
    11: ['BTC', 'SOL'],
    12: ['BTC', 'DOGE'],
    13: ['ETH', 'RUB'],
    14: ['ETH', 'BTC'],
    15: ['ETH', 'USDT'],
    16: ['ETH', 'USDC'],
    17: ['ETH', 'SOL'],
    18: ['ETH', 'DOGE'],
    19: ['USDT', 'RUB'],
    20: ['USDT', 'BTC'],
    21: ['USDT', 'ETH'],
    22: ['USDT', 'USDC'],
    23: ['USDT', 'SOL'],
    24: ['USDT', 'DOGE'],
    25: ['USDC', 'RUB'],
    26: ['USDC', 'BTC'],
    27: ['USDC', 'ETH'],
    28: ['USDC', 'USDT'],
    29: ['USDC', 'SOL'],
    30: ['USDC', 'DOGE'],
    31: ['SOL', 'RUB'],
    32: ['SOL', 'BTC'],
    33: ['SOL', 'ETH'],
    34: ['SOL', 'USDT'],
    35: ['SOL', 'USDC'],
    36: ['SOL', 'DOGE'],
    37: ['DOGE', 'RUB'],
    38: ['DOGE', 'BTC'],
    39: ['DOGE', 'ETH'],
    40: ['DOGE', 'USDT'],
    41: ['DOGE', 'USDC'],
    42: ['DOGE', 'SOL']
}

# COURSE FOR RUB

rub_courses = {
    "RUB": 1,
    "BTC": 0.1,
    "ETH": 0.2,
    "USDT": 2,
    "SOL": 8,
    "DOGE": 0.5,
    "USDC": 5
}

class CryptoExchangeAPI:
    session = requests.Session()
    def __init__(self, host: str, port: int):
        self.url = f"http://{host}:{port}"
    
    def _generate_random_str(self, length: int):
        result = ''
        
        for i in range(length):
            result += random.choice(string.ascii_letters + string.digits)
        
        return result
    
    def register(self, username, random_len=0) -> tuple[int, str]:
        if random_len != 0:
            username += self._generate_random_str(random_len)
        
        response = self.session.post(self.url + '/user', json={"username": username})
        response.raise_for_status()
        data = response.json()
        self.session.headers.setdefault('X-USER-TOKEN', data["token"])
        
        return data["id"], data["token"]
    
    def get_balance(self):
        response = self.session.get(self.url + '/balance')
        response.raise_for_status()
        return response.json()
    
    def create_order(self, pair_id: int, quantity: float, price: float, type: str):
        if type not in ['buy','sell']:
            raise ValueError("Invalid order type. Must be 'buy' or'sell'.")
    
        request_data = {
            "pair_id": pair_id,
            "quantity": quantity,
            "price": price,
            "type": type,
        }
    
        response = self.session.post(self.url + '/order', json=request_data)
        response.raise_for_status()
        return response.json()
    
    def get_orders(self):
        response = self.session.get(self.url + '/order')
        response.raise_for_status()
        return response.json()