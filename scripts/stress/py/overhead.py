import subprocess
import time
import timeit

cmd = ["grpcurl", "-version"]

n = 50
durations = []


durations.append(
    timeit.timeit(
        lambda: subprocess.run(
            cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL
        ),
        number=n,
    )
)
# subprocess.run(cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
# duration_ms = (time.monotonic() - start) * 1000
# durations.append(duration_ms)

avg_overhead = sum(durations) / n
print(f"Average grpcurl subprocess overhead: {avg_overhead:.2f} ms")
