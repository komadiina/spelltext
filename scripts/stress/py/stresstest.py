import json
import subprocess
import time
from concurrent.futures import ThreadPoolExecutor, as_completed

from tabulate import tabulate
from grpcmetrics import GRPCMetrics

SERVICES = {}
results = []
n_calls = 50
interval = 10  # 10ms
connect_timeout = 1  # seconds


def call_grpc(svc_name, pkg, svc, method_name, params, target, call_index):
    json_data = json.dumps(params) if params else "{}"
    cmd = [
        "grpcurl",
        "-plaintext",
        "-connect-timeout",
        str(connect_timeout),
        "-d",
        json_data,
        target,
        f"{pkg}.{svc}/{method_name}",
    ]

    st = time.monotonic()
    result = subprocess.run(cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    duration_ms = int((time.monotonic() - st) * 1000)
    # return f"{svc_name},{pkg},{svc},{method_name},{call_index},{duration_ms},{result.returncode}"
    return [svc_name, pkg, svc, method_name, call_index, duration_ms, result.returncode]


def main():
    with open("data.json") as f:
        SERVICES = json.load(f)

    tasks = []

    with ThreadPoolExecutor(max_workers=20) as executor:
        for svc_name, svc_info in SERVICES.items():
            print(
                f"> stress testing svc {svc_name} ({len(svc_info['methods'])} methods)"
            )

            pkg = svc_info["pkg"]
            svc = svc_info["svc"]
            method_info = svc_info["methods"]
            target = svc_info["target"]

            for method_info in svc_info["methods"]:
                method_name = method_info["method"]
                params = method_info.get("parameters", {})
                for i in range(n_calls):
                    tasks.append(
                        executor.submit(
                            call_grpc,
                            svc_name,
                            pkg,
                            svc,
                            method_name,
                            params,
                            target,
                            i,
                        )
                    )
                    time.sleep(interval)

        for future in as_completed(tasks):
            results.append(future.result())


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-n",
        "--n-calls",
        type=int,
        default=n_calls,
        help=f"number of calls per method. default={n_calls}",
    )
    parser.add_argument(
        "-i",
        "--interval",
        type=float,
        default=interval / 1000,
        help=f"interval between calls, in milliseconds. default={interval}",
    )
    parser.add_argument(
        "-c",
        "--connect-timeout",
        type=float,
        default=connect_timeout,
        help=f"maximum connect timeout/wait-per-response, in seconds. default={connect_timeout}",
    )
    args = parser.parse_args()

    n_calls = args.n_calls
    interval = args.interval / 1000
    connect_timeout = args.connect_timeout

    print(
        f"initialized stresstest.py (n_calls={n_calls}, interval={interval}, connect_timeout={connect_timeout})"
    )

    main()

    metrics = GRPCMetrics(results)
    print("> finished!")
    print(f"> total calls: {metrics.total_calls()}")
    print(f"> total time taken: {metrics.total_execution_time()/10000}s")

    print("\n> per-service time taken:")
    table = []
    for svc_name, svc_time in metrics.per_service_times().items():
        table.append((svc_name, f"{(svc_time / 10000):.2f}"))
    print(tabulate(table, headers=["service", "duration (s)"]))

    print("\n> per-service method call details:")
    table = []
    for svc_name, methods in metrics.per_service_per_method_average().items():
        for method, avg in methods.items():
            table.append((svc_name, method, f"{avg}"))
    print(tabulate(table, headers=["service", "method", "avg (ms)"]))

    print("\n> heaviest methods:")
    table = []
    for svc_name, method in metrics.per_service_heaviest_method().items():
        table.append((svc_name, method[0], f"{method[1]}"))
    print(tabulate(table, headers=["service", "method", "duration (ms)"]))

    print(f"\n> per-service deviation:")
    table = []
    for svc_name, dev in metrics.service_deviation().items():
        table.append((svc_name, f"{dev:.2f}"))
    print(tabulate(table, headers=["service", "deviation (%)"]))

    results = ["host,pkg,service,method,call_n,duration_ms,resp_status"] + results
    with open("results.csv", "w") as f:
        f.write("\n".join([str(x)[1:-1] for x in results]))
        print("> results written to [./results.csv]")
