from skywalking import agent, config

# 定义Skywalking代理
config.service_name = "my-service"
config.logging_level = "DEBUG"
agent.start()


# 定义一个简单的函数
def hello(name):
    print(f"Hello, {name}!")


# 将函数包装到Skywalking跨度中
with agent.start_span(operation_name="test-span") as span:
    span.set_tag("Hello", "World")
    hello("Skywalking")
    # 获取所有跨度对象的列表
    span_list = agent.get_spans()
    print(span_list)

# 停止Skywalking代理
agent.stop()
