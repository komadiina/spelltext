import json
import os
import subprocess
import time
from concurrent.futures import ThreadPoolExecutor, as_completed

from tabulate import tabulate
from grpcmetrics import GRPCMetrics

SERVICES = {}
PROJECT_ROOT = "~/Desktop/Moje/Faks/diplomski"
results = []
ghz_outputs = []
n_calls = 500
interval = 10  # 10ms
connect_timeout = 1  # seconds
rps = 400
imports = ["proto"]
import_paths_arg = ""
output_file = "ghz_results"
output_format = "csv"


def call_grpc(svc_name, pkg, svc, method_name, params, target, proto_dir):
    json_data = json.dumps(params) if params else "{}"
    # resolve proto path passed in JSON (proto_dir may be "proto/build/build.proto")
    proto_path = proto_dir
    if not os.path.isabs(proto_path):
        proto_path = os.path.join(PROJECT_ROOT, proto_path)
    proto_path = os.path.normpath(os.path.expanduser(proto_path))

    cmd = [
        "ghz",
        "--insecure",
        "--proto",
        proto_path,
        "--import-paths",
        f"{import_paths_arg}",
        "--rps",
        str(rps),
        "-c",
        str(n_calls),
        "-d",
        json_data,
        target,
        "--format",
        "csv",
        "--call",
        f"{pkg}.{svc}/{method_name}",
    ]

    results = subprocess.run(cmd, capture_output=True, text=True).stdout.strip()

    if output_format == "json":
        ghz_outputs.append(json.loads(results))
    elif output_format == "csv":
        lines = [ln.strip() for ln in results.splitlines() if ln.strip()]
        prefixed_rows = []

        for ln in lines:
            # skip ghz header lines
            if ln.startswith("duration"):
                continue
            prefixed_rows.append(f"{svc_name},{method_name},{ln}")

        ghz_outputs.extend(prefixed_rows)


def main():
    with open(f"{os.getcwd()}/scripts/stress/py/data.json") as f:
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
                            svc_info["protoDir"],
                        )
                    )
                    time.sleep(interval)

        for future in as_completed(tasks):
            pass


if __name__ == "__main__":
    calculate = False
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
    parser.add_argument(
        "-r",
        "--rps",
        type=int,
        default=rps,
        help=f"requests per second. default={rps}",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        default=output_file,
        help=f"ghz results output file. default={output_file}",
    )
    parser.add_argument(
        "-d",
        "--imports",
        type=str,
        default=imports,
        help=f"ghz imports. default={imports}",
    )

    args = parser.parse_args()

    n_calls = args.n_calls
    interval = args.interval / 1000
    connect_timeout = args.connect_timeout

    import_paths = []
    for p in imports:
        p_full = p
        if not os.path.isabs(p_full):
            p_full = os.path.join(PROJECT_ROOT, p_full)
        import_paths.append(os.path.normpath(os.path.expanduser(p_full)))
    import_paths_arg = ",".join(import_paths)

    print(
        f"initialized stresstest.py (n_calls={n_calls}, interval={interval*1000}s, rps={rps}, connect_timeout={connect_timeout}, output_format={output_format})"
    )

    st = time.monotonic()
    main()
    dur = time.monotonic() - st

    print("> finished!")
    print(f"> time taken: {dur:.2f}s")
    output_file = f"{output_file}.{output_format}"
    fname = f"{os.getcwd()}/scripts/stress/py/{output_file}"
    with open(fname, "w") as f:
        if output_format == "json":
            f.write(json.dumps(ghz_outputs))
        elif output_format == "csv":
            f.write("service,method,duration (ms),status,error\n")
            f.write("\n".join(ghz_outputs))

        print(f"> results written to {fname}")
