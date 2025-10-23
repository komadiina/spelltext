import json


def generate_sql(input, output):
    with open(input) as f:
        data = json.load(f)

    sql = []
    for table, entries in data.items():
        for entry in entries:
            columns = []
            values = []
            for k, v in entry.items():
                columns.append(k)

                if v is str:
                    v = v.replace("'", "''")
                    v = f"'{v}'"

                values.append(v)

            for i in range(len(columns)):
                sql.append(
                    f"INSERT INTO {table} ({', '.join(columns)}) "
                    f"VALUES ({', '.join(values)});"
                )

    with open(output, "w") as f:
        f.write(sql)


if __name__ == "__main__":
    # parse cmd line args (--input X --output Y)
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument("-i", "--input")
    parser.add_argument("-o", "--output")
    args = parser.parse_args()
    print(args.input, args.output)

    generate_sql(args.input, args.output)
