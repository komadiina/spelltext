import csv
import json
from matplotlib import pyplot as plt
from tabulate import tabulate

from grpcmetrics import GRPCMetrics

DATAFILE = "ghz_results.csv"

# schema: csv(service_name, method_name, duration_ms, status, error|"")


def crunch() -> GRPCMetrics:
    with open(DATAFILE) as f:
        reader = csv.reader(f)
        return GRPCMetrics(reader)


def log(metrics):
    counts = metrics["counts"]
    avgs = metrics["averages"]
    maxs = metrics["max_times"]
    heaviest = metrics["heaviest_methods"]

    table = []
    for svc_name, svc_methods in counts.items():
        for method_name, count in svc_methods.items():
            heaviest_method = heaviest[svc_name]

            table.append(
                (
                    svc_name,
                    method_name,
                    count,
                    avgs[svc_name][method_name],
                    maxs[svc_name][method_name],
                    heaviest_method["pct"],
                )
            )

    table.sort(key=lambda x: x[3], reverse=True)
    tabulated = tabulate(
        table,
        headers=["service", "method", "count", "avg", "max", "pct"],
        colalign=["left", "left", "right", "right", "right", "right"],
    )

    print(tabulated)

    with open("table.txt", "w") as f:
        f.write(tabulated)


def main():
    metrics = crunch()
    log(metrics.summary())


if __name__ == "__main__":
    main()
