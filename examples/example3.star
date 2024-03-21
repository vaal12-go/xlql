

def factorial(n):
    ret = 1
    for i in range(1, n+1):
        ret = ret * i
    return ret

def pow(x, n):
    if n==0: return 1
    ret = 1
    for i in range(n):
        ret = ret * x
    return ret

def taylor_sin(x, polynom_grade):
    ret = 0
    for n in range(polynom_grade):
        sign = pow(-1, n)
        nominator = pow(x, (2*n+1))
        denominator = factorial(2*n+1)
        ret = ret + (sign*nominator)/denominator
    return ret

sinDB = open_db("sin_db.sqlite")

sin_q = sinDB.create_table(
    name = "sinus_table",
    columns = {
        "x": "NUMERIC",
        "sin_of_x": "NUMERIC",
    }
)

for x1 in range(-300, 400):
    real_x = x1*0.01
    sin_x = taylor_sin(real_x, 20)
    sin_q.add_row(real_x, sin_x)

sin_q.save_to_excel("sinus_graph_data.xlsx", "sinus_data")