import asyncio

async def foo():
    raise ValueError("Foo value error")
    return("Foo finished")

async def bar():
    await asyncio.sleep(1)
    return("Bar finished")

async def main():
    foo_task = asyncio.create_task(foo(), name="Exception task")
    bar_task = asyncio.create_task(bar(), name="Waiting task")
    try:
        done, pending = await asyncio.wait(
            [foo_task, bar_task],
            return_when=asyncio.ALL_COMPLETED
        )
        for task in done:
            name = task.get_name()
            print(f"Done: {name}")
        for task in pending:
            task.cancel()
    except Exception as e:
        print(f"Exception caught: {e}")

asyncio.run(main())