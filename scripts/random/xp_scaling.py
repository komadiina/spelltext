import math

formulas = {
    "linear": lambda l: l * 100 * 1.25,
    "lin": lambda l: l * 100 * 1.25 + l**2 * 66,
    "quadratic": lambda l: 25 * l * (5 * l - 4),
}

for name, formula in formulas.items():
    print(f'"{name}"')
    for i in range(1, 10):
        print(f"{i}: {formula(i)}")
