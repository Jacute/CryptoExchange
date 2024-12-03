import argparse

from richbot.bot import start_bot as start_buybot
from randombot.bot import start_bot as start_randombot

def main():
    parser = argparse.ArgumentParser(description="CLI для управления торговыми ботами.")
    subparsers = parser.add_subparsers(dest="command", help="Команды для выполнения")
    
    def add_common_arguments(subparser):
        subparser.add_argument("--host", required=True, help="Хост для подключения")
        subparser.add_argument("--port", required=True, type=int, help="Порт для подключения")
        subparser.add_argument("--timeout", type=float, default=10, help="Таймаут работы бота")

    parser_sellbot = subparsers.add_parser("randombot", help="Запуск randombot")
    add_common_arguments(parser_sellbot)
    parser_sellbot.set_defaults(func=start_randombot)

    parser_buybot = subparsers.add_parser("richbot", help="Запуск richbot")
    add_common_arguments(parser_buybot)
    parser_buybot.set_defaults(func=start_buybot)

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
    elif args.command == "randombot":
        args.func(args.host, args.port, args.timeout)
    elif args.command == "richbot":
        args.func(args.host, args.port, args.timeout)
    else:
        args.func(args.host, args.port)

if __name__ == "__main__":
    main()
