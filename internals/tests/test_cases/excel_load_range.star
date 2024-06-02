db = open_db("file:memdb1?mode=memory&cache=shared")
fName = "./test_cases/01.Simple_examples.xlsx"
shName = "TableRange"

q1 = db.load_excel_sheet(
    file_name = fName,
    sheet_name = shName,
    table_range = "C4:F8",
    drop_table = True
)

# q1.print()

starlark_result_sha = q1.get_sha()

# print("Sha1:", sha1)

# test_result = "sha of query:"+starlark_result_sha

# if starlark_result_sha == "2d3e7fc8284ea9805a41fd54cbbf07aec3365b925fbabe90a29769dc7a997109":
#     test_result = "OK"

