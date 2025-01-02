# CryptoExchange

Биржа, использующая СУБД [JacuteSQL](https://github.com/Jacute/JacuteSQL). Из-за отсутствия механизма транзакций могут быть потенциальные race condition между запросами.

## Запуск

```bash
git clone --recurse-submodules https://github.com/Jacute/CryptoExchange
docker compose -f docker-compose-balance.yml up
```

## Эндпоинты API

### Регистрация пользователя

Запрос:

*POST /user*
```json
{
    "username": string
}
```

Ответ:

```json
{
    "id": string,
    "token": string
}
```

### Получение списка ордеров

Запрос:

*GET /order*

Ответ:

```json
[
    {
        "order_id": int,
        "user_id": int,
        "lot_id": int,
        "quantity": float,
        "type": "sell" | "buy",
        "price": float,
        "closed": string
    },
    ...
]
```

### Удаление ордера

Запрос:

*DELETE /order*

*X-USER-KEY: string*

```json
{
    "order_id": int
}
```

### Получение информации о лотах

Запрос:

*GET /lot*

Ответ:

```json
[
    {
        "lot_id": int,
        "name": string
    },
    ...
]
```

### Получение информации о парах

Запрос:

*GET /pair*

Ответ:

```json
[
    {
        "pair_id": int,
        "sale_lot_id": int,
        "buy_lot_id": int
    },
    ...
]
```

### Баланс пользователя

*GET /balance*

*X-USER-KEY: string*

Ответ:

```json
[
    {
        "lot_id": int,
        "quantity": float
    },
    ...
]
```