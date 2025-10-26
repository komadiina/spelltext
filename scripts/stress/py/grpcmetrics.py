from collections import defaultdict


class GRPCMetrics:
    def __init__(self, data):
        self.data = data
        self.total_time = 0
        self.service_times = defaultdict(int)
        self.method_totals = defaultdict(int)
        self.method_counts = defaultdict(int)
        self._process_data()

    def _process_data(self):
        for row in self.data:
            service = row[0]
            method = row[3]
            duration = row[5]

            # total time
            self.total_time += duration

            # per-service time
            self.service_times[service] += duration

            # per-method calcs
            self.method_totals[method] += duration
            self.method_counts[method] += 1

    def total_calls(self):
        return len(self.data)

    def total_execution_time(self):
        return self.total_time

    def per_service_times(self):
        return dict(self.service_times)

    def per_method_average(self):
        return {
            m: self.method_totals[m] / self.method_counts[m] for m in self.method_totals
        }

    def per_service_per_method_average(self):
        service_method_totals = defaultdict(lambda: defaultdict(int))
        service_method_counts = defaultdict(lambda: defaultdict(int))

        for row in self.data:
            service = row[0]
            method = row[3]
            duration = row[5]

            service_method_totals[service][method] += duration
            service_method_counts[service][method] += 1

        result = {}
        for svc in service_method_totals:
            result[svc] = {
                m: service_method_totals[svc][m] / service_method_counts[svc][m]
                for m in service_method_totals[svc]
            }

        return result

    def per_service_heaviest_method(self):
        service_method_totals = defaultdict(lambda: defaultdict(int))
        for row in self.data:
            service = row[0]
            method = row[3]
            duration = row[5]
            service_method_totals[service][method] += duration

        result = {}
        for svc, methods in service_method_totals.items():
            if methods:
                heaviest = max(methods.items(), key=lambda x: x[1])
                result[svc] = heaviest  # (method_name, total_duration)
            else:
                result[svc] = (None, 0)

        return result

    def service_deviation(self):
        n_services = len(self.service_times)
        if n_services == 0:
            return {}
        avg_time = self.total_time / n_services
        return {
            svc: ((time - avg_time) / avg_time) * 100
            for svc, time in self.service_times.items()
        }
