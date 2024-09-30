db = open_db("file:memdb1?mode=memory&cache=shared")
fName = "./test_cases/01.Simple_examples.xlsx"
shName = "DateTime"

q1 = db.load_excel_sheet(
    fName,
    shName,
    drop_table = True,
    cols = {
        "DateTime" : {
            "type" : "date",
            "format" : "2006-01-02 03:04:05"
        },
        "NumKey" : {
            "type" : "numeric"
        }
    }
)

# q1.print()

sha1 = q1.get_sha()

# print("Sha1:", sha1)

test_result = "sha of query:"+sha1

# q1.save_to_excel("test_xl1.xlsx")

starlark_result_sha = sha1
