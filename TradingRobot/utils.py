from api.api import rub_courses


def count_rubles(balance):
    rubles = 0
    
    for el in balance:
        rubles += rub_courses[list(rub_courses.keys())[el['lot_id'] - 1]] * el['quantity']
    
    return rubles