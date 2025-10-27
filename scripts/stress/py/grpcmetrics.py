import json
from typing import Dict, List
from csv import reader

# { "service_a": {"method_a": [0.1, 0.2, 0.3]} }


class GRPCMetrics:
    def __init__(self, data: reader):
        self.data = data
        self.service_methods: Dict[str, Dict[str, List[float]]] = {}
        self._aggregate()

    def _aggregate(self):
        for row in self.data:
            if not row:
                continue
            # trim whitespace
            row = [c.strip() for c in row]
            service, method, duration_s = row[0], row[1], row[2]

            try:
                duration = float(duration_s)
            except ValueError as e:
                continue

            if service not in self.service_methods:
                self.service_methods[service] = {}

            if method not in self.service_methods[service]:
                self.service_methods[service][method] = []

            self.service_methods[service][method].append(duration)

    def get_invocation_counts(self):
        def f(service_methods: dict):
            d = {}
            for method_name, method_durations in service_methods.items():
                d[method_name] = len(method_durations)

            return d

        counts = {}
        for svc_name, svc_methods in self.service_methods.items():
            counts[svc_name] = f(svc_methods)

        return counts

    def get_average_times(self):
        def f(service_methods: dict):
            d = {}
            for name, durations in service_methods.items():
                d[name] = sum(durations) / len(durations)
            return d

        averages = {}
        for svc_name, svc_methods in self.service_methods.items():
            averages[svc_name] = f(svc_methods)

        return averages

    def get_max_times(self):
        def f(service_methods: dict):
            d = {}
            for name, durations in service_methods.items():
                d[name] = max(durations)
            return d

        maximums = {}
        for svc_name, svc_methods in self.service_methods.items():
            maximums[svc_name] = f(svc_methods)

        return maximums

    def get_heaviest_methods(self, max_times, averages):
        t = {}
        for svc_name, svc_times in max_times.items():
            name, dur = max(max_times[svc_name].items(), key=lambda x: x[1])
            pct = dur / averages[svc_name][name]
            t[svc_name] = {"method": name, "duration": dur, "pct": pct}
        return t

    def summary(self) -> dict:
        counts = self.get_invocation_counts()
        averages = self.get_average_times()
        max_times = self.get_max_times()
        heaviest_methods = self.get_heaviest_methods(max_times, averages)

        return {
            "counts": counts,
            "averages": averages,
            "max_times": max_times,
            "heaviest_methods": heaviest_methods,
        }
