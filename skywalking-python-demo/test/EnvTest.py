import os


def test1():
    env_var = os.environ.get('SW_AGENT_COLLECTOR_BACKEND_SERVICES')
    print(env_var)


if __name__ == "__main__":
    test1()

