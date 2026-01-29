import asyncio
import aiohttp
import time

URL = "http://a29954c0e8ae449e792cc60bac7fb00d-361890160.us-east-1.elb.amazonaws.com:3000/"
CONCURRENT_REQUESTS = 500 
TOTAL_REQUESTS = 10000  
RAMP_UP = 5             

async def fetch(session, url):
    start = time.time()
    try:
        async with session.get(url) as response:
            latency = time.time() - start
            return response.status, latency
    except Exception as e:
        return "ERROR", 0

async def worker(name, session, queue, results):
    while not queue.empty():
        await queue.get()
        status, latency = await fetch(session, URL)
        results.append((status, latency))
        queue.task_done()
        await asyncio.sleep(1) 

async def main():
    queue = asyncio.Queue()
    for _ in range(TOTAL_REQUESTS):
        queue.put_nowait(None)

    results = []

    async with aiohttp.ClientSession() as session:
        tasks = []
        for i in range(CONCURRENT_REQUESTS):
            await asyncio.sleep(RAMP_UP / CONCURRENT_REQUESTS) 
            task = asyncio.create_task(worker(f"worker-{i}", session, queue, results))
            tasks.append(task)

        start_time = time.time()
        await queue.join()  
        end_time = time.time()

    success = sum(1 for r in results if isinstance(r[0], int) and r[0] < 400)
    errors = sum(1 for r in results if r[0] == "ERROR" or (isinstance(r[0], int) and r[0] >= 400))
    latencies = [r[1] for r in results if r[1] > 0]

    print(f"Total de requisições: {TOTAL_REQUESTS}")
    print(f"Sucesso: {success}, Erros: {errors}")
    if latencies:
        print(f"Latência média: {sum(latencies)/len(latencies):.3f} s")
        print(f"Latência máxima: {max(latencies):.3f} s")
    print(f"Tempo total do teste: {end_time - start_time:.2f} s")

if __name__ == "__main__":
    asyncio.run(main())
